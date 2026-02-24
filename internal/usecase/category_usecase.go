package usecase

import (
	"context"
	"errors"
	"product-listing/internal/domain"

	"github.com/google/uuid"
)

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, c domain.CategoryInput) error
	GetCategories(ctx context.Context, limit, offset int) ([]domain.Category, error)
	GetCategoryById(ctx context.Context, id string) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	GetCategoryCount(ctx context.Context) (int, error)
	UpdateCategory(ctx context.Context, id string, c domain.CategoryInput) error
	DeleteCategory(ctx context.Context, id string) error
}

type categoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (u *categoryUsecase) CreateCategory(ctx context.Context, c domain.CategoryInput) error {
	if c.Name == "" {
		return errors.New("Category name cannot be empty")
	}

	if c.Slug == "" {
		return errors.New("Category slug cannot be empty")
	}

	err := u.repo.Create(ctx, c)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (u *categoryUsecase) GetCategories(ctx context.Context, page, limit int) ([]domain.Category, error) {
	offset := (page - 1) * limit
	categories, err := u.repo.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (u *categoryUsecase) GetCategoryCount(ctx context.Context) (int, error) {
	total, err := u.repo.FetchCount(ctx)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (u *categoryUsecase) GetCategoryById(ctx context.Context, id string) (*domain.Category, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	category, err := u.repo.FetchById(ctx, uid)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *categoryUsecase) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := u.repo.FetchBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *categoryUsecase) UpdateCategory(ctx context.Context, id string, c domain.CategoryInput) error {
	uid, _ := uuid.Parse(id)

	err := u.repo.Update(ctx, uid, c)
	if err != nil {
		return err
	}

	return nil
}

func (u *categoryUsecase) DeleteCategory(ctx context.Context, id string) error {
	uid, _ := uuid.Parse(id)

	err := u.repo.Delete(ctx, uid)
	if err != nil {
		return err
	}

	return nil
}
