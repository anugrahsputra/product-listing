package dto

import (
	"product-listing/internal/domain"
	"time"
)

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

func ToProductDTO(p *domain.Product) ProductResp {
	categories := make([]CategoryResp, 0, len(p.Categories))
	for _, c := range p.Categories {
		item := ToCategoryDTO(&c)
		categories = append(categories, item)

	}

	return ProductResp{
		ID:              p.ID.String(),
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		Price:           p.Price,
		PrimaryImageURL: p.PrimaryImageURL,
		Categories:      categories,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}
