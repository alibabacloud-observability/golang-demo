package cli

import (
	"fmt"
	"github.com/alibabacloud-observability/golang-demo/jaeger-demo/xtrace"
	"github.com/opentracing/opentracing-go"
	"log"
	"time"
)

var tracer opentracing.Tracer = nil

func child(pSpan opentracing.Span) {
	span := tracer.StartSpan("child", opentracing.ChildOf(pSpan.Context()))
	span.SetTag("type", "child")
	defer span.Finish()
	log.Println("cli child start")
	time.Sleep(2 * time.Second)
	log.Println("cli child end")
}

func parent(pSpan opentracing.Span) {
	span := tracer.StartSpan("parent", opentracing.ChildOf(pSpan.Context()))
	defer span.Finish()
	log.Println("cli parent start")
	span.SetTag("type", "parent")
	time.Sleep(1 * time.Second)
	child(span)
	log.Println("cli parent end")
	fmt.Println("cli done")
}

func CliMain() {
	tracer = xtrace.NewJaegerTracer(xtrace.CliTracerServiceName)
	span := tracer.StartSpan("cli_main")
	span.SetTag("type", "http")
	log.Println("cli http")
	go parent(span)
	span.Finish()
}
