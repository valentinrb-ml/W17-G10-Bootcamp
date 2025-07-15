package repository

import (
	"context"
	"database/sql"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

type GeographyRepository interface {
	CreateCountry(ctx context.Context, exec Executor, c models.Country) (*models.Country, error)
	FindCountryByName(ctx context.Context, exec Executor, name string) (*models.Country, error)
	CreateProvince(ctx context.Context, exec Executor, p models.Province) (*models.Province, error)
	FindProvinceByName(ctx context.Context, exec Executor, name string, countryId int) (*models.Province, error)
	CreateLocality(ctx context.Context, exec Executor, l models.Locality) (*models.Locality, error)
	FindLocalityById(ctx context.Context, exec Executor, id string) (*models.Locality, error)
	CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error)

	BeginTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error
}

type geographyRepository struct {
	mysql *sql.DB
}

func NewGeographyRepository(mysql *sql.DB) GeographyRepository {
	return &geographyRepository{
		mysql: mysql,
	}
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *geographyRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.mysql.BeginTx(ctx, nil)
}

func (r *geographyRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *geographyRepository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}
