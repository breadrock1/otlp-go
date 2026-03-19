package otlp_go

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func InitTracer(config TracerConfig) (trace.Tracer, error) {
	return InitTraceProvider(config)
}

func InitLocalLogger(config LoggerConfig) *slog.Logger {
	return InitLocalLoggerProvider(config)
}

func InitRemoteLogger(config LoggerConfig) *slog.Logger {
	return InitLokiLoggerProvider(config)
}
