## 通过OpenTelemetry上报Go应用数据

- jaeger-exporter: 使用OpenTelemetry协议采集数据，并使用Jaeger Exporter上报数据

- oltp-exporter: 使用OpenTelemetry协议采集数据，并使用OTLP(OpenTelemetry Protocol Exporter) Exporter上报数据
  - server: 服务端demo（通过gRPC协议上报链路数据）
  - client_grpc_export: 客户端demo1（通过gRPC协议上报链路数据）
  - client_http_export: 客户端demo2（通过HTTP协议上报链路数据）
