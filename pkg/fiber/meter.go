package fiber

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/breadrock1/otlp-go/otlp"
	"github.com/gofiber/fiber/v2"
)

const (
	metricsEndpoint = "/metrics"
)

func PrometheusMeterMiddleware(app *fiber.App, config otlp_go.OtlpConfig) fiber.Handler {
	prometheus := fiberprometheus.New(config.AppName)
	prometheus.SetSkipPaths(otlp_go.ExcludedPaths)
	prometheus.RegisterAt(app, metricsEndpoint)
	return prometheus.Middleware
}
