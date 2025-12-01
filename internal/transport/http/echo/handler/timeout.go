package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type TimeoutHandler struct{}

func NewTimeoutHandler() *TimeoutHandler {
	return &TimeoutHandler{}
}

func (h *TimeoutHandler) Timeout(c echo.Context) error {
	time.Sleep(6 * time.Second)
	return c.JSON(http.StatusOK, echo.Map{"message": "Finished successfully"})
}
