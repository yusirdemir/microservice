package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/repository"
)

type memoryProductRepository struct {
	products map[string]*domain.Product
	mu       sync.RWMutex
}

func NewProductRepository() repository.ProductRepository {
	return &memoryProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (r *memoryProductRepository) Create(ctx context.Context, product *domain.Product) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return errors.New("product already exists")
	}

	r.products[product.ID] = product
	return nil
}

func (r *memoryProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	return product, nil
}

func (r *memoryProductRepository) FindAllByUserID(ctx context.Context, userID string) ([]*domain.Product, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var products []*domain.Product
	for _, p := range r.products {
		if p.UserID == userID {
			products = append(products, p)
		}
	}
	return products, nil
}

func (r *memoryProductRepository) Update(ctx context.Context, product *domain.Product) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}

	r.products[product.ID] = product
	return nil
}

func (r *memoryProductRepository) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)
	return nil
}
