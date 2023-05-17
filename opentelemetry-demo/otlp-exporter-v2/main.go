package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"time"
)

func main() {
	shutdown := initOpenTelemetry()
	defer shutdown()

	ctx := context.Background()
	clientParentSpan(ctx)

	time.Sleep(10 * time.Second)
}

func clientParentSpan(ctx context.Context) {
	tracer := otel.Tracer("otel-go-tracer")
	ctx, span := tracer.Start(ctx, "parent span")
	//span.SetStatus(codes.Error, "error")
	span.SetStatus(codes.Ok, "Success")
	clientChildSpan(ctx)
	fmt.Println(span.SpanContext().TraceID())
	fmt.Println(span.SpanContext().SpanID())
	span.End()
}

func clientChildSpan(ctx context.Context) {
	tracer := otel.Tracer("otel-go-tracer")
	ctx, span := tracer.Start(ctx, "child span")
	//span.SetStatus(codes.Error, "error")
	span.SetStatus(codes.Ok, "Success")
	fmt.Println(span.SpanContext().TraceID())
	fmt.Println(span.SpanContext().SpanID())
	span.End()
}
