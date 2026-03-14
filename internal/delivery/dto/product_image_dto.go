package dto

import (
	"product-listing/internal/domain"
	"time"
)

type ProductImageReq struct {
	ProductID string `json:"product_id"`
	Url       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductImageResp struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Url       string    `json:"url"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

func ToProductImageDTO(img *domain.ProductImage) ProductImageResp {
	return ProductImageResp{
		ID:        img.ID.String(),
		ProductID: img.ProductID.String(),
		Url:       img.Url,
		IsPrimary: img.IsPrimary,
		CreatedAt: img.CreatedAt,
	}
}
