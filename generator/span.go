package generator

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net/url"
	"otelij/config"
	"strings"
)

type SpanGenerator struct{}

func (se SpanGenerator) GenerateData(ctx context.Context) (any, error) {

	status, err := config.GetSpanStatus()
	if err != nil {
		return nil, err
	}

	attributes, err := convertMapToAttribute(config.GetSpanAttributes())
	if err != nil {
		return nil, err
	}

	traceResource, err := resource.New(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(trace.WithResource(traceResource))

	tracer := tracerProvider.Tracer("Test Trace")

	_, span := tracer.Start(ctx, config.GetSpanName(), oteltrace.WithSpanKind(getSpanKind(config.GetSpanKind())))

	span.SetStatus(codes.Code(status), "")
	span.SetAttributes(attributes...)
	span.End()

	return span, nil
}

// Copied from /go.opentelemetry.io/otel/sdk/resource/env.go with minor modifications
func convertMapToAttribute(s string) ([]attribute.KeyValue, error) {
	if s == "" {
		return []attribute.KeyValue{}, nil
	}
	pairs := strings.Split(s, ",")
	attrs := []attribute.KeyValue{}
	var invalid []string
	for _, p := range pairs {
		field := strings.SplitN(p, "=", 2)
		if len(field) != 2 {
			invalid = append(invalid, p)
			continue
		}
		k := strings.TrimSpace(field[0])
		v, err := url.QueryUnescape(strings.TrimSpace(field[1]))
		if err != nil {
			// Retain original value if decoding fails, otherwise it will be
			// an empty string.
			v = field[1]
			otel.Handle(err)
		}
		attrs = append(attrs, attribute.String(k, v))
	}

	return attrs, nil
}

func getSpanKind(kind string) oteltrace.SpanKind {
	kindList := []oteltrace.SpanKind{oteltrace.SpanKindUnspecified, oteltrace.SpanKindInternal,
		oteltrace.SpanKindServer, oteltrace.SpanKindClient, oteltrace.SpanKindProducer, oteltrace.SpanKindConsumer}

	for _, option := range kindList {
		if option.String() == kind {
			return option
		}
	}

	return oteltrace.SpanKindUnspecified
}
