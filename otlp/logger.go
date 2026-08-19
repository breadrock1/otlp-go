package otlp_go

import (
	"fmt"
	"log/slog"
	"log/syslog"
	"os"
	"time"

	"github.com/Marlliton/slogpretty"

	slogloki "github.com/samber/slog-loki/v2"
	slogsyslog "github.com/samber/slog-syslog/v2"
)

func InitLocalLoggerProvider(config OtlpConfig) *slog.Logger {
	slogPrettyOpts := &slogpretty.Options{
		Level: getSlogLevelType(config.Logger),
	}

	textHandler := slogpretty.New(os.Stdout, slogPrettyOpts)
	localLogger := slog.New(textHandler)
	return localLogger
}

func InitSyslogLoggerProvider(config OtlpConfig) *slog.Logger {
	level := getSlogLevelType(config.Logger)
	syslogLevel := getSyslogLevelType(config.Logger)

	var address, network string
	if config.Logger.Syslog.Address != "" {
		network = "udp"
		address = config.Logger.Syslog.Address
	}

	writer, err := syslog.Dial(network, address, syslogLevel, config.AppName)
	if err != nil {
		slog.Error("failed to init syslog handler", slog.String("err", err.Error()))
	}

	syslogOpts := slogsyslog.Option{
		Level:  slog.LevelDebug,
		Writer: writer,
	}

	logger := slog.New(syslogOpts.NewSyslogHandler()).
		With("service_name", config.AppName).
		With("service", config.AppName).
		With("detected_level", level).
		With("level", level)

	return logger
}

func InitLokiLoggerProvider(config OtlpConfig) *slog.Logger {
	level := getSlogLevelType(config.Logger)
	lokiConfig := slogloki.Option{
		Endpoint:           fmt.Sprintf("%s/api/prom/push", config.Logger.Loki.Address),
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

func getSyslogLevelType(config LoggerConfig) syslog.Priority {
	var level syslog.Priority
	switch config.Level {
	case "debug":
		level = syslog.LOG_DEBUG
	case "info":
		level = syslog.LOG_INFO
	case "warn":
		level = syslog.LOG_WARNING
	case "error":
		level = syslog.LOG_ERR
	default:
		level = syslog.LOG_INFO
	}

	return level | syslog.LOG_LOCAL0
}

func getSlogLevelType(config LoggerConfig) slog.Level {
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
