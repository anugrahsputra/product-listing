package repository

import (
	"context"
	"errors"
	"product-listing/config"
	"product-listing/internal/db"
	"product-listing/internal/domain"

	"github.com/google/uuid"
)

type categoryRepository struct {
	db *db.Queries
}

func NewCategoryRepository(database *config.Database) domain.CategoryRepository {
	return &categoryRepository{
		db: db.New(database.Pool),
	}
}

func (r *categoryRepository) Create(ctx context.Context, c domain.CategoryInput) error {
	params := db.CreateCategoryParams{
		Name: c.Name,
		Slug: c.Slug,
	}

	_, err := r.db.CreateCategory(ctx, params)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *categoryRepository) Fetch(ctx context.Context, limit, offset int) ([]domain.Category, error) {
	params := db.GetCategoriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	categories, err := r.db.GetCategories(ctx, params)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var result []domain.Category
	for _, c := range categories {
		result = append(result, toCategoryEntity(&c))
	}

	return result, nil
}

func (r *categoryRepository) FetchById(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category, err := r.db.GetCategoryById(ctx, id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	result := toCategoryEntity(&category)

	return &result, nil
}

func (r *categoryRepository) FetchBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := r.db.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	result := toCategoryEntity(&category)

	return &result, nil
}

func (r *categoryRepository) FetchCount(ctx context.Context) (int, error) {
	total, err := r.db.GetCategoriesCount(ctx)
	if err != nil {
		return 0, err
	}

	return int(total), err
}

func (r *categoryRepository) Update(ctx context.Context, id uuid.UUID, c domain.CategoryInput) error {
	params := db.UpdateCategoryParams{
		ID:   id,
		Name: c.Name,
		Slug: c.Slug,
	}

	err := r.db.UpdateCategory(ctx, params)
	if err != nil {
		return errors.New(err.Error())

	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.DeleteCategory(ctx, id)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func toCategoryEntity(c *db.Category) domain.Category {
	return domain.Category{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}
