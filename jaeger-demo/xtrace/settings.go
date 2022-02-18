package xtrace

// TracingAnalysisEndpoint SDK上报需要：设置链路追踪的网关（不同region对应不同的值，从http://tracing.console.aliyun.com/ 的配置查看中获取）
const TracingAnalysisEndpoint = "http://tracing-analysis-dc-hz.aliyuncs.com/adapt_xxxxxx_xxxxx/api/traces"

// Agent上报：true （需要本地启动 jaeger-agent） SDK上报：false（默认值）
const AgentSwitch = false

// 链路追踪控制台对应的应用名
const (
	CliTracerServiceName = "cliDemo"

	GrpcServerName = "grpcServer"
	GrpcClientName = "grpcClient"

	HttpServerName = "httpServer"
	HttpClientName = "httpClient"
)
