package mocks

import (
	"context"
	"database/sql"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

type GeographyRepositoryMock struct {
	FuncCreateCountry                 func(ctx context.Context, exec repository.Executor, c models.Country) (*models.Country, error)
	FuncFindCountryByName             func(ctx context.Context, name string) (*models.Country, error)
	FuncCreateProvince                func(ctx context.Context, exec repository.Executor, p models.Province) (*models.Province, error)
	FuncFindProvinceByName            func(ctx context.Context, name string, countryId int) (*models.Province, error)
	FuncCreateLocality                func(ctx context.Context, exec repository.Executor, l models.Locality) (*models.Locality, error)
	FuncFindLocalityById              func(ctx context.Context, id string) (*models.Locality, error)
	FuncCountSellersByLocality        func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error)
	FuncCountSellersGroupedByLocality func(ctx context.Context) ([]models.ResponseLocalitySellers, error)
	FuncBeginTx                       func(ctx context.Context) (*sql.Tx, error)
	FuncCommitTx                      func(tx *sql.Tx) error
	FuncRollbackTx                    func(tx *sql.Tx) error
	FuncGetDB                         func() *sql.DB
	SetLoggerFn                       func(l logger.Logger)
}

func (m *GeographyRepositoryMock) CreateCountry(ctx context.Context, exec repository.Executor, c models.Country) (*models.Country, error) {
	if m.FuncCreateCountry != nil {
		return m.FuncCreateCountry(ctx, exec, c)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) FindCountryByName(ctx context.Context, name string) (*models.Country, error) {
	if m.FuncFindCountryByName != nil {
		return m.FuncFindCountryByName(ctx, name)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) CreateProvince(ctx context.Context, exec repository.Executor, p models.Province) (*models.Province, error) {
	if m.FuncCreateProvince != nil {
		return m.FuncCreateProvince(ctx, exec, p)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) FindProvinceByName(ctx context.Context, name string, countryId int) (*models.Province, error) {
	if m.FuncFindProvinceByName != nil {
		return m.FuncFindProvinceByName(ctx, name, countryId)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) CreateLocality(ctx context.Context, exec repository.Executor, l models.Locality) (*models.Locality, error) {
	if m.FuncCreateLocality != nil {
		return m.FuncCreateLocality(ctx, exec, l)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) FindLocalityById(ctx context.Context, id string) (*models.Locality, error) {
	if m.FuncFindLocalityById != nil {
		return m.FuncFindLocalityById(ctx, id)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
	if m.FuncCountSellersByLocality != nil {
		return m.FuncCountSellersByLocality(ctx, id)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) CountSellersGroupedByLocality(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
	if m.FuncCountSellersGroupedByLocality != nil {
		return m.FuncCountSellersGroupedByLocality(ctx)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if m.FuncBeginTx != nil {
		return m.FuncBeginTx(ctx)
	}
	return nil, nil
}

func (m *GeographyRepositoryMock) CommitTx(tx *sql.Tx) error {
	if m.FuncCommitTx != nil {
		return m.FuncCommitTx(tx)
	}
	return nil
}

func (m *GeographyRepositoryMock) RollbackTx(tx *sql.Tx) error {
	if m.FuncRollbackTx != nil {
		return m.FuncRollbackTx(tx)
	}
	return nil
}

func (m *GeographyRepositoryMock) GetDB() *sql.DB {
	if m.FuncGetDB != nil {
		return m.FuncGetDB()
	}
	return nil
}

func (m *GeographyRepositoryMock) SetLogger(l logger.Logger) {
	if m.SetLoggerFn != nil {
		m.SetLoggerFn(l)
	}
}
