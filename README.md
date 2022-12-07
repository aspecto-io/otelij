# otelij
Open telemetry data generator (and exporter), currently only traces (spans) are supported

## How to run 
You can run the app with docker as following 
```
docker run -e OTEL_TRACES_EXPORTER=otlp -e OTEL_EXPORTER_OTLP_PROTOCOL=grpc somedocker/docker:tag
```
> **Note**
> 
> In case you want to send spans to localhost, you can solve it with 2 options: 
> 1. Add ```--network host``` to docker command, eg: ```docker run --network host somedocker/docker:tag```
> 2. Replace ```localhost``` (in the relevant endpoint) with ```host.docker.internal``` instead, eg: ```docker run -e OTEL_EXPORTER_OTLP_ENDPOINT==https://host.docker.internal:4317 somedocker/docker:tag```

## Exporting spans
This exporter is using the same exporter env specification as described here: https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/. 
In order to control the endpoint/protocol/headers/resources and other built in exporter stuff, you can just use the same env as you see in the specification.

### Usage
In order to use this tool there are some env that you might need to pass: 
1. ```OTEL_TRACES_EXPORTER``` **(optional)** - valid options are [otlp, jaeger, zipkin, stdout], default is otlp
2. ```OTEL_EXPORTER_OTLP_PROTOCOL``` (optional) - valid options are [grpc, http/protobuf, http/json], default is grpc
3. ```OTEL_EXPORTER_JAEGER_PROTOCOL``` (optional [relevant for jaeger exporter]) - valid options are [http/thrift.binary, udp/thrift.compact] - default is http/thrift.binary
4. Dynamic span data (all optional):
   1. ```OTEL_SPAN_NAME``` - span name, regular string. default will be ```Otelij debug span```
   2. ```OTEL_SPAN_ATTRIBUTES``` - key values pair - comma delimited, e.g: ```attribute1=value1,attribute2=value2```. default is no attributes
   3. ```OTEL_SPAN_KIND``` - span kind, valid options are: [internal, server, client, producer, consumer]. default is internal
   4. ```OTEL_SPAN_STATUS``` - span status, valid options are: [Unset, Error, Ok]. default is Unset.
   5. ```OTEL_SPAN_STATUS_MESSAGE``` - free text for description of span status. default is empty.
   6. ```OTEL_SPAN_DURATION_SEC``` - span duration to set in seconds. default is 1s.
   7. ```OTEL_SPAN_LINK_TRACE_ID``` - in case you want to link this span to another trace. default is none.
   8. ```OTEL_SPAN_LINK_SPAN_ID``` - span id to link to. default is none.
   9. ```OTEL_SPAN_LINK_TRACE_FLAGS``` - byte (currently only 1 bit to represent sampled/not sampled). default is 1.
   10. ```OTEL_SPAN_LINK_REMOTE``` - boolean, if link propagated from a remote parent. 
   11. ```OTEL_SPAN_LINK_ATTRIBUTES``` - link attributes (same structure as OTEL_SPAN_ATTRIBUTES)

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
 somedocker/docker:tag
```

otlp with http (can be protobuf or json)
```bash
docker run -e OTEL_TRACES_EXPORTER=otlp \
 -e OTEL_EXPORTER_OTLP_ENDPOINT=https://my-endpoint.io:4318 \
 -e OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf \
 -e OTEL_EXPORTER_OTLP_HEADERS=Authorization\=SOME_TOKEN \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=MyServiceName \
 somedocker/docker:tag
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
 -e OTEL_SPAN_STATUS=1 \
 -e OTEL_SPAN_ATTRIBUTES=span.attr\=val1,span.attr2\=val2 \
 somedocker/docker:tag
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
 somedocker/docker:tag
```

with http/thrift.binary
```bash
docker run -e OTEL_TRACES_EXPORTER=jaeger \
 -e OTEL_EXPORTER_JAEGER_PROTOCOL=http/thrift.binary \
 -e OTEL_EXPORTER_JAEGER_ENDPOINT=http://localhost:14268/api/traces \
 -e OTEL_EXPORTER_JAEGER_AGENT_PORT=6831 \
 -e OTEL_RESOURCE_ATTRIBUTES=attribute1\=value1,attribute2\=value2 \
 -e OTEL_SERVICE_NAME=JaegerServiceName \
 somedocker/docker:tag
```

#### zipkin [env specification](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#zipkin-exporter)

```bash
docker run -e OTEL_TRACES_EXPORTER=zipkin \
 -e OTEL_EXPORTER_ZIPKIN_ENDPOINT=http://localhost:9411/api/v2/spans \
 -e OTEL_SERVICE_NAME=JaegerServiceName \
 somedocker/docker:tag
```
