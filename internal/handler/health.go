package handler

import (
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Register(router fiber.Router) {
	router.Get("/health", h.Check)
	router.Get("/error", h.Error)
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return c.SendString("System is running perfectly.")
}

func (h *HealthHandler) Error(c *fiber.Ctx) error {
	return c.SendStatus(500)
}
