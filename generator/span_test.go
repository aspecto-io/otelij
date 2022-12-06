package generator

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"os"
	"otelij/config"
	"testing"
)

func TestSpanGeneratedCorrectly(t *testing.T) {
	os.Setenv(config.SpanKind, "server")
	os.Setenv(config.SpanName, "best-span-ever")
	os.Setenv(config.SpanAttributes, "first=one,second=two")
	os.Setenv(config.SpanStatus, "2")

	spanGen := &SpanGenerator{}
	data, err := spanGen.GenerateData(context.Background())

	assert.Nil(t, err, "Error shouldn't be returned.")
	span := data.(trace.ReadOnlySpan)

	assert.Equal(t, oteltrace.SpanKindServer, span.SpanKind(), "Span kind should be server.")
	assert.Equal(t, "best-span-ever", span.Name(), "Span name should be 'best-span-ever'.")
	assert.Equal(t, 2, int(span.Status().Code), "Span status should be 2.")

	assert.Equal(t, 2, len(span.Attributes()), "expected 2 span attributes.")
	assert.Equal(t, "first", string(span.Attributes()[0].Key))
	assert.Equal(t, "one", span.Attributes()[0].Value.AsString())
	assert.Equal(t, "second", string(span.Attributes()[1].Key))
	assert.Equal(t, "two", span.Attributes()[1].Value.AsString())

}
