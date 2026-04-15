package otlp_go

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

var (
	GlobalTracer    trace.Tracer
	TracePropagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
)

func InitTraceProvider(appName string, config TracerConfig) (trace.Tracer, error) {
	sampler := sdktrace.AlwaysSample()
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.TelemetrySDKLanguageGo,
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to merge trace resources: %w", err)
	}

	var traceOpts []sdktrace.TracerProviderOption
	traceOpts = append(traceOpts, sdktrace.WithResource(res))
	traceOpts = append(traceOpts, sdktrace.WithSampler(sampler))

	if config.EnableJaeger {
		ctx := context.Background()
		client := otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(config.Address),
			otlptracegrpc.WithInsecure(),
		)

		traceExporter, err := otlptrace.New(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to otlp trace server: %w", err)
		}
		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		traceOpts = append(traceOpts, sdktrace.WithSpanProcessor(bsp))
	}

	tp := sdktrace.NewTracerProvider(traceOpts...)
	otel.SetTextMapPropagator(TracePropagator)
	GlobalTracer = tp.Tracer(appName)
	otel.SetTracerProvider(tp)
	return GlobalTracer, nil
}
