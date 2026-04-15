package otlp_go

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func InitLocalLogger(config OtlpConfig) *slog.Logger {
	return InitLocalLoggerProvider(config)
}

func InitRemoteLogger(config OtlpConfig) *slog.Logger {
	return InitLokiLoggerProvider(config)
}

func InitTracer(config OtlpConfig) (trace.Tracer, error) {
	return InitTraceProvider(config)
}
