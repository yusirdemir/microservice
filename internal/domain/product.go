package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProduct(id string, userID string, name string, price int, stock int) (*Product, error) {

	if id == "" {
		id = uuid.New().String()
	}

	if userID == "" {
		return nil, errors.New("user_id cannot be empty")
	}

	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	now := time.Now()
	return &Product{
		ID:        id,
		UserID:    userID,
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func ReconstituteProduct(id string, userID string, name string, price int, stock int, createdAt time.Time, updatedAt time.Time) *Product {
	return &Product{
		ID:        id,
		UserID:    userID,
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
