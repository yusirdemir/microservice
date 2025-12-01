package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"status": "UP"})
}

func (h *HealthHandler) Ready(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"status": "UP"})
}
