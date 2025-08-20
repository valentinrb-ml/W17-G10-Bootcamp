package mocks_test

import (
	"context"
	"database/sql"
	"testing"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func TestGeographyRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.GeographyRepositoryMock{
		FuncCreateCountry: func(ctx context.Context, exec repository.Executor, c models.Country) (*models.Country, error) {
			return nil, nil
		},
		FuncFindCountryByName: func(ctx context.Context, name string) (*models.Country, error) { return nil, nil },
		FuncCreateProvince: func(ctx context.Context, exec repository.Executor, p models.Province) (*models.Province, error) {
			return nil, nil
		},
		FuncFindProvinceByName: func(ctx context.Context, name string, countryId int) (*models.Province, error) { return nil, nil },
		FuncCreateLocality: func(ctx context.Context, exec repository.Executor, l models.Locality) (*models.Locality, error) {
			return nil, nil
		},
		FuncFindLocalityById:              func(ctx context.Context, id string) (*models.Locality, error) { return nil, nil },
		FuncCountSellersByLocality:        func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) { return nil, nil },
		FuncCountSellersGroupedByLocality: func(ctx context.Context) ([]models.ResponseLocalitySellers, error) { return nil, nil },
		FuncBeginTx:                       func(ctx context.Context) (*sql.Tx, error) { return nil, nil },
		FuncCommitTx:                      func(tx *sql.Tx) error { return nil },
		FuncRollbackTx:                    func(tx *sql.Tx) error { return nil },
		FuncGetDB:                         func() *sql.DB { return nil },
		SetLoggerFn:                       func(l logger.Logger) {},
	}

	var exec repository.Executor
	var tx *sql.Tx

	m.CreateCountry(context.TODO(), exec, models.Country{})
	m.FindCountryByName(context.TODO(), "")
	m.CreateProvince(context.TODO(), exec, models.Province{})
	m.FindProvinceByName(context.TODO(), "", 0)
	m.CreateLocality(context.TODO(), exec, models.Locality{})
	m.FindLocalityById(context.TODO(), "")
	m.CountSellersByLocality(context.TODO(), "")
	m.CountSellersGroupedByLocality(context.TODO())
	m.BeginTx(context.TODO())
	m.CommitTx(tx)
	m.RollbackTx(tx)
	m.GetDB()
	m.SetLogger(nil)
}

func TestGeographyServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.GeographyServiceMock{
		CreateFn: func(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error) {
			return nil, nil
		},
		CountSellersByLocalityFn:        func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) { return nil, nil },
		CountSellersGroupedByLocalityFn: func(ctx context.Context) ([]models.ResponseLocalitySellers, error) { return nil, nil },
		SetLoggerFn:                     func(l logger.Logger) {},
	}

	m.Create(context.TODO(), models.RequestGeography{})
	m.CountSellersByLocality(context.TODO(), "")
	m.CountSellersGroupedByLocality(context.TODO())
	m.SetLogger(nil)
}
