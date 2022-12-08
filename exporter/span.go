package exporter

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"go.opentelemetry.io/otel/sdk/trace"
	"otelij/config"
)

type SpanExporter struct {
	internalExporter trace.SpanExporter
}

func (se *SpanExporter) Export(ctx context.Context, span any) error {
	readonlySpan := span.(trace.ReadOnlySpan)
	return se.internalExporter.ExportSpans(ctx, []trace.ReadOnlySpan{readonlySpan})
}

func newSpanExporter(ctx context.Context) (*SpanExporter, error) {
	exporterType, err := config.GetExporter()
	if err != nil {
		return nil, err
	}

	var exporter trace.SpanExporter

	if exporterType == config.OTLP {
		exporter, err = newOtlpExporter(ctx)
		if err != nil {
			return nil, err
		}
	} else if exporterType == config.JAEGER {
		exporter, err = newJaegerExporter()
		if err != nil {
			return nil, err
		}
	} else if exporterType == config.ZIPKIN {
		exporter, err = zipkin.New("")
		if err != nil {
			return nil, err
		}
	} else if exporterType == config.STDOUT {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	} else {
		//Shouldn't happen because config already has a validation, but just in case
		return nil, errors.New(fmt.Sprintf("No such exporter: %s", exporterType))
	}

	return &SpanExporter{exporter}, nil
}

func newOtlpExporter(ctx context.Context) (trace.SpanExporter, error) {
	protocol, err := config.GetOtlpProtocol()
	if err != nil {
		return nil, err
	}

	var client otlptrace.Client
	if protocol == config.HttpProtobuf || protocol == config.HttpJson {
		client = otlptracehttp.NewClient()
	} else if protocol == config.GRPC {
		client = otlptracegrpc.NewClient()
	}

	return otlptrace.New(ctx, client)
}

func newJaegerExporter() (trace.SpanExporter, error) {
	protocol, err := config.GetJaegerProtocol()
	if err != nil {
		return nil, err
	}

	var endpointOption jaeger.EndpointOption

	if protocol == config.HttpThriftBinary {
		endpointOption = jaeger.WithCollectorEndpoint()
	} else if protocol == config.UdpThriftCompact {
		endpointOption = jaeger.WithAgentEndpoint()
	}

	return jaeger.New(endpointOption)
}
