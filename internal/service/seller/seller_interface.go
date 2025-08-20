package service

import (
	"context"

	geographyRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	sellerRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

// SellerService defines the business logic methods for working with sellers.
// It provides high-level operations to create, update, delete, and retrieve sellers.
type SellerService interface {
	// Create adds a new seller using the provided request data.
	// Returns the newly created seller response, or an error if the operation fails.
	Create(ctx context.Context, reqs models.RequestSeller) (*models.ResponseSeller, error)

	// Update modifies an existing seller identified by id using the provided request data.
	// Returns the updated seller response, or an error if the operation fails.
	Update(ctx context.Context, id int, reqs models.RequestSeller) (*models.ResponseSeller, error)

	// Delete removes a seller identified by id.
	// Returns an error if the seller does not exist or the operation fails.
	Delete(ctx context.Context, id int) error

	// FindAll retrieves all sellers.
	// Returns a slice of seller responses or an error if the operation fails.
	FindAll(ctx context.Context) ([]models.ResponseSeller, error)

	// FindById retrieves a seller by their unique id.
	// Returns the seller response or an error if not found.
	FindById(ctx context.Context, id int) (*models.ResponseSeller, error)

	// SetLogger allows injecting the logger after creation
	SetLogger(l logger.Logger)
}

// sellerService implements the SellerService interface.
type sellerService struct {
	sellerRepo    sellerRepo.SellerRepository
	geographyRepo geographyRepo.GeographyRepository
	logger        logger.Logger
}

// NewSellerService creates a new SellerService for managing sellers.
// sellerRepo provides access to seller data; geoRepo provides access to geographical data.
func NewSellerService(sellerRepo sellerRepo.SellerRepository, geographyRepo geographyRepo.GeographyRepository) SellerService {
	return &sellerService{
		sellerRepo:    sellerRepo,
		geographyRepo: geographyRepo,
	}
}

// SetLogger allows you to inject the logger after creation
func (s *sellerService) SetLogger(l logger.Logger) {
	s.logger = l
}
