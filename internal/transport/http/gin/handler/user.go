package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/dto"
	"github.com/yusirdemir/microservice/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Register(r gin.IRouter) {
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.PUT("/users/:id", h.UpdateUser)
	r.DELETE("/users/:id", h.DeleteUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.CreateUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.UpdateUser(ctx, id, req.Name, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	if err := h.service.DeleteUser(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID(),
		Name:      u.Name(),
		Email:     u.Email(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}
