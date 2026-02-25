package dto

import "time"

type ProductResp struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Slug            string         `json:"slug"`
	Description     string         `json:"Description"`
	Price           float64        `json:"price"`
	PrimaryImageURL string         `json:"primary_image_url"`
	Categories      []CategoryResp `json:"categories"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type ProductReq struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"Description"`
	CategoryIDs []string `json:"category_ids"`
	Price       float64  `json:"price"`
}
