package main

import (
	"context"
	"encoding/json"
	"github.com/openzipkin/zipkin-go"
	"time"
)

func doSomeWork(context.Context) {}

func ExampleNewTracer() {
	tracer := GetTracer("demoService", "172.20.23.100:80")
	// tracer can now be used to create spans.
	span := tracer.StartSpan("some_operation")
	// ... do some work ...
	span.Finish()

	childSpan := tracer.StartSpan("some_operation2", zipkin.Parent(span.Context()))
	// ... do some work ...
	var events = make(map[string]string)
	events["event"] = "error"
	events["stack"] = "Runtime Exception: unable to find userid"
	jsonStr, err := json.Marshal(events)
	if err == nil {
		childSpan.Annotate(time.Now(), string(jsonStr))
	}
	childSpan.Finish()

	span.Finish()

	// Output:
}
