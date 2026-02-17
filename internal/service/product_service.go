package service

import (
	"context"

	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/repository"
)

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
	product, err := domain.NewProduct("", userID, name, price, stock)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) GetAllProductsByUserID(ctx context.Context, userID string) ([]*domain.Product, error) {
	return s.repo.FindAllByUserID(ctx, userID)
}

func (s *productService) UpdateProduct(ctx context.Context, id, name string, price int, stock int) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
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
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
