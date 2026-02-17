package repository

import (
	"context"

	"github.com/yusirdemir/microservice/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
}
