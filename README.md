# otelij
Open telemetry data generator (and exporter), currently only traces (spans) are supported

## How to run 
You can run the tool with docker as following
```
docker run -e OTEL_TRACES_EXPORTER=otlp -e OTEL_EXPORTER_OTLP_PROTOCOL=grpc public.ecr.aws/x3s3n8k7/otelij
```
This will start the otel debugger, send a single span with the configured exporter, and terminate when exporting is settled

> **Note**
> 
> In case you want to send spans to localhost, you can solve it with 2 options: 
> 1. Add ```--network host``` to docker command, e.g., ```docker run --network host public.ecr.aws/x3s3n8k7/otelij```
> 2. Replace ```localhost``` (in the relevant endpoint) with ```host.docker.internal``` instead, e.g., ```docker run -e OTEL_EXPORTER_OTLP_ENDPOINT==https://host.docker.internal:4317 public.ecr.aws/x3s3n8k7/otelij```

## Exporting spans
Otelij exporter is using the same exporter environment variables specification as described in [OpenTelemetry specification](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables).
In order to control the endpoint/protocol/headers/resources and other built in exporter stuff, you can just use the same env as you see in the specification.

### Usage
In order to use this tool there are some environment variables that you might need to set: 
1. ```OTEL_TRACES_EXPORTER``` (optional) - valid options are [otlp, jaeger, zipkin, stdout], default is otlp
2. ```OTEL_EXPORTER_OTLP_PROTOCOL``` (optional [relevant for otlp exporter]) - valid options are [grpc, http/protobuf, http/json], default is grpc
3. ```OTEL_EXPORTER_JAEGER_PROTOCOL``` (optional [relevant for jaeger exporter]) - valid options are [http/thrift.binary, udp/thrift.compact] - default is http/thrift.binary
4. Dynamic span data (all optional):
   1. ```OTEL_SPAN_NAME``` - span name, regular string. default will be ```Otelij debug span```
   2. ```OTEL_SPAN_ATTRIBUTES``` - key values pair - comma delimited, e.g: ```attribute1=value1,attribute2=value2```. default is no attributes
   3. ```OTEL_SPAN_KIND``` - span kind, valid options are: [internal, server, client, producer, consumer]. default is internal
   4. ```OTEL_SPAN_STATUS``` - span status, valid options are: [Unset, Error, Ok]. default is Unset.
   5. ```OTEL_SPAN_STATUS_MESSAGE``` - free text for description of span status. default is empty.
   6. ```OTEL_SPAN_DURATION_SEC``` - span duration to set in seconds. default is 1, usage as following OTEL_SPAN_DURATION_SEC=5.
5. Add link to another trace (optional)
   1. ```OTEL_SPAN_LINK_TRACE_ID``` **(mandatory)** - in case you want to link this span to another trace. default is none.
   2. ```OTEL_SPAN_LINK_SPAN_ID``` **(mandatory)** - span id to link to. default is none.
   3. ```OTEL_SPAN_LINK_TRACE_FLAGS``` - byte (currently only 1 bit to represent sampled/not sampled). default is 1.
   4. ```OTEL_SPAN_LINK_REMOTE``` - boolean, if link propagated from a remote parent. 
   5. ```OTEL_SPAN_LINK_ATTRIBUTES``` - link attributes (same structure as OTEL_SPAN_ATTRIBUTES)

### Examples
#### otlp [env specification](https://opentelemetry.io/docs/reference/specification/protocol/exporter/)

otlp with grpc
```bash
docker run -e OTEL_TRACES_EXPORTER=otlp \
 -e OTEL_EXPORTER_OTLP_ENDPOINT=https://my-endpoint.io:4317 \
 -e OTEL_EXPORTER_OTLP_PROTOCOL=grpc \
 -e OTEL_EXPORTER_OTLP_HEADERS=Authorization\=SOME_TOKEN \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=MyServiceName \
 public.ecr.aws/x3s3n8k7/otelij
```

otlp with http/protobuf
```bash
docker run -e OTEL_TRACES_EXPORTER=otlp \
 -e OTEL_EXPORTER_OTLP_ENDPOINT=https://my-endpoint.io:4318 \
 -e OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf \
 -e OTEL_EXPORTER_OTLP_HEADERS=Authorization\=SOME_TOKEN \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=MyServiceName \
 public.ecr.aws/x3s3n8k7/otelij
```

with span data 
```bash
docker run -e OTEL_TRACES_EXPORTER=otlp \
 -e OTEL_EXPORTER_OTLP_ENDPOINT=https://my-endpoint.io:4318 \
 -e OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf \
 -e OTEL_EXPORTER_OTLP_HEADERS=Authorization\=SOME_TOKEN \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=MyServiceName \
 -e OTEL_SPAN_NAME=TestSpan \
 -e OTEL_SPAN_KIND=server \
 -e OTEL_SPAN_STATUS=OK \
 -e OTEL_SPAN_ATTRIBUTES=span.attr\=val1,span.attr2\=val2 \
 public.ecr.aws/x3s3n8k7/otelij
```

#### jaeger [env specification](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#jaeger-exporter)
with udp/thrift.compact
```bash
docker run -e OTEL_TRACES_EXPORTER=jaeger \
 -e OTEL_EXPORTER_JAEGER_PROTOCOL=udp/thrift.compact \
 -e OTEL_EXPORTER_JAEGER_AGENT_HOST=localhost \
 -e OTEL_EXPORTER_JAEGER_AGENT_PORT=6831 \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=JaegerServiceName \
 public.ecr.aws/x3s3n8k7/otelij
```

with http/thrift.binary
```bash
docker run -e OTEL_TRACES_EXPORTER=jaeger \
 -e OTEL_EXPORTER_JAEGER_PROTOCOL=http/thrift.binary \
 -e OTEL_EXPORTER_JAEGER_ENDPOINT=http://localhost:14268/api/traces \
 -e OTEL_EXPORTER_JAEGER_AGENT_PORT=6831 \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=JaegerServiceName \
 public.ecr.aws/x3s3n8k7/otelij
```

#### zipkin [env specification](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#zipkin-exporter)

```bash
docker run -e OTEL_TRACES_EXPORTER=zipkin \
 -e OTEL_EXPORTER_ZIPKIN_ENDPOINT=http://localhost:9411/api/v2/spans \
 -e OTEL_SERVICE_NAME=JaegerServiceName \
 public.ecr.aws/x3s3n8k7/otelij
```


#### with link 
```bash
docker run -e OTEL_TRACES_EXPORTER=otlp \
 -e OTEL_EXPORTER_OTLP_ENDPOINT=https://my-endpoint.io:4318 \
 -e OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf \
 -e OTEL_EXPORTER_OTLP_HEADERS=Authorization\=SOME_TOKEN \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=MyServiceName \
 -e OTEL_SPAN_NAME=TestSpan \
 -e OTEL_SPAN_KIND=server \
 -e OTEL_SPAN_STATUS=OK \
 -e OTEL_SPAN_ATTRIBUTES=span.attr\=val1,span.attr2\=val2 \
 -e OTEL_SPAN_LINK_SPAN_ID=d6583451bafe66cb \
 -e OTEL_SPAN_LINK_TRACE_ID=d18ae83289fb43df3f8570bcb5c3177c \
 public.ecr.aws/x3s3n8k7/otelij
```
