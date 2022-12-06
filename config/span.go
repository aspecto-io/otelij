package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	OtelTracesExporter       = "OTEL_TRACES_EXPORTER"        //otlp, jaeger, zipkin, stdout
	OtelExporterOtlpProtocol = "OTEL_EXPORTER_OTLP_PROTOCOL" //grpc, http/protobuf, http/json

	OtelExporterJaegerProtocol = "OTEL_EXPORTER_JAEGER_PROTOCOL" //http/thrift.binary, grpc, udp/thrift.compact, udp/thrift.binary

	SpanAttributes = "SPAN_ATTRIBUTES"
	SpanName       = "SPAN_NAME"
	SpanKind       = "SPAN_KIND"
	SpanStatus     = "SPAN_STATUS"
)

func GetExporter() (Exporter, error) {
	exporter, ok := os.LookupEnv(OtelTracesExporter)

	if !ok {
		return "", ParameterError{fmt.Sprintf("%s must be provided. options are: %s", OtelTracesExporter, getExporters())}
	}

	if !isValidExporter(exporter) {
		return "", ParameterError{fmt.Sprintf("Invalid exporter: %s. valid options are: %s", exporter, getExporters())}
	}

	return Exporter(exporter), nil
}

func GetOtlpProtocol() (OtlpProtocol, error) {
	protocol, ok := os.LookupEnv(OtelExporterOtlpProtocol)

	if ok {
		if !isValidOtlpProtocol(protocol) {
			return "", ParameterError{fmt.Sprintf("Invalid otlp protocol: %s. valid options are: %s", protocol, getOtlpProtocols())}
		}
		return OtlpProtocol(protocol), nil
	}

	return HttpJson, nil
}

func GetJaegerProtocol() (JaegerProtocol, error) {
	protocol, ok := os.LookupEnv(OtelExporterJaegerProtocol)

	if ok {
		if !isValidJaegerProtocol(protocol) {
			return "", ParameterError{fmt.Sprintf("Invalid jaeger protocol: %s. valid options are: %s", protocol, getJaegerProtocols())}
		}
		return JaegerProtocol(protocol), nil
	}

	return HttpThriftBinary, nil
}

func GetSpanAttributes() string {
	return os.Getenv(SpanAttributes)
}

func GetSpanName() string {
	name := os.Getenv(SpanName)

	if name == "" {
		return "test-span"
	}

	return name
}

func GetSpanStatus() (uint64, error) {
	status, ok := os.LookupEnv(SpanStatus)
	if !ok {
		return 0, nil
	}

	return strconv.ParseUint(status, 10, 32)
}

func GetSpanKind() string {
	return os.Getenv(SpanKind)
}
