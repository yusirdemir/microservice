package handler

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Live(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString(`{"status": "UP"}`)
}

func (h *HealthHandler) Ready(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString(`{"status": "UP"}`)
}
