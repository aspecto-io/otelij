package main

import (
	"context"
	"fmt"
	"os"
	"otelij/config"
	"otelij/exporter"
	"otelij/generator"
)

func main() {

	ctx := context.Background()
	otlpType := config.GetType()
	handleRequest(ctx, otlpType)

	os.Exit(0)
}

func handleRequest(ctx context.Context, otlpType config.Type) {

	dataGenerator, err := generator.NewGenerator(otlpType)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error while trying to create generator: %s", err.Error()))
		return
	}

	dataExporter, err := exporter.NewExporter(ctx, otlpType)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error while trying to create exporter: %s", err.Error()))
		return
	}

	data, err := dataGenerator.GenerateData(ctx)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error while trying to generate data: %s", err.Error()))
		return
	}

	err = dataExporter.Export(ctx, data)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error while trying to export data: %s", err.Error()))
	}
}
