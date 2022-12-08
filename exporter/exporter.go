package exporter

import (
	"context"
	"errors"
	"fmt"
	"otelij/config"
)

type Exporter interface {
	Export(ctx context.Context, span any) error
}

func NewExporter(ctx context.Context, otlpType config.Type) (Exporter, error) {
	if otlpType == config.TRACES {
		return newSpanExporter(ctx)
	}

	return nil, errors.New(fmt.Sprintf("Type %s is not yet supported.", otlpType))
}
