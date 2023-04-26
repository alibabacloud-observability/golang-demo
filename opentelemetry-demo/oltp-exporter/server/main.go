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

// Sample contains a simple http server that exports to the OpenTelemetry agent.
package main

import (
	"context"
	"fmt"
	"github.com/alibabacloud-observability/golang-demo/opentelemetry-demo/otlp-exporter/common"
	"go.opentelemetry.io/otel/codes"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// SpanNameVariety SpanName 发散程度（多少个不同值）
const SpanNameVariety = 1000

// AttrValueVariety 属性值发散程度（多少个不同值）
const AttrValueVariety = 10000

// AttrMaxLen AttrMinLen tag value 长度范围
const AttrMaxLen = 10000
const AttrMinLen = 1000

// SpanNameMaxLen SpanNameMinLen span name 长度范围
const SpanNameMaxLen = 64
const SpanNameMinLen = 32

var avaAttrValue = [AttrValueVariety]string{}
var avaSpanName = [SpanNameVariety]string{}

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
	log.Println("start to connect to server")
	traceExp, err := otlptrace.New(ctx, traceClient)
	handleErr(err, "Failed to create the collector trace exporter")

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(common.ServerServiceName),
			semconv.HostNameKey.String(common.ServerServiceHostName),
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

func initTraceDemoData() {
	for i := 0; i < len(avaAttrValue); i++ {
		avaAttrValue[i] = common.GenStrWithRandomLen(AttrMinLen, AttrMaxLen)
	}
	for i := 0; i < len(avaSpanName); i++ {
		avaSpanName[i] = common.GenStrWithRandomLen(SpanNameMinLen, SpanNameMaxLen)
	}
}

func main() {
	shutdown := initProvider()
	defer shutdown()

	//meter := global.Meter("demo-server-meter")
	serverAttribute := attribute.String("server-attribute", "foo")
	fmt.Println("start to gen chars for trace data")
	initTraceDemoData()
	fmt.Println("gen trace data done")
	tracer := otel.Tracer(common.TraceInstrumentationName)

	// create a handler wrapped in OpenTelemetry instrumentation
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//  random sleep to simulate latency
		var sleep int64
		switch modulus := time.Now().Unix() % 5; modulus {
		case 0:
			sleep = rng.Int63n(2000)
		case 1:
			sleep = rng.Int63n(15)
		case 2:
			sleep = rng.Int63n(917)
		case 3:
			sleep = rng.Int63n(87)
		case 4:
			sleep = rng.Int63n(1173)
		}
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(serverAttribute)

		actionChild(tracer, ctx, sleep)

		w.Write([]byte("Hello World"))
	})
	wrappedHandler := otelhttp.NewHandler(handler, "/hello")

	// serve up the wrapped handler
	http.Handle("/hello", wrappedHandler)
	http.ListenAndServe(":7080", nil)

}

func actionChild(tracer trace.Tracer, ctx context.Context, sleep int64) {
	_, subSpan := tracer.Start(ctx, getRandomSpanName())
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	subSpan.SetStatus(codes.Ok, "success")
	serverAttribute := attribute.String("attr1", getRandomAttrValue())
	subSpan.SetAttributes(serverAttribute)
	subSpan.End()
}

func getRandomAttrValue() string {
	n := rand.Int63n(int64(len(avaAttrValue)))
	return avaAttrValue[n]
}

func getRandomSpanName() string {
	n := rand.Int63n(int64(len(avaSpanName)))
	return avaSpanName[n]
}
