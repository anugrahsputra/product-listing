package usecase

import (
	"context"
	"product-listing/internal/domain"

	"github.com/google/uuid"
)

type ProductImageUsecase interface {
	AddImage(ctx context.Context, input domain.ProductImageInput) (*domain.ProductImage, error)
	GetProductImages(ctx context.Context, productID string) ([]domain.ProductImage, error)
	DeleteImage(ctx context.Context, id string) error
	SetPrimary(ctx context.Context, productID string, imageID string) error
}

type productImageUsecase struct {
	repo domain.ProductImageRepository
}

func NewProductImageUsecase(repo domain.ProductImageRepository) ProductImageUsecase {
	return &productImageUsecase{repo: repo}
}

func (u *productImageUsecase) AddImage(ctx context.Context, input domain.ProductImageInput) (*domain.ProductImage, error) {
	return u.repo.Create(ctx, input)
}

func (u *productImageUsecase) GetProductImages(ctx context.Context, productID string) ([]domain.ProductImage, error) {
	uid, err := uuid.Parse(productID)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByProductID(ctx, uid)
}

func (u *productImageUsecase) DeleteImage(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, uid)
}

func (u *productImageUsecase) SetPrimary(ctx context.Context, productID string, imageID string) error {
	puid, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	iuid, err := uuid.Parse(imageID)
	if err != nil {
		return err
	}
	return u.repo.SetPrimary(ctx, puid, iuid)
}
