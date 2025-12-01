package handler

import (
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type TimeoutHandler struct{}

func NewTimeoutHandler() *TimeoutHandler {
	return &TimeoutHandler{}
}

func (h *TimeoutHandler) Timeout(ctx *fasthttp.RequestCtx) {
	time.Sleep(6 * time.Second)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString(`{"message": "Finished successfully"}`)
}
