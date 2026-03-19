package fiber

import (
	"strings"

	"github.com/breadrock1/otlp-go/otlp"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

func OtlpJaegerTracerMiddleware() fiber.Handler {
	return otelfiber.Middleware(
		otelfiber.WithNext(traceURLSkipper),
	)
}

func traceURLSkipper(eCtx *fiber.Ctx) bool {
	for _, excluded := range otlp_go.ExcludedPaths {
		if strings.HasPrefix(eCtx.Path(), excluded) {
			return true
		}
	}

	return eCtx.Request().Header.IsOptions()
}
