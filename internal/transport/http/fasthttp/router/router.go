package router

import (
	"github.com/fasthttp/router"
	"github.com/yusirdemir/microservice/internal/transport/http/fasthttp/handler"
)

type Router struct {
	router   *router.Router
	handlers []RouteHandler
}

type RouteHandler interface{}

func New(r *router.Router, handlers []RouteHandler) *Router {
	return &Router{
		router:   r,
		handlers: handlers,
	}
}

func (r *Router) SetupRoutes() {
	for _, h := range r.handlers {
		switch v := h.(type) {
		case *handler.UserHandler:
			r.router.POST("/users", v.CreateUser)
			r.router.GET("/users/{id}", v.GetUser)
			r.router.PUT("/users/{id}", v.UpdateUser)
			r.router.DELETE("/users/{id}", v.DeleteUser)
		case *handler.HealthHandler:
			r.router.GET("/health/live", v.Live)
			r.router.GET("/health/ready", v.Ready)
		case *handler.TimeoutHandler:
			r.router.GET("/timeout", v.Timeout)
		}
	}
}
