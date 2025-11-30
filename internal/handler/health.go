package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Register(r fiber.Router) {
	r.Get("/", h.Check)
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	time.Sleep(5 * time.Second)
	return c.SendString("System is running perfectly.")
}
