package service

import (
	"context"

	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var productTracer = otel.Tracer("microservice/service/product")

type ProductService interface {
	CreateProduct(ctx context.Context, userID string, name string, price int, stock int) (*domain.Product, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	GetAllProductsByUserID(ctx context.Context, userID string) ([]*domain.Product, error)
	UpdateProduct(ctx context.Context, id, name string, price int, stock int) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, userID string, name string, price int, stock int) (*domain.Product, error) {
	ctx, span := productTracer.Start(ctx, "ProductService.CreateProduct")
	defer span.End()

	span.SetAttributes(attribute.String("app.user.id", userID))

	product, err := domain.NewProduct("", userID, name, price, stock)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(
		attribute.String("app.product.id", product.ID),
		attribute.Int("app.product.price", product.Price),
	)

	if err := s.repo.Create(ctx, product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	ctx, span := productTracer.Start(ctx, "ProductService.GetProduct")
	defer span.End()

	span.SetAttributes(attribute.String("app.product.id", id))

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return product, nil
}

func (s *productService) GetAllProductsByUserID(ctx context.Context, userID string) ([]*domain.Product, error) {
	ctx, span := productTracer.Start(ctx, "ProductService.GetAllByUserID")
	defer span.End()

	span.SetAttributes(attribute.String("app.user.id", userID))

	products, err := s.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Int("app.product.count", len(products)))

	return products, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id, name string, price int, stock int) (*domain.Product, error) {
	ctx, span := productTracer.Start(ctx, "ProductService.UpdateProduct")
	defer span.End()

	span.SetAttributes(attribute.String("app.product.id", id))

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if name != "" {
		product.Name = name
	}
	if price > 0 {
		product.Price = price
	}
	if stock >= 0 {
		product.Stock = stock
	}

	if err := s.repo.Update(ctx, product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	ctx, span := productTracer.Start(ctx, "ProductService.DeleteProduct")
	defer span.End()

	span.SetAttributes(attribute.String("app.product.id", id))

	if err := s.repo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
