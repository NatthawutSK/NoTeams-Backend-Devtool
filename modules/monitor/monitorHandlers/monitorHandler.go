package monitorHandlers

import (
	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/entities"
	"github.com/NatthawutSK/NoTeams-Backend/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	// get coookies
	// cookies := c.Cookies("accessToken")

	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
		// Token:   cookies,
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		res,
	).Res()
}
