package generator

import (
	"context"
	"errors"
	"fmt"
	"otelij/config"
)

type Generator interface {
	GenerateData(ctx context.Context) (any, error)
}

func NewGenerator(otlpType config.Type) (Generator, error) {
	if otlpType == config.TRACES {
		return SpanGenerator{}, nil
	}

	return nil, errors.New(fmt.Sprintf("Type %s is not yet supported.", otlpType))
}
