package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID              uuid.UUID
	Name            string
	Slug            string
	Description     string
	Price           float64
	PrimaryImageURL string
	Categories      []Category
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ProductInput struct {
	Name        string
	Slug        string
	Description string
	CategoryIDs []uuid.UUID
	Price       float64
}

type ProductRepository interface {
	Create(ctx context.Context, p ProductInput) error
	Fetch(ctx context.Context, limit, offset int) ([]Product, error)
	FetchById(ctx context.Context, id uuid.UUID) (*Product, error)
	FetchByCategory(ctx context.Context, cID uuid.UUID) ([]Product, error)
	FetchCount(ctx context.Context) (int, error)
	Update(ctx context.Context, id uuid.UUID, p ProductInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}
