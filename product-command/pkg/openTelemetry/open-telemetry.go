package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTracerProvider(ctx context.Context, serviceName, serviceVersion, jaegerURL string) (*sdktrace.TracerProvider, error) {
	// 1. Define o recurso (Identificação do Serviço)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	// 2. Configura o exportador OTLP via gRPC para o Jaeger
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(jaegerURL), // ex: "localhost:4317" ou "jaeger:4317"
	)
	if err != nil {
		return nil, err
	}

	// 3. Cria o TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // Envia spans em lotes para alta performance
		sdktrace.WithResource(res),
	)

	// 4. Define o TracerProvider e o Propagador de Contexto Globais
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // Garante a propagação do Trace ID (Kafka, HTTP)
			propagation.Baggage{},
		),
	)

	return tp, nil
}
