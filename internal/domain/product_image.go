package domain

import (
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

func ToProductImageEntity(pi *db.ProductImage) ProductImage {
	return ProductImage{
		ID:        pi.ID,
		ProductID: pi.ProductID,
		Url:       pi.Url,
		IsPrimary: pi.IsPrimary.Bool,
		CreatedAt: pi.CreatedAt.Time,
	}
}
