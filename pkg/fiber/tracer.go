package fiber

import (
	"strings"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"

	"otlp-go/pkg/provider"
)

func OtlpJaegerTracerMiddleware() fiber.Handler {
	return otelfiber.Middleware(
		otelfiber.WithNext(traceURLSkipper),
	)
}

func traceURLSkipper(eCtx *fiber.Ctx) bool {
	for _, excluded := range provider.ExcludedPaths {
		if strings.HasPrefix(eCtx.Path(), excluded) {
			return true
		}
	}

	return eCtx.Request().Header.IsOptions()
}
