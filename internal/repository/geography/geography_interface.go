package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

// GeographyRepository defines the contract for operations related to countries, provinces, and localities.
type GeographyRepository interface {
	// CreateCountry inserts a new country record using the provided Executor.
	// Returns the created Country model or an error if the operation fails.
	CreateCountry(ctx context.Context, exec Executor, c models.Country) (*models.Country, error)

	// FindCountryByName retrieves a country by its name.
	// Returns the Country model or an error if no country is found.
	FindCountryByName(ctx context.Context, name string) (*models.Country, error)

	// CreateProvince inserts a new province associated with a country using the provided Executor.
	// Returns the created Province model or an error if the operation fails.
	CreateProvince(ctx context.Context, exec Executor, p models.Province) (*models.Province, error)

	// FindProvinceByName retrieves a province by its name and country ID.
	// Returns the Province model or an error if no province is found.
	FindProvinceByName(ctx context.Context, name string, countryId int) (*models.Province, error)

	// CreateLocality inserts a new locality using the provided Executor.
	// Returns the created Locality model or an error if the operation fails.
	CreateLocality(ctx context.Context, exec Executor, l models.Locality) (*models.Locality, error)

	// FindLocalityById retrieves a locality by its unique string identifier.
	// Returns the Locality model or an error if no locality is found.
	FindLocalityById(ctx context.Context, id string) (*models.Locality, error)

	// CountSellersByLocality returns the number of sellers registered in the specified locality.
	// Returns a response model or an error if the operation fails.
	CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error)

	// CountSellersGroupedByLocality returns the number of sellers grouped by each locality.
	// Returns a slice of response models, or an error if the operation fails.
	CountSellersGroupedByLocality(ctx context.Context) ([]models.ResponseLocalitySellers, error)

	// BeginTx starts a new database transaction and returns the transaction object.
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// CommitTx commits the provided transaction.
	CommitTx(tx *sql.Tx) error

	// RollbackTx aborts the provided transaction.
	RollbackTx(tx *sql.Tx) error

	// SetLogger allows injecting the logger after creation
	SetLogger(l logger.Logger)
}

// Executor wraps types that can execute SQL commands or queries within a context.
// Used for supporting operations both in transactions and out of them.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// geographyRepository implements the GeographyRepository interface using MySQL as the backend.
type geographyRepository struct {
	mysql  *sql.DB
	logger logger.Logger
}

// NewGeographyRepository creates a new GeographyRepository backed by a MySQL database.
func NewGeographyRepository(mysql *sql.DB) GeographyRepository {
	return &geographyRepository{
		mysql: mysql,
	}
}

// SetLogger allows you to inject the logger after creation
func (s *geographyRepository) SetLogger(l logger.Logger) {
	s.logger = l
}
