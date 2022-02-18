package http

import (
	"flag"
	"fmt"
	"github.com/alibabacloud-observability/golang-demo/jaeger-demo/xtrace"
)

var (
	serverPort = flag.String("port", "8000", "server port")
)

func HttpTest() {
	go startServer()

	runClient(xtrace.NewJaegerTracer(xtrace.HttpClientName))
	fmt.Println("http done")

}

func startServer() {
	runServer(xtrace.NewJaegerTracer(xtrace.HttpServerName))
}
