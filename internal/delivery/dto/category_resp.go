package dto

import (
	"product-listing/internal/domain"
	"time"
)

type CategoryResp struct {
	ID        string    `json:"id"`
	Name      string    `json:"Name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryReq struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func ToCategoryDTO(c *domain.Category) CategoryResp {
	return CategoryResp{
		ID:        c.ID.String(),
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
