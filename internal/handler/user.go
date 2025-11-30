package handler

import "github.com/gofiber/fiber/v2"

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Register(r fiber.Router) {
	r.Get("/users/:id", h.GetUser)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.SendString("User ID: " + id)
}
