package fiber

import (
	"otlp-go/pkg/provider"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

const (
	metricsEndpoint = "/metrics"
)

func PrometheusMeterMiddleware(app *fiber.App) fiber.Handler {
	prometheus := fiberprometheus.New(provider.AppName)
	prometheus.SetSkipPaths(provider.ExcludedPaths)
	prometheus.RegisterAt(app, metricsEndpoint)
	return prometheus.Middleware
}
