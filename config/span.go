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

	OtelServiceName = "OTEL_SERVICE_NAME"

	SpanAttributes     = "OTEL_SPAN_ATTRIBUTES"
	SpanName           = "OTEL_SPAN_NAME"
	SpanKind           = "OTEL_SPAN_KIND"
	SpanStatus         = "OTEL_SPAN_STATUS"
	SpanStatusMessage  = "OTEL_SPAN_STATUS_MESSAGE"
	SpanDuration       = "OTEL_SPAN_DURATION_SEC" //duration in seconds
	SpanLinkTraceId    = "OTEL_SPAN_LINK_TRACE_ID"
	SpanLinkSpanId     = "OTEL_SPAN_LINK_SPAN_ID"
	SpanLinkFlags      = "OTEL_SPAN_LINK_TRACE_FLAGS" //byte
	SpanLinkRemote     = "OTEL_SPAN_LINK_REMOTE"      //true, false
	SpanLinkAttributes = "OTEL_SPAN_LINK_ATTRIBUTES"  //same as span attributes
)

func GetExporter() (Exporter, error) {
	exporter, ok := os.LookupEnv(OtelTracesExporter)

	if ok {
		if !isValidExporter(exporter) {
			return "", ParameterError{fmt.Sprintf("Invalid exporter: %s. valid options are: %s", exporter, getExporters())}
		}
		return Exporter(exporter), nil
	}

	return OTLP, nil
}

func GetOtlpProtocol() (OtlpProtocol, error) {
	protocol, ok := os.LookupEnv(OtelExporterOtlpProtocol)

	if ok {
		if !isValidOtlpProtocol(protocol) {
			return "", ParameterError{fmt.Sprintf("Invalid otlp protocol: %s. valid options are: %s", protocol, getOtlpProtocols())}
		}
		return OtlpProtocol(protocol), nil
	}

	return GRPC, nil
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
		return "Otelij debug span"
	}

	return name
}

func GetSpanStatus() string {
	return os.Getenv(SpanStatus)
}

func GetSpanStatusMessage() string {
	return os.Getenv(SpanStatusMessage)
}

func GetSpanKind() string {
	return os.Getenv(SpanKind)
}

func GetSpanDuration() (uint64, error) {
	duration, ok := os.LookupEnv(SpanDuration)
	if ok {
		return strconv.ParseUint(duration, 10, 32)
	}

	return 1, nil
}

func GetServiceNameOrDefault() string {
	service, ok := os.LookupEnv(OtelServiceName)
	if !ok {
		name := "Otelij debugger"
		os.Setenv(OtelServiceName, "Otelij debugger")
		return name
	}

	return service
}

func GetLinkTraceId() string {
	return os.Getenv(SpanLinkTraceId)
}

func GetLinkSpanId() string {
	return os.Getenv(SpanLinkSpanId)
}

func GetLinkFlags() (uint64, error) {
	flags, ok := os.LookupEnv(SpanLinkFlags)

	if ok {
		return strconv.ParseUint(flags, 10, 8)
	}

	return 1, nil
}

func GetLinkRemote() (bool, error) {
	remote, ok := os.LookupEnv(SpanLinkRemote)

	if ok {
		return strconv.ParseBool(remote)
	}

	return true, nil
}

func GetLinkAttributes() string {
	return os.Getenv(SpanLinkAttributes)
}
