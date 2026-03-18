package simple

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	HelloEndpoint = "/hello"
)

func (s *Server) CreateGroup(group fiber.Router) {
	group.Get(HelloEndpoint, adaptor.HTTPHandler(promhttp.Handler()))
}

func (s *Server) Hello(eCtx *fiber.Ctx) error {
	return eCtx.SendString("hello")
}
