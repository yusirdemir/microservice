package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type TimeoutHandler struct{}

func NewTimeoutHandler() *TimeoutHandler {
	return &TimeoutHandler{}
}

func (h *TimeoutHandler) Register(r fiber.Router) {
	r.Get("/timeout", h.TestTimeout)
	r.Post("/timeout", h.TestTimeout)
}

func (h *TimeoutHandler) TestTimeout(c *fiber.Ctx) error {
	sleepStr := c.Query("sleep", "0s")
	sleep, err := time.ParseDuration(sleepStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid sleep duration")
	}

	if sleep > 0 {
		select {
		case <-time.After(sleep):
		case <-c.UserContext().Done():
			return c.Status(fiber.StatusRequestTimeout).SendString("Context Deadline Exceeded")
		}
	}

	if len(c.Body()) > 0 {
		return c.SendString("Received body with size: " + string(rune(len(c.Body()))))
	}

	return c.SendString("Finished processing after " + sleep.String())
}
