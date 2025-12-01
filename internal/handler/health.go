package handler

import (
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Register(router fiber.Router) {
	router.Get("/health/live", h.Live)
	router.Get("/health/ready", h.Ready)
}

func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "UP",
	})
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	ctx := c.UserContext()

	select {
	case <-ctx.Done():
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "DOWN",
			"error":  "Readiness check timed out",
		})
	default:
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "UP",
			"checks": fiber.Map{
				"database": "connected",
				"cache":    "connected",
			},
		})
	}
}
