package dto

import "time"

type ProductResp struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"Description"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductReq struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"Description"`
	CategoryID  string  `json:"category_id"`
	Price       float64 `json:"price"`
}
