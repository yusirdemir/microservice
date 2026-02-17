package router

import (
	"github.com/gofiber/fiber/v2"
)

type RouteHandler interface {
	Register(r fiber.Router)
}

type Router struct {
	App      *fiber.App
	Handlers []RouteHandler
}

func New(app *fiber.App, handlers []RouteHandler) *Router {
	return &Router{
		App:      app,
		Handlers: handlers,
	}
}

func (r *Router) SetupRoutes() {
	for _, h := range r.Handlers {
		h.Register(r.App)
	}
}
