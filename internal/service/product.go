package service

import (
	"context"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type productService struct{ repo repository.ProductRepository }

func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (s *productService) GetAll(ctx context.Context) ([]models.ProductResponse, error) {
	domainList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed to get all products")
	}
	return mappers.FromDomainList(domainList), nil
}

func (s *productService) Create(ctx context.Context, prod models.Product) (models.ProductResponse, error) {
	if err := validators.ValidateProductBusinessRules(prod); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return models.ProductResponse{}, err // Ya es AppError
		}
		return models.ProductResponse{}, apperrors.NewAppError(apperrors.CodeBadRequest, err.Error())
	}

	savedProduct, err := s.repo.Save(ctx, prod)
	if err != nil {
		return models.ProductResponse{}, err
	}

	return mappers.FromDomain(savedProduct), nil
}

func (s *productService) GetByID(ctx context.Context, id int) (models.ProductResponse, error) {
	currentProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.ProductResponse{}, err
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

func (s *productService) Patch(ctx context.Context, id int, req models.ProductPatchRequest) (models.ProductResponse, error) {
	updatedProduct, err := s.repo.Patch(ctx, id, req)
	if err != nil {
		return models.ProductResponse{}, err
	}

	return mappers.FromDomain(updatedProduct), nil
}
