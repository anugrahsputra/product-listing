package repository

import (
	"context"
	"product-listing/config"
	"product-listing/internal/db"
	"product-listing/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type productImageRepository struct {
	db *db.Queries
}

func NewProductImageRepository(database *config.Database) domain.ProductImageRepository {
	return &productImageRepository{
		db: db.New(database.Pool),
	}
}

func (r *productImageRepository) Create(ctx context.Context, input domain.ProductImageInput) (*domain.ProductImage, error) {
	params := db.CreateProductImageParams{
		ProductID: input.ProductID,
		Url:       input.Url,
		IsPrimary: pgtype.Bool{Bool: input.IsPrimary, Valid: true},
	}

	pi, err := r.db.CreateProductImage(ctx, params)
	if err != nil {
		return nil, err
	}

	entity := toProductImageEntity(&pi)
	return &entity, nil
}

func (r *productImageRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error) {
	images, err := r.db.GetProductImages(ctx, productID)
	if err != nil {
		return nil, err
	}

	var result []domain.ProductImage
	for _, img := range images {
		result = append(result, toProductImageEntity(&img))
	}
	return result, nil
}

func (r *productImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteProductImage(ctx, id)
}

func (r *productImageRepository) SetPrimary(ctx context.Context, productID uuid.UUID, imageID uuid.UUID) error {
	params := db.SetProductPrimaryImageParams{
		ProductID: productID,
		ID:        imageID,
	}
	return r.db.SetProductPrimaryImage(ctx, params)
}

func toProductImageEntity(pi *db.ProductImage) domain.ProductImage {
	return domain.ProductImage{
		ID:        pi.ID,
		ProductID: pi.ProductID,
		Url:       pi.Url,
		IsPrimary: pi.IsPrimary.Bool,
		CreatedAt: pi.CreatedAt.Time,
	}
}
