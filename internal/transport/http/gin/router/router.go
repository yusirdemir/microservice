package router

import (
	"github.com/gin-gonic/gin"
)

type RouteHandler interface {
	Register(r gin.IRouter)
}

type Router struct {
	App      *gin.Engine
	Handlers []RouteHandler
}

func New(app *gin.Engine, handlers []RouteHandler) *Router {
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
