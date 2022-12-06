# otelij
Open telemetry data generator (and exporter), currently only traces (spans) are supported

## How to run 
You can run the app with docker as following 
```
docker run -e OTEL_TRACES_EXPORTER=otlp -e OTEL_EXPORTER_OTLP_PROTOCOL=grpc somedocker/docker:tag
```
> **Note**
> 
> In case you want to send spans to localhost, you might to change it to ```host.docker.internal``` instead to make it work inside docker

## Exporting spans
This exporter is using the same exporter env specification as described here: https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/. 
In order to control the endpoint/protocol/headers/resources and other built in exporter stuff, you can just use the same env as you see in the specification.

### Usage
In order to use this app there are some env that you might need to pass: 
1. ```OTEL_TRACES_EXPORTER``` **(mandatory)** - valid options are [otlp, jaeger, zipkin, stdout]
2. ```OTEL_EXPORTER_OTLP_PROTOCOL``` (optional) - valid options are [grpc, http/protobuf, http/json], default is http/json
3. ```OTEL_EXPORTER_JAEGER_PROTOCOL``` (optional [relevant for jaeger exporter]) - valid options are [http/thrift.binary, udp/thrift.compact] - default is http/thrift.binary
4. Dynamic span data (all optional):
   1. ```SPAN_ATTRIBUTES``` - key values pair - comma delimited, e.g: ```attribute1=value1,attribute2=value2```
   2. ```SPAN_NAME``` - span name, regular string
   3. ```SPAN_KIND``` - span kind, valid options are: [internal, server, client, producer, consumer]
   4. ```SPAN_STATUS``` - number (between 0-2) stating the span status

### Examples
#### otlp 
otlp exporter full environments variables - https://opentelemetry.io/docs/reference/specification/protocol/exporter/

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
 -e SPAN_NAME=TestSpan \
 -e SPAN_KIND=server \
 -e SPAN_STATUS=1 \
 -e SPAN_ATTRIBUTES=span.attr\=val1,span.attr2\=val2 \
 somedocker/docker:tag
```

#### jaeger 
