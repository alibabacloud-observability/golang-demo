module github.com/alibabacloud-observability/golang-demo/opentelemetry-demo/single

go 1.16

require (
	go.opentelemetry.io/otel v1.15.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.15.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.15.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.15.1
	go.opentelemetry.io/otel/sdk v1.15.1
	go.opentelemetry.io/otel/trace v1.15.1
	google.golang.org/grpc v1.55.0
)
