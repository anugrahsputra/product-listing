package repository

import (
	"context"
	"encoding/json"
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
		Price:       p.Price,
	}
	product, err := r.db.CreateProduct(ctx, params)
	if err != nil {
		return errors.New(err.Error())
	}

	for _, catID := range p.CategoryIDs {
		err := r.db.AddProductCategory(ctx, db.AddProductCategoryParams{
			ProductID:  product.ID,
			CategoryID: catID,
		})
		if err != nil {
			return err
		}
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

	result := toProductEntityByID(&product)

	return &result, nil
}

func (r *productRepository) FetchByCategory(ctx context.Context, cID uuid.UUID) ([]domain.Product, error) {
	products, err := r.db.GetProductsByCategoryID(ctx, cID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var result []domain.Product
	for _, p := range products {
		result = append(result, toProductEntityByCategoryID(&p))
	}

	return result, nil
}

func (r *productRepository) Update(ctx context.Context, id uuid.UUID, p domain.ProductInput) error {
	params := db.UpdateProductParams{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	err := r.db.UpdateProduct(ctx, params)
	if err != nil {
		return errors.New(err.Error())
	}

	if len(p.CategoryIDs) > 0 {
		err = r.db.ClearProductCategories(ctx, id)
		if err != nil {
			return err
		}

		for _, catID := range p.CategoryIDs {
			err = r.db.AddProductCategory(ctx, db.AddProductCategoryParams{
				ProductID:  id,
				CategoryID: catID,
			})
			if err != nil {
				return err
			}
		}
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

func parseCategories(data []byte) []domain.Category {
	if len(data) == 0 {
		return nil
	}
	var cats []domain.Category
	_ = json.Unmarshal(data, &cats)
	return cats
}

func toProductEntity(p *db.GetAllProductsRow) domain.Product {
	return domain.Product{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		Price:           p.Price,
		PrimaryImageURL: p.PrimaryImageUrl.String,
		Categories:      parseCategories(p.Categories),
		CreatedAt:       p.CreatedAt.Time,
		UpdatedAt:       p.UpdatedAt.Time,
	}
}

func toProductEntityByID(p *db.GetProductByIDRow) domain.Product {
	return domain.Product{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		Price:           p.Price,
		PrimaryImageURL: p.PrimaryImageUrl.String,
		Categories:      parseCategories(p.Categories),
		CreatedAt:       p.CreatedAt.Time,
		UpdatedAt:       p.UpdatedAt.Time,
	}
}

func toProductEntityByCategoryID(p *db.GetProductsByCategoryIDRow) domain.Product {
	return domain.Product{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		Price:           p.Price,
		PrimaryImageURL: p.PrimaryImageUrl.String,
		Categories:      parseCategories(p.Categories),
		CreatedAt:       p.CreatedAt.Time,
		UpdatedAt:       p.UpdatedAt.Time,
	}
}
