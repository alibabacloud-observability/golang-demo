# Tracing Analysis Jaeger Demo

## 使用指南

#### SDK 上报

1. 配置 `examples/settings.go` 中的 `TracingAnalysisEndpoint`, 可以从 `https://tracing.console.aliyun.com/` => `集群配置`
   => `接入点信息` 中获得 Jaeger 的接入点（通过HTTP上报）
2. 使用 Go MOD 整理依赖 `go mod tidy`
3. 运行 Demo `go run tracingdemo`

#### Agent 上报

1. 修改 `examples/settings.go` 中的 `AgentSwitch` 为 `true`
2. 使用 Go MOD 整理依赖 `go mod tidy`
3. 启动 `Jaeger Agent` ，如果您使用 `Docker` 可以执行以下命令。其中 `reporter.grpc.host-port` 和 `agent.tags`
   参数值可以从 `https://tracing.console.aliyun.com/` => `集群配置` => `接入点信息` 中 Jaeger 的接入点（通过 Agent 上报）获得

```shell
sudo docker run \
  --rm \
  -p5775:5775/udp \
  -p6831:6831/udp \
  -p6832:6832/udp \
  -p5778:5778/tcp \
  jaegertracing/jaeger-agent:1.23 \
  --reporter.grpc.host-port=XXXXXX \
  --agent.tags=Authentication=XXXXXX
```

4. 运行 Demo `go run tracingdemo`

## 目录结构

```shell
|- jaeger-demo
  |- helloworld  # grpc pb 文件
  |- cli_demo.go # 纯函数 Demo
  |- grpc_client.go # grpc client 的 demo
  |- grpc_runner.go # grpc demo 入口
  |- http_client.go # http client 的 demo
  |- http_runner.go # http demo 入口
  |- http_server.go # http server 的 demo
  |- settings.go # 相关配置
  |- utils.go # 常用函数
|- go.mod # GO MOD 文件
|- http.go # 主入口
```