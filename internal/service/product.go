package service

import (
	"context"
	"errors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]product.ProductResponse, error)
	Create(ctx context.Context, prod product.Product) (product.ProductResponse, error)
	GetByID(ctx context.Context, id int) (product.ProductResponse, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, req product.ProductPatchRequest) (product.ProductResponse, error)
}

type productService struct{ repo repository.ProductRepository }

func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (s *productService) GetAll(ctx context.Context) ([]product.ProductResponse, error) {
	domainList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed to get all products")
	}
	return mappers.FromDomainList(domainList), nil
}

func (s *productService) Create(ctx context.Context, prod product.Product) (product.ProductResponse, error) {
	if err := validators.ValidateProductBusinessRules(prod); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return product.ProductResponse{}, err // Ya es AppError
		}
		return product.ProductResponse{}, apperrors.BadRequest(err.Error())
	}

	savedProduct, err := s.repo.Save(ctx, prod)
	if err != nil {
		return product.ProductResponse{}, err
	}

	return mappers.FromDomain(savedProduct), nil
}

func (s *productService) GetByID(ctx context.Context, id int) (product.ProductResponse, error) {
	currentProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return product.ProductResponse{}, err
	}

	return mappers.FromDomain(currentProduct), nil
}

func (s *productService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) Patch(ctx context.Context, id int, req product.ProductPatchRequest) (product.ProductResponse, error) {
	updatedProduct, err := s.repo.Patch(ctx, id, req)
	if err != nil {
		return product.ProductResponse{}, err
	}

	return mappers.FromDomain(updatedProduct), nil
}
