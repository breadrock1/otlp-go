package fiber

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/breadrock1/otlp-go/otlp"
	"github.com/gofiber/fiber/v2"
)

const (
	defaultMetricsEndpoint = "/metrics"
)

func PrometheusMeterMiddleware(app *fiber.App, config otlp_go.OtlpConfig) fiber.Handler {
	prometheus := fiberprometheus.New(config.AppName)

	excludedPaths := otlp_go.AppendExcludedPath(config.Meter.ExcludedPaths)
	prometheus.SetSkipPaths(excludedPaths)

	endpointURI := defaultMetricsEndpoint
	if config.Meter.EndpointURI != "" {
		endpointURI = config.Meter.EndpointURI
	}
	prometheus.RegisterAt(app, endpointURI)

	return prometheus.Middleware
}
