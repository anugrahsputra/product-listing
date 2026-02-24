package domain

import (
	"context"
	"product-listing/internal/db"
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Url       string    `json:"url"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductImageInput struct {
	ProductID uuid.UUID
	Url       string
	IsPrimary bool
}

type ProductImageRepository interface {
	Create(ctx context.Context, input ProductImageInput) (*ProductImage, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]ProductImage, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetPrimary(ctx context.Context, productID uuid.UUID, imageID uuid.UUID) error
}

func ToProductImageEntity(pi *db.ProductImage) ProductImage {
	return ProductImage{
		ID:        pi.ID,
		ProductID: pi.ProductID,
		Url:       pi.Url,
		IsPrimary: pi.IsPrimary.Bool,
		CreatedAt: pi.CreatedAt.Time,
	}
}
