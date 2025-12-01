package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
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

func (h *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	var req dto.CreateUserRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString(`{"error": "Invalid request body"}`)
		return
	}

	c := context.Background()
	user, err := h.service.CreateUser(c, req.Name, req.Email, req.Password)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(`{"error": "` + err.Error() + `"}`)
		return
	}

	resp, _ := json.Marshal(toUserResponse(user))
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(resp)
}

func (h *UserHandler) GetUser(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	c := context.Background()
	user, err := h.service.GetUser(c, id)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString(`{"error": "User not found"}`)
		return
	}

	resp, _ := json.Marshal(toUserResponse(user))
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(resp)
}

func (h *UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	var req dto.UpdateUserRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString(`{"error": "Invalid request body"}`)
		return
	}

	c := context.Background()
	user, err := h.service.UpdateUser(c, id, req.Name, "")
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(`{"error": "` + err.Error() + `"}`)
		return
	}

	resp, _ := json.Marshal(toUserResponse(user))
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(resp)
}

func (h *UserHandler) DeleteUser(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	c := context.Background()
	if err := h.service.DeleteUser(c, id); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(`{"error": "` + err.Error() + `"}`)
		return
	}

	ctx.SetStatusCode(http.StatusNoContent)
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
