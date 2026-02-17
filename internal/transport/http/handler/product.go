package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/dto"
	"github.com/yusirdemir/microservice/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) Register(r fiber.Router) {
	r.Post("/products", h.CreateProduct)
	r.Get("/products/:id", h.GetProduct)
	r.Get("/users/:id/products", h.GetUserProducts)
	r.Put("/products/:id", h.UpdateProduct)
	r.Delete("/products/:id", h.DeleteProduct)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx := c.UserContext()

	userID := c.Get("X-User-ID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "X-User-ID header is required"})
	}

	product, err := h.service.CreateProduct(ctx, userID, req.Name, req.Price, req.Stock)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(toProductResponse(product))
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(toProductResponse(product))
}

func (h *ProductHandler) GetUserProducts(c *fiber.Ctx) error {
	userID := c.Params("id")
	ctx := c.UserContext()

	products, err := h.service.GetAllProductsByUserID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		response[i] = toProductResponse(p)
	}

	return c.JSON(response)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx := c.UserContext()
	product, err := h.service.UpdateProduct(ctx, id, req.Name, req.Price, req.Stock)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(toProductResponse(product))
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()

	if err := h.service.DeleteProduct(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func toProductResponse(p *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
	}
}
