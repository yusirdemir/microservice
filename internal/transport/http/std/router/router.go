package router

import (
	"net/http"

	"github.com/yusirdemir/microservice/internal/transport/http/std/handler"
)

type Router struct {
	mux      *http.ServeMux
	handlers []RouteHandler
}

type RouteHandler interface{}

func New(mux *http.ServeMux, handlers []RouteHandler) *Router {
	return &Router{
		mux:      mux,
		handlers: handlers,
	}
}

func (r *Router) SetupRoutes() {
	for _, h := range r.handlers {
		switch v := h.(type) {
		case *handler.UserHandler:
			r.mux.HandleFunc("POST /users", v.CreateUser)
			r.mux.HandleFunc("GET /users/{id}", v.GetUser)
			r.mux.HandleFunc("PUT /users/{id}", v.UpdateUser)
			r.mux.HandleFunc("DELETE /users/{id}", v.DeleteUser)
		case *handler.HealthHandler:
			r.mux.HandleFunc("GET /health/live", v.Live)
			r.mux.HandleFunc("GET /health/ready", v.Ready)
		case *handler.TimeoutHandler:
			r.mux.HandleFunc("GET /timeout", v.Timeout)
		}
	}
}
