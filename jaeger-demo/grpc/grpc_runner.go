package grpc

import (
	"context"
	"fmt"
	"github.com/alibabacloud-observability/golang-demo/jaeger-demo/xtrace"
	"github.com/opentracing/opentracing-go"
	"time"
	pb "github.com/alibabacloud-observability/golang-demo/jaeger-demo/grpc/helloworld"
)

const (
	addr = "0.0.0.0:19090"
)

func GrpcTest() {
	tracer := xtrace.NewJaegerTracer(xtrace.GrpcServerName)
	opentracing.SetGlobalTracer(tracer)
	go StartGrpcServer(tracer)
	time.Sleep(1 * time.Second)
	c, _ := NewClient(xtrace.NewJaegerTracer(xtrace.GrpcClientName))
	ctx := context.Background()
	c.SayHello(ctx, &pb.HelloRequest{Name: "123"})
	fmt.Println("grpc done")
}
