package fiber

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/breadrock1/otlp-go/otlp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	XRequestIDHeaderKey = "X-Request-ID"
	ContextRequestIDKey = "request_id"
)

func StdoutLoggerMiddleware(config otlp_go.OtlpConfig) fiber.Handler {
	logger := otlp_go.InitLocalLoggerProvider(config)
	return func(eCtx *fiber.Ctx) error {
		if checkFilteredURI(eCtx.Path()) {
			return eCtx.Next()
		}

		requestID := eCtx.Get(XRequestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.New().String()
			eCtx.Set(XRequestIDHeaderKey, requestID)
		}

		startTime := time.Now()

		//nolint
		ctx := context.WithValue(eCtx.UserContext(), ContextRequestIDKey, requestID)
		eCtx.SetUserContext(ctx)

		err := eCtx.Next()

		latency := time.Since(startTime)

		statusCode := eCtx.Response().StatusCode()
		if err != nil {
			//nolint
			if fiberErr, ok := err.(*fiber.Error); ok {
				statusCode = fiberErr.Code
			}
		}

		var responseMsg = "Ok"
		var level = slog.LevelInfo
		if statusCode >= 300 {
			level = slog.LevelError
			responseMsg = string(eCtx.Response().Body())
		}

		slogAttrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", eCtx.Method()),
			slog.String("uri", eCtx.OriginalURL()),
			slog.Int("status", statusCode),
			slog.String("message", responseMsg),
			slog.Int("bytes_received", len(eCtx.Request().Body())),
			slog.Int("bytes_sent", len(eCtx.Response().Body())),
			slog.Duration("latency", latency),
			slog.String("referer", eCtx.Get("Referer")),
			slog.String("client_ip", eCtx.IP()),
			slog.String("user_agent", eCtx.Get("User-Agent")),
		}

		for _, attribute := range extractContextLocals(eCtx, config.Logger.Attributes) {
			slogAttrs = append(slogAttrs, attribute)
		}

		logger.LogAttrs(ctx, level, "http-request", slogAttrs...)

		return err
	}
}

func SyslogLoggerMiddleware(config otlp_go.OtlpConfig) fiber.Handler {
	logger := otlp_go.InitSyslogLoggerProvider(config)
	return func(eCtx *fiber.Ctx) error {
		if checkFilteredURI(eCtx.Path()) {
			return eCtx.Next()
		}

		requestID := eCtx.Get(XRequestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.New().String()
			eCtx.Set(XRequestIDHeaderKey, requestID)
		}

		ctx := context.WithValue(eCtx.UserContext(), ContextRequestIDKey, requestID)
		eCtx.SetUserContext(ctx)

		err := eCtx.Next()

		statusCode := eCtx.Response().StatusCode()
		if err != nil {
			//nolint
			if fiberErr, ok := err.(*fiber.Error); ok {
				statusCode = fiberErr.Code
			}
		}

		var responseMsg = "Ok"
		var level = slog.LevelInfo
		if statusCode >= 300 {
			level = slog.LevelError
			responseMsg = string(eCtx.Response().Body())
		}

		slogAttrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", eCtx.Method()),
			slog.String("uri", eCtx.OriginalURL()),
			slog.Int("status", statusCode),
			slog.String("message", responseMsg),
			slog.String("client_ip", eCtx.IP()),
			slog.String("user_agent", eCtx.Get("User-Agent")),
		}

		for _, attribute := range extractContextLocals(eCtx, config.Logger.Attributes) {
			slogAttrs = append(slogAttrs, attribute)
		}

		logger.LogAttrs(ctx, level, "http-request", slogAttrs...)

		return err
	}
}

func RemoteLokiLoggerMiddleware(config otlp_go.OtlpConfig) fiber.Handler {
	logger := otlp_go.InitLokiLoggerProvider(config)
	return func(eCtx *fiber.Ctx) error {
		if checkFilteredURI(eCtx.Path()) {
			return eCtx.Next()
		}

		start := time.Now()

		err := eCtx.Next()
		if err != nil {
			return err
		}

		latency := time.Since(start)

		var responseMsg = "Ok"
		var logLevel = slog.LevelInfo
		if eCtx.Response().StatusCode() >= 300 {
			logLevel = slog.LevelError
			responseMsg = string(eCtx.Response().Body())
		}

		logMessage := map[string]interface{}{
			"message":    responseMsg,
			"latency":    latency.String(),
			"status":     eCtx.Response().StatusCode(),
			"method":     eCtx.Method(),
			"uri":        eCtx.Path(),
			"client_ip":  eCtx.IP(),
			"user_agent": eCtx.Request(),
		}

		for attrKey, attribute := range extractContextLocals(eCtx, config.Logger.Attributes) {
			logMessage[attrKey] = attribute.Value.String()
		}

		jsonMessage, _ := json.Marshal(logMessage)

		ctx, cancel := context.WithTimeout(eCtx.Context(), 5*time.Second)
		logger.Log(ctx, logLevel, string(jsonMessage))
		defer cancel()

		return err
	}
}

func extractContextLocals(eCtx *fiber.Ctx, locals []string) map[string]slog.Attr {
	attributes := make(map[string]slog.Attr)
	for _, localKey := range locals {
		val, ok := eCtx.Locals(localKey).(string)
		if !ok || val == "" {
			continue
		}

		attributes[localKey] = slog.String(localKey, val)
	}

	return attributes
}

func checkFilteredURI(urlPath string) bool {
	for _, filtered := range otlp_go.ExcludedPaths {
		if strings.HasPrefix(urlPath, filtered) {
			return true
		}
	}

	return false
}
