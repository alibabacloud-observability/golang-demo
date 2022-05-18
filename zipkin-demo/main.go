package main

import (
	"time"
)

func main() {
	//手动埋点
	ExampleNewTracer()
	// 用http框架埋点
	HttpExample()
	time.Sleep(5 * time.Second)
}
