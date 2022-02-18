package grpc

import (
	"fmt"
	"github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	pb "github.com/alibabacloud-observability/golang-demo/jaeger-demo/grpc/helloworld"
)

func NewClient(tracer opentracing.Tracer) (pb.GreeterClient, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to service: %v", err)
	}
	return pb.NewGreeterClient(conn), nil
}
