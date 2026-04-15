package simple

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/breadrock1/otlp-go/otlp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.opentelemetry.io/otel/trace"

	otlppfiber "github.com/breadrock1/otlp-go/pkg/fiber"
)

const AppName = "simple-app"

type Server struct {
	tracer trace.Tracer

	Server *fiber.App
}

func SetupServer(otlpConfig otlp_go.OtlpConfig) *Server {
	tracer, err := otlp_go.InitTracer(otlpConfig)
	if err != nil {
		slog.Warn("failed to init tracer", slog.String("err", err.Error()))
	}

	serverApp := &Server{
		tracer: tracer,
	}

	serverApp.Server = fiber.New()

	serverApp.Server.Use(cors.New(cors.Config{}))
	serverApp.Server.Use(recover.New())

	serverApp.Server.Use(otlppfiber.PrometheusMeterMiddleware(serverApp.Server, otlpConfig))
	serverApp.Server.Use(otlppfiber.OtlpJaegerTracerMiddleware())
	serverApp.Server.Use(otlppfiber.StdoutLoggerMiddleware(otlpConfig))
	serverApp.Server.Use(otlppfiber.RemoteLokiLoggerMiddleware(otlpConfig))

	serverApp.Server.Get("/monitor", monitor.New())

	api := serverApp.Server.Group("/api")

	v1 := api.Group("/v1")
	serverApp.CreateGroup(v1)

	return serverApp
}

func (s *Server) Start(_ context.Context, address string) error {
	if err := s.Server.Listen(address); err != nil {
		return fmt.Errorf("failed to start Server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(_ context.Context) error {
	return s.Server.Shutdown()
}

//nolint
func main() {
	otlpConfig := otlp_go.OtlpConfig{
		Logger: otlp_go.LoggerConfig{
			Level:      "debug",
			Address:    "http://loki:3100",
			EnableLoki: true,
		},
		Tracer: otlp_go.TracerConfig{
			Address:      "http://jaeger:4317",
			EnableJaeger: true,
		},
	}

	ctx := context.Background()
	cCtx, cancel := context.WithCancel(ctx)
	httpServer := SetupServer(otlpConfig)
	go func() {
		if err := httpServer.Start(cCtx, "0.0.0.0:8080"); err != nil {
			slog.Error("http server start failed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cancel()
}
