package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

// SellerRepository defines the business contract for the seller storage layer.
// It provides methods for creating, updating, deleting, and fetching seller records.
type SellerRepository interface {
	// Create inserts a new seller record into the storage.
	// Returns the created seller model or an error if the operation fails.
	Create(ctx context.Context, s models.Seller) (*models.Seller, error)

	// Update modifies the data of the seller specified by id.
	// Returns an error if the seller cannot be updated or does not exist.
	Update(ctx context.Context, id int, s models.Seller) error

	// Delete removes a seller record by its unique id.
	// Returns an error if the seller cannot be deleted or is not found.
	Delete(ctx context.Context, id int) error

	// FindAll retrieves all seller records from the storage.
	// Returns a slice of sellers, or an error if the operation fails.
	FindAll(ctx context.Context) ([]models.Seller, error)

	// FindById fetches a seller record by its unique id.
	// Returns the seller model or an error if the seller is not found.
	FindById(ctx context.Context, id int) (*models.Seller, error)

	// SetLogger allows injecting the logger after creation
	SetLogger(l logger.Logger)
}

// sellerRepository implements the SellerRepository interface using a MySQL backend.
type sellerRepository struct {
	mysql  *sql.DB
	logger logger.Logger
}

// NewSellerRepository creates a new SellerRepository backed by MySQL.
func NewSellerRepository(mysql *sql.DB) SellerRepository {
	return &sellerRepository{
		mysql: mysql,
	}
}

// SetLogger allows you to inject the logger after creation
func (s *sellerRepository) SetLogger(l logger.Logger) {
	s.logger = l
}
