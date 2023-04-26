// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample contains a simple client that periodically makes a simple http request
// to a server and exports to the OpenTelemetry service.
package main

import (
	"context"
	"github.com/alibabacloud-observability/golang-demo/opentelemetry-demo/otlp-exporter/common"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"

	//"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() func() {
	ctx := context.Background()

	otelAgentAddr, xtraceToken, ok := common.ObtainXTraceInfo()

	if !ok {
		log.Fatalf("Cannot init OpenTelemetry, exit")
		os.Exit(-1)
	}

	headers := map[string]string{"Authentication": xtraceToken}
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
		otlptracegrpc.WithHeaders(headers), // 鉴权信息
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	//traceClientHttp := otlptracehttp.NewClient(
	//	otlptracehttp.WithEndpoint("127.0.0.1:8080"),
	//	otlptracehttp.WithURLPath("/adapt_xxxxx/api/otlp/traces"),
	//	otlptracehttp.WithInsecure())
	//otlptracehttp.WithCompression(1))

	traceExp, err := otlptrace.New(ctx, traceClient)
	handleErr(err, "Failed to create the collector trace exporter")

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(common.ClientServiceName),
			semconv.HostNameKey.String(common.ClientServiceHostName),
		),
	)
	handleErr(err, "failed to create resource")

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	log.Println("OTEL init success")

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func main() {
	log.Printf("client start")
	shutdown := initProvider()
	defer shutdown()

	tracer := otel.Tracer(common.TraceInstrumentationName)

	method, _ := baggage.NewMember("method", "repl")
	client, _ := baggage.NewMember("client", "cli")
	bag, _ := baggage.New(method, client)

	defaultCtx := baggage.ContextWithBaggage(context.Background(), bag)
	for {
		ctx, span := tracer.Start(defaultCtx, "ExecuteRequest")
		makeRequest(ctx)
		span.End()
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func makeRequest(ctx context.Context) {
	demoServerAddr, ok := os.LookupEnv("DEMO_SERVER_ENDPOINT")
	if !ok {
		demoServerAddr = common.DefaultServerEndpoint
	}

	// Trace an HTTP client by wrapping the transport
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Make sure we pass the context to the request to avoid broken traces.
	req, err := http.NewRequestWithContext(ctx, "GET", demoServerAddr, nil)
	if err != nil {
		handleErr(err, "failed to http request")
	}

	// All requests made with this client will create spans.
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		res.Body.Close()
	}
}
