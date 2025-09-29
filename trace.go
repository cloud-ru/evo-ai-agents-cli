package main

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initTracing(ctx context.Context, config *config.Config) func() {
	//// Create OTLP exporter
	//exporter, err := stdouttrace
	//if err != nil {
	//	log.Printf("Failed to create exporter: %v", err)
	//	return func() {}
	//}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceConfig.ServiceName),
			semconv.ServiceVersion(config.ServiceConfig.Version),
			semconv.DeploymentEnvironment(string(config.ServiceConfig.AppEnvironment)),
		),
	)
	if err != nil {
		log.Printf("Failed to create resource: %v", err)
		return func() {}
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)

	// Set global trace provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}
