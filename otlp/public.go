package otlp_go

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func InitTracer(appName string, config TracerConfig) (trace.Tracer, error) {
	return InitTraceProvider(appName, config)
}

func InitLocalLogger(config LoggerConfig) *slog.Logger {
	return InitLocalLoggerProvider(config)
}

func InitRemoteLogger(config LoggerConfig) *slog.Logger {
	return InitLokiLoggerProvider(config)
}
