package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/dto"
	"github.com/yusirdemir/microservice/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	ctx := c.Request().Context()
	user, err := h.service.CreateUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	ctx := c.Request().Context()
	user, err := h.service.UpdateUser(ctx, id, req.Name, "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	if err := h.service.DeleteUser(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func toUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}
