package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryInput struct {
	Name string
	Slug string
}

type CategoryRepository interface {
	Create(ctx context.Context, c CategoryInput) error
	Fetch(ctx context.Context, limit, offset int) ([]Category, error)
	FetchById(ctx context.Context, id uuid.UUID) (*Category, error)
	FetchBySlug(ctx context.Context, slug string) (*Category, error)
	FetchCount(ctx context.Context) (int, error)
	Update(ctx context.Context, id uuid.UUID, c CategoryInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}
