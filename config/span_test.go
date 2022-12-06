package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	TestName = "TestName"
)

func TestGetExporter(t *testing.T) {
	validateExporter(t, OTLP)
	validateExporter(t, JAEGER)
	validateExporter(t, ZIPKIN)
	validateExporter(t, STDOUT)

	os.Setenv(OtelTracesExporter, TestName)
	_, err := GetExporter()

	assert.NotNilf(t, err, "Expected to get an error while trying to create exporter %s", TestName)
}

func TestGetOtlpProtocol(t *testing.T) {
	validateOtlpProtocol(t, HttpProtobuf)
	validateOtlpProtocol(t, HttpJson)
	validateOtlpProtocol(t, GRPC)

	os.Setenv(OtelExporterOtlpProtocol, TestName)
	_, err := GetOtlpProtocol()

	assert.NotNilf(t, err, "Expected to get an error while trying to create protocol %s", TestName)
}

func TestGetJaegerProtocol(t *testing.T) {
	validateJaegerProtocol(t, HttpThriftBinary)
	validateJaegerProtocol(t, UdpThriftCompact)

	os.Setenv(OtelExporterJaegerProtocol, TestName)
	_, err := GetOtlpProtocol()

	assert.NotNilf(t, err, "Expected to get an error while trying to create jaeger protocol %s", TestName)

}

func TestGetSpanName(t *testing.T) {
	os.Unsetenv(SpanName)
	name := GetSpanName()

	assert.Equal(t, name, "test-span", "Expected to get default span name when no value exists.")

	os.Setenv(SpanName, TestName)
	name = GetSpanName()

	assert.Equal(t, name, TestName, "Expected to get span name value from env.")
}

func validateExporter(t *testing.T, exporter Exporter) {
	os.Setenv(OtelTracesExporter, string(exporter))
	created, err := GetExporter()

	assert.Nilf(t, err, "Error while trying to get exporter: %s", err)
	assert.Equalf(t, exporter, created, "Expected exporter: %s, got: %s", exporter, created)
}

func validateOtlpProtocol(t *testing.T, protocol OtlpProtocol) {
	os.Setenv(OtelExporterOtlpProtocol, string(protocol))
	created, err := GetOtlpProtocol()

	assert.Nilf(t, err, "Error while trying to get otlp protocol: %s", err)
	assert.Equalf(t, protocol, created, "Expected protocol: %s, got: %s", protocol, created)
}

func validateJaegerProtocol(t *testing.T, protocol JaegerProtocol) {
	os.Setenv(OtelExporterJaegerProtocol, string(protocol))
	created, err := GetJaegerProtocol()

	assert.Nilf(t, err, "Error while trying to get jaeger protocol: %s", err)
	assert.Equalf(t, protocol, created, "Expected jaeger protocol: %s, got: %s", protocol, created)
}
