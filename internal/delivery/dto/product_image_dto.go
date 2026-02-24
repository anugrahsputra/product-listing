package dto

import "time"

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
