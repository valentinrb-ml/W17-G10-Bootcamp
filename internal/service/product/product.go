package service

import (
	"context"
	"errors"
	productMappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product"
	productRepository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

const productServiceName = "product-service" // [LOG]

type productService struct {
	repo   productRepository.ProductRepository
	logger logger.Logger
}

func NewProductService(r productRepository.ProductRepository) ProductService {
	return &productService{repo: r}
}

// SetLogger allows you to inject the logger after creation
func (s *productService) SetLogger(l logger.Logger) {
	s.logger = l
}

// --- logging helpers (avoid repeating nil-check and set the service name) ---
func (s *productService) debug(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if s.logger != nil {
		s.logger.Debug(ctx, productServiceName, msg, md...) // [LOG]
	}
}
func (s *productService) info(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if s.logger != nil {
		s.logger.Info(ctx, productServiceName, msg, md...) // [LOG]
	}
}

func (s *productService) GetAll(ctx context.Context) ([]models.ProductResponse, error) {
	s.debug(ctx, "Fetching all products from repository") // [LOG]

	domainList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed to get all products")
	}

	s.debug(ctx, "Products fetched from repository", map[string]interface{}{ // [LOG]
		"count": len(domainList), // [LOG]
	})

	return productMappers.FromDomainList(domainList), nil
}

func (s *productService) Create(ctx context.Context, prod models.Product) (models.ProductResponse, error) {
	s.debug(ctx, "Validating product business rules", map[string]interface{}{ // [LOG]
		"product_code": prod.Code, // [LOG]
	})

	if err := validators.ValidateProductBusinessRules(prod); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return models.ProductResponse{}, err // Ya es AppError
		}
		return models.ProductResponse{}, apperrors.NewAppError(apperrors.CodeBadRequest, err.Error())
	}

	s.info(ctx, "Creating product", map[string]interface{}{ // [LOG]
		"product_code": prod.Code, // [LOG]
	})

	savedProduct, err := s.repo.Save(ctx, prod)
	if err != nil {
		return models.ProductResponse{}, err
	}

	s.info(ctx, "Product created", map[string]interface{}{ // [LOG]
		"product_id":   savedProduct.ID, // [LOG]
		"product_code": savedProduct.Code,
	})

	return productMappers.FromDomain(savedProduct), nil
}

func (s *productService) GetByID(ctx context.Context, id int) (models.ProductResponse, error) {
	s.debug(ctx, "Fetching product by id from repository", map[string]interface{}{ // [LOG]
		"product_id": id, // [LOG]
	})

	currentProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.ProductResponse{}, err
	}

	s.debug(ctx, "Product fetched from repository", map[string]interface{}{ // [LOG]
		"product_id": id, // [LOG]
	})

	return productMappers.FromDomain(currentProduct), nil
}

func (s *productService) Delete(ctx context.Context, id int) error {
	s.info(ctx, "Deleting product", map[string]interface{}{ // [LOG]
		"product_id": id, // [LOG]
	})

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.info(ctx, "Product deleted", map[string]interface{}{ // [LOG]
		"product_id": id, // [LOG]
	})

	return nil
}

func (s *productService) Patch(ctx context.Context, id int, req models.ProductPatchRequest) (models.ProductResponse, error) {
	s.info(ctx, "Patching product", map[string]interface{}{ // [LOG]
		"product_id": id, // [LOG]
	})
	updatedProduct, err := s.repo.Patch(ctx, id, req)
	if err != nil {
		return models.ProductResponse{}, err
	}

	s.info(ctx, "Product patched", map[string]interface{}{ // [LOG]
		"product_id": id, // [LOG]
	})

	return productMappers.FromDomain(updatedProduct), nil
}
