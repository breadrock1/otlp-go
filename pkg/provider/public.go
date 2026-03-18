package provider

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func InitTracer(config OtlpConfig) (trace.Tracer, error) {
	return InitTraceProvider(config.Tracer)
}

func InitLocalLogger(config OtlpConfig) *slog.Logger {
	return InitLocalLoggerProvider(config)
}

func InitRemoteLogger(config OtlpConfig) *slog.Logger {
	return InitLokiLoggerProvider(config)
}
