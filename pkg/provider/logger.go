package provider

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Marlliton/slogpretty"

	slogloki "github.com/samber/slog-loki/v2"
)

func InitLocalLoggerProvider(otlpConfig OtlpConfig) *slog.Logger {
	slogPrettyOpts := &slogpretty.Options{
		Level: getLevelType(otlpConfig.Logger),
	}

	textHandler := slogpretty.New(os.Stdout, slogPrettyOpts)
	localLogger := slog.New(textHandler)
	return localLogger
}

func InitLokiLoggerProvider(otlpConfig OtlpConfig) *slog.Logger {
	level := getLevelType(otlpConfig.Logger)
	lokiConfig := slogloki.Option{
		Endpoint:           fmt.Sprintf("%s/api/prom/push", otlpConfig.Logger.Address),
		Level:              slog.LevelInfo,
		BatchWait:          time.Second * 5,
		BatchEntriesNumber: 10,
	}

	lokiHandler := lokiConfig.NewLokiHandler()

	logger := slog.New(lokiHandler).
		With("service_name", AppName).
		With("service", AppName).
		With("detected_level", level).
		With("level", level)

	return logger
}

func getLevelType(config LoggerConfig) slog.Level {
	var logLevel = slog.LevelInfo
	switch config.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return logLevel
}
