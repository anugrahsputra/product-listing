package repository

import (
	"context"
	"errors"
	"product-listing/config"
	"product-listing/internal/db"
	"product-listing/internal/domain"

	"github.com/google/uuid"
)

type productRepository struct {
	db *db.Queries
}

func NewProductRepository(database *config.Database) domain.ProductRepository {
	return &productRepository{
		db: db.New(database.Pool),
	}
}

func (r *productRepository) Create(ctx context.Context, p domain.ProductInput) error {
	params := db.CreateProductParams{
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		CategoryID:  p.CategoryID,
		Price:       p.Price,
	}
	_, err := r.db.CreateProduct(ctx, params)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *productRepository) Fetch(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	params := db.GetAllProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	products, err := r.db.GetAllProducts(ctx, params)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var result []domain.Product
	for _, p := range products {
		result = append(result, toProductEntity(&p))
	}

	return result, nil
}

func (r *productRepository) FetchCount(ctx context.Context) (int, error) {
	total, err := r.db.GetProductsCount(ctx)
	if err != nil {
		return int(0), err
	}

	return int(total), nil
}

func (r *productRepository) FetchById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := r.db.GetProductByID(ctx, id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	result := toProductEntity(&product)

	return &result, nil
}

func (r *productRepository) FetchByCategory(ctx context.Context, cID uuid.UUID) ([]domain.Product, error) {
	products, err := r.db.GetProductByCategory(ctx, cID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var result []domain.Product
	for _, p := range products {
		result = append(result, toProductEntity(&p))
	}

	return result, nil
}

func (r *productRepository) Update(ctx context.Context, id uuid.UUID, p domain.ProductInput) error {
	params := db.UpdateProductParams{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		CategoryID:  p.CategoryID,
	}

	err := r.db.UpdateProduct(ctx, params)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.DeleteProduct(ctx, id)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func toProductEntity(p *db.Product) domain.Product {
	return domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		CategoryID:  p.CategoryID,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}
