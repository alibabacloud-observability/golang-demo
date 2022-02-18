package main

import (
	"github.com/alibabacloud-observability/golang-demo/jaeger-demo/cli"
	"github.com/alibabacloud-observability/golang-demo/jaeger-demo/grpc"
	"github.com/alibabacloud-observability/golang-demo/jaeger-demo/http"
	"time"
)

func main() {
	http.HttpTest()
	grpc.GrpcTest()
	cli.CliMain()
	time.Sleep(30 * time.Second)
}
