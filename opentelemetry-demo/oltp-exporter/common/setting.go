package common

const (
	TraceInstrumentationName = "otlp-demo-tracer"
)

// Server端配置
const (
	ServerServiceName     = "otlp-demo-server"
	ServerServiceHostName = "server.host.name"
	DefaultServerEndpoint = "http://0.0.0.0:7080/hello"
)

// 通过gRPC协议上报Trace数据的Client
const (
	ClientGrpcExportServiceName     = "otlp-grpc-export-demo-client"
	ClientGrpcExportServiceHostName = "client.grpc.host.name"
)

// 通过HTTP协议上报Trace数据的Client
const (
	ClientHttpExportServiceName     = "otlp-http-export-demo-client"
	ClientHttpExportServiceHostName = "client.http.host.name"
	TraceExportHttpEndpoint         = "tracing-analysis-dc-hz.aliyuncs.com" // 请勿添加"http://"
	TraceExportHttpURLPath          = "/adapt_xxxx@xxxx@xxxx/api/otlp/traces"
)
