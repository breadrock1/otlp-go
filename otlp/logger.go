package otlp_go

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Marlliton/slogpretty"

	slogloki "github.com/samber/slog-loki/v2"
)

func InitLocalLoggerProvider(config OtlpConfig) *slog.Logger {
	slogPrettyOpts := &slogpretty.Options{
		Level: getLevelType(config.Logger),
	}

	textHandler := slogpretty.New(os.Stdout, slogPrettyOpts)
	localLogger := slog.New(textHandler)
	return localLogger
}

func InitLokiLoggerProvider(config OtlpConfig) *slog.Logger {
	level := getLevelType(config.Logger)
	lokiConfig := slogloki.Option{
		Endpoint:           fmt.Sprintf("%s/api/prom/push", config.Logger.Address),
		Level:              slog.LevelInfo,
		BatchWait:          time.Second * 5,
		BatchEntriesNumber: 10,
	}

	lokiHandler := lokiConfig.NewLokiHandler()

	logger := slog.New(lokiHandler).
		With("service_name", config.AppName).
		With("service", config.AppName).
		With("detected_level", level).
		With("level", level)

	return logger
}

func getLevelType(config LoggerConfig) slog.Level {
	switch config.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
