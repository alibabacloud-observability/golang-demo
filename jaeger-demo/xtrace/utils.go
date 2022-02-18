package xtrace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"time"
)

func NewJaegerTracer(service string) opentracing.Tracer {
	if AgentSwitch == true {
		return NewJaegerTracerAgent(service)
	}
	return NewJaegerTracerDirect(service)
}

// NewJaegerTracerDirect 通过 SDK 直接上报 (HTTP)
func NewJaegerTracerDirect(service string) opentracing.Tracer {
	sender := transport.NewHTTPTransport(
		TracingAnalysisEndpoint,
		transport.HTTPTimeout(10 * time.Second), // defaults to 5 seconds
		transport.HTTPBatchSize(50), // defaults to 100
	)
	tracer, _ := jaeger.NewTracer(service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger)),
	)
	return tracer
}

// NewJaegerTracerAgent 通过 Jaeger Agent 上报
func NewJaegerTracerAgent(service string) opentracing.Tracer {
	sender, _ := jaeger.NewUDPTransport("",0)
	tracer, _ := jaeger.NewTracer(service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger)),
	)
	return tracer
}