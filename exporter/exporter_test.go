package exporter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"otelij/config"
	"testing"
)

func TestNewExporterMetricsNotSupported(t *testing.T) {
	_, err := NewExporter(context.Background(), config.METRICS)

	assert.NotNil(t, err, "Expected metrics exporter to return error.")
	assert.Contains(t, err.Error(), "not yet supported")
}

func TestNewExporterLogsNotSupported(t *testing.T) {
	_, err := NewExporter(context.Background(), config.LOGS)

	assert.NotNil(t, err, "Expected logs exporter to return error.")
	assert.Contains(t, err.Error(), "not yet supported")
}

func TestTraceNewExporterWithoutOtlpProtocol(t *testing.T) {
	//Put empty value
	os.Unsetenv(config.OtelTracesExporter)

	_, err := NewExporter(context.Background(), config.TRACES)

	assert.NotNil(t, err, "Expected trace exporter to return error.")
	assert.Contains(t, err.Error(), "must be provided")
}

func TestTraceNewExporterWithInvalidProtocolProtocol(t *testing.T) {
	//Put empty value
	os.Setenv(config.OtelTracesExporter, "test")

	_, err := NewExporter(context.Background(), config.TRACES)

	assert.NotNil(t, err, "Expected trace exporter to return error.")
	assert.Contains(t, err.Error(), "valid options are:")
}

func TestTraceNewExporterValid(t *testing.T) {
	//Put empty value
	os.Setenv(config.OtelTracesExporter, string(config.OTLP))

	_, err := NewExporter(context.Background(), config.TRACES)

	assert.Nilf(t, err, "Expected trace exporter to not return error.")
}
