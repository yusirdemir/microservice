package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TimeoutHandler struct{}

func NewTimeoutHandler() *TimeoutHandler {
	return &TimeoutHandler{}
}

func (h *TimeoutHandler) Register(r gin.IRouter) {
	r.GET("/timeout", h.Timeout)
}

func (h *TimeoutHandler) Timeout(c *gin.Context) {
	time.Sleep(6 * time.Second)
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
