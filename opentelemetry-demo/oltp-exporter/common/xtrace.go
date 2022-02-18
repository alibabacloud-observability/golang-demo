package common

import (
	"log"
	"os"
)

const OtlpEndpointEnvName = "OTEL_EXPORTER_OTLP_ENDPOINT"
const XTraceTokenEnvName = "XTRACE_TOKEN"

// ObtainXTraceInfo 从环境变量获取 XTrace 配置信息
func ObtainXTraceInfo() (string, string, bool) {
	otelAgentAddr, ok := os.LookupEnv(OtlpEndpointEnvName)
	if !ok {
		log.Printf("invalid otlp endpiont from os.env: %s\n", OtlpEndpointEnvName)
		return "", "", false
	}
	xtraceToken, ok := os.LookupEnv(XTraceTokenEnvName)
	if !ok {
		log.Printf("invalid xtrace token from os.env: %s\n", XTraceTokenEnvName)
		return "", "", false
	}
	log.Printf("OTLP Endpoint: %s\t", otelAgentAddr)
	log.Printf("XTRACE Token: %s\n", xtraceToken)

	return otelAgentAddr, xtraceToken, true
}
