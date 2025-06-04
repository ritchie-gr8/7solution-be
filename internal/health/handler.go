package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ritchie-gr8/7solution-be/internal/config"
	"github.com/ritchie-gr8/7solution-be/pkg/response"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func NewMonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{cfg: cfg}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	return response.NewResponse(c).Success(fiber.StatusOK, map[string]string{
		"name":    h.cfg.App().Name(),
		"version": h.cfg.App().Version(),
		"status":  "ok",
	}).Response()
}
