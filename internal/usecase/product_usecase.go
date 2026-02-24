package usecase

import (
	"context"
	"product-listing/internal/domain"

	"github.com/google/uuid"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, p domain.ProductInput) error
	GetProducts(ctx context.Context, page, limit int) ([]domain.Product, error)
	GetProductCount(ctx context.Context) (int, error)
	GetProductsById(ctx context.Context, id string) (*domain.Product, error)
	GetProductsByCategory(ctx context.Context, cID string) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, id string, p domain.ProductInput) error
	DeleteProduct(ctx context.Context, id string) error
}

type productUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) CreateProduct(ctx context.Context, p domain.ProductInput) error {
	err := u.repo.Create(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

func (u *productUsecase) GetProducts(ctx context.Context, page, limit int) ([]domain.Product, error) {
	offset := (page - 1) * limit
	products, err := u.repo.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *productUsecase) GetProductCount(ctx context.Context) (int, error) {
	total, err := u.repo.FetchCount(ctx)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (u *productUsecase) GetProductsById(ctx context.Context, id string) (*domain.Product, error) {
	uid, _ := uuid.Parse(id)

	product, err := u.repo.FetchById(ctx, uid)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUsecase) GetProductsByCategory(ctx context.Context, cID string) ([]domain.Product, error) {
	uid, _ := uuid.Parse(cID)
	products, err := u.repo.FetchByCategory(ctx, uid)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *productUsecase) UpdateProduct(ctx context.Context, id string, p domain.ProductInput) error {
	uid, _ := uuid.Parse(id)
	err := u.repo.Update(ctx, uid, p)
	if err != nil {
		return err
	}

	return nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	uid, _ := uuid.Parse(id)
	err := u.repo.Delete(ctx, uid)
	if err != nil {
		return err
	}

	return nil
}
