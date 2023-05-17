package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"log"
	"os"
	"time"
)

const (
	SERVICE_NAME = "<your-service-name>" // 应用名称
	HOST_NAME    = ""                    // 未设置时则默认为本机主机名
)

const (
	// 是否使用gRPC上报数据
	EXPORT_BY_GRPC = true
	GRPC_ENDPOINT  = "<gRPC-endpoint>"
	GRPC_TOKEN     = "<gRPC-token>"
)

const (
	HTTP_ENDPOINT = "<HTTP-endpoint>"
	HTTP_URL_PATH = "<HTTP-URL-path>"
)

// 设置应用资源
func newResource(ctx context.Context, serviceName string, hostName string) *resource.Resource {
	// hostname默认值为本机主机名
	if len(hostName) == 0 {
		hostName, _ = os.Hostname()
	}

	r, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(), // runtime信息 process.runtime.name: go/gc, process.runtime.version: go1.20.1s
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.HostNameKey.String(hostName),
		),
	)

	if err != nil {
		log.Fatalf("%s: %v", "Failed to create OpenTelemetry resource", err)
	}
	return r
}

// 通过
func newGrpcExporterAndSpanProcessor(ctx context.Context) (*otlptrace.Exporter, sdktrace.SpanProcessor) {
	headers := map[string]string{"Authentication": GRPC_TOKEN}

	traceExporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(GRPC_ENDPOINT),
			otlptracegrpc.WithHeaders(headers),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
			otlptracegrpc.WithCompressor(gzip.Name)),
	)

	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the OpenTelemetry trace exporter", err)
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	return traceExporter, batchSpanProcessor
}

func newHTTPExporterAndSpanProcessor(ctx context.Context) (*otlptrace.Exporter, sdktrace.SpanProcessor) {

	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(HTTP_ENDPOINT),
		otlptracehttp.WithURLPath(HTTP_URL_PATH),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(1)))

	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the OpenTelemetry trace exporter", err)
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	return traceExporter, batchSpanProcessor
}

// OpenTelemetry初始化方法
func initOpenTelemetry() func() {
	ctx := context.Background()

	var traceExporter *otlptrace.Exporter
	var batchSpanProcessor sdktrace.SpanProcessor

	if EXPORT_BY_GRPC {
		traceExporter, batchSpanProcessor = newGrpcExporterAndSpanProcessor(ctx)
	} else {
		traceExporter, batchSpanProcessor = newHTTPExporterAndSpanProcessor(ctx)
	}

	otelResource := newResource(ctx, SERVICE_NAME, HOST_NAME)

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(otelResource),
		sdktrace.WithSpanProcessor(batchSpanProcessor))

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExporter.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
