package grpc

import (
	"context"
	"fmt"
	"github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"net"
	"time"
	pb "github.com/alibabacloud-observability/golang-demo/jaeger-demo/grpc/helloworld"
)

func StartGrpcServer(tracer opentracing.Tracer) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	// Register reflection service on gRPC server.
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer)),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer)))
	pb.RegisterGreeterServer(s, &grpcServer{})

	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}

type grpcServer struct{}

// 实现 helloworld.GreeterServer
func (s *grpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	//fmt.Println("receive req : %v \n", *in)

	//start a new span, eg.(mysql)
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			mysqlSpan := tracer.StartSpan("SQL FindUserTable", opentracing.ChildOf(pctx))
			mysqlSpan.SetTag("db.statement", "select * from user ...")
			//do mysql operations
			time.Sleep(time.Millisecond * 100)

			mysqlSpan.Finish()
		}

	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
