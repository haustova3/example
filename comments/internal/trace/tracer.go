package trace

import (
	"context"
	"example/comments/internal/app/config"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracerProvider trace.TracerProvider
	tracer         trace.Tracer
)

func CreateTracerProvider(ctx context.Context, config *config.Config) {
	jaegerResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("loms_service"),
		),
	)
	if err != nil {
		panic(err)
	}
	jaegerAddress := fmt.Sprintf("http://%s:%s", config.JaegerConf.Host, config.JaegerConf.Port)
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(jaegerAddress))
	if err != nil {
		panic(err)
	}

	tracerProvider = trace2.NewTracerProvider(
		trace2.WithBatcher(exp),
		trace2.WithResource(jaegerResource),
	)
	otel.SetTracerProvider(tracerProvider)
	tracer = otel.GetTracerProvider().Tracer("cart_service")
}

func Tracer() trace.Tracer {
	return tracer
}
