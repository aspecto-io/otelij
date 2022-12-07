package generator

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net/url"
	"otelij/config"
	"strings"
	"time"
)

type SpanGenerator struct{}

type spanconf struct {
	name          string
	status        codes.Code
	kind          oteltrace.SpanKind
	statusMessage string
	attributes    []attribute.KeyValue
	duration      time.Duration
	link          *oteltrace.Link
}

func (se SpanGenerator) GenerateData(ctx context.Context) (any, error) {

	spanConfig, err := getSpanConfiguration()
	if err != nil {
		return nil, err
	}

	traceResource, err := resource.New(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(trace.WithResource(traceResource))

	tracer := tracerProvider.Tracer("Test Trace")

	startOptions := []oteltrace.SpanStartOption{oteltrace.WithSpanKind(spanConfig.kind), oteltrace.WithTimestamp(time.Now())}
	if spanConfig.link != nil {
		startOptions = append(startOptions, oteltrace.WithLinks(*spanConfig.link))
	}
	_, span := tracer.Start(ctx, config.GetSpanName(), startOptions...)

	span.SetStatus(spanConfig.status, spanConfig.statusMessage)
	span.SetAttributes(spanConfig.attributes...)

	span.End(oteltrace.WithTimestamp(time.Now().Add(time.Second * spanConfig.duration)))

	return span, nil
}

func getSpanConfiguration() (*spanconf, error) {

	var errs []string

	_ = config.GetServiceNameOrDefault()

	status := getSpanStatus()
	statusMessage := config.GetSpanStatusMessage()

	attributes, err := convertStringToAttribute(config.GetSpanAttributes())
	if err != nil {
		errs = append(errs, err.Error())
	}

	duration, err := config.GetSpanDuration()
	if err != nil {
		errs = append(errs, err.Error())
	}

	spanKind := getSpanKind()

	link, err := createLinkIfNeeded()
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	return &spanconf{name: config.GetSpanName(),
		status:        status,
		kind:          spanKind,
		statusMessage: statusMessage,
		attributes:    attributes,
		duration:      time.Duration(duration),
		link:          link,
	}, nil
}

func getSpanKind() oteltrace.SpanKind {
	kind := config.GetSpanKind()
	kindList := []oteltrace.SpanKind{oteltrace.SpanKindUnspecified, oteltrace.SpanKindInternal,
		oteltrace.SpanKindServer, oteltrace.SpanKindClient, oteltrace.SpanKindProducer, oteltrace.SpanKindConsumer}

	for _, option := range kindList {
		if strings.EqualFold(option.String(), kind) {
			return option
		}
	}

	return oteltrace.SpanKindUnspecified
}

func getSpanStatus() codes.Code {
	kind := config.GetSpanStatus()
	kindList := []codes.Code{codes.Unset, codes.Error, codes.Ok}

	for _, option := range kindList {
		if strings.EqualFold(option.String(), kind) {
			return option
		}
	}

	return codes.Unset
}

func createLinkIfNeeded() (*oteltrace.Link, error) {
	linkTraceId := config.GetLinkTraceId()
	linkSpanId := config.GetLinkSpanId()

	if linkTraceId == "" && linkSpanId == "" {
		return nil, nil
	}

	if linkTraceId == "" || linkSpanId == "" {
		return nil, errors.New("you must provide both link trace id and span id in order to create a link")
	}

	traceId := oteltrace.TraceID{}

	decodedTraceId, err := hex.DecodeString(linkTraceId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to decode linked trace id: %s, error: %s", linkTraceId, err.Error()))
	}

	for index, aByte := range decodedTraceId {
		traceId[index] = aByte
	}

	spanId := oteltrace.SpanID{}
	decodedSpanId, err := hex.DecodeString(linkSpanId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to decode linked span id: %s, error: %s", linkTraceId, err.Error()))
	}
	for index, aByte := range decodedSpanId {
		spanId[index] = aByte
	}

	flags, err := config.GetLinkFlags()
	if err != nil {
		return nil, err
	}

	remote, err := config.GetLinkRemote()
	if err != nil {
		return nil, err
	}

	attributes, err := convertStringToAttribute(config.GetLinkAttributes())
	if err != nil {
		return nil, err
	}

	return &oteltrace.Link{SpanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		TraceID:    traceId,
		SpanID:     spanId,
		TraceFlags: oteltrace.TraceFlags(flags),
		Remote:     remote,
	}),
		Attributes: attributes}, nil
}

// Copied from /go.opentelemetry.io/otel/sdk/resource/env.go with minor modifications
func convertStringToAttribute(s string) ([]attribute.KeyValue, error) {
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
