package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yusirdemir/microservice/internal/transport/http/echo/handler"
)

type Router struct {
	app      *echo.Echo
	handlers []RouteHandler
}

type RouteHandler interface{}

func New(app *echo.Echo, handlers []RouteHandler) *Router {
	return &Router{
		app:      app,
		handlers: handlers,
	}
}

func (r *Router) SetupRoutes() {
	for _, h := range r.handlers {
		switch v := h.(type) {
		case *handler.UserHandler:
			r.app.POST("/users", v.CreateUser)
			r.app.GET("/users/:id", v.GetUser)
			r.app.PUT("/users/:id", v.UpdateUser)
			r.app.DELETE("/users/:id", v.DeleteUser)
		case *handler.HealthHandler:
			r.app.GET("/health/live", v.Live)
			r.app.GET("/health/ready", v.Ready)
		case *handler.TimeoutHandler:
			r.app.GET("/timeout", v.Timeout)
		}
	}
}
