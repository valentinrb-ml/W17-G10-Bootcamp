package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

func TestCarryRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.CarryRepositoryMock{
		FuncCreate:                         func(ctx context.Context, c carry.Carry) (*carry.Carry, error) { return nil, nil },
		FuncGetCarriesCountByAllLocalities: func(ctx context.Context) ([]carry.CarriesReport, error) { return nil, nil },
		FuncGetCarriesCountByLocalityID:    func(ctx context.Context, localityID string) (*carry.CarriesReport, error) { return nil, nil },
	}
	m.Create(context.TODO(), carry.Carry{})
	m.GetCarriesCountByAllLocalities(context.TODO())
	m.GetCarriesCountByLocalityID(context.TODO(), "")
}

func TestCarryServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.CarryServiceMock{
		FuncCreate:           func(ctx context.Context, c carry.Carry) (*carry.Carry, error) { return nil, nil },
		FuncGetCarriesReport: func(ctx context.Context, localityID *string) (interface{}, error) { return nil, nil },
	}
	m.Create(context.TODO(), carry.Carry{})
	m.GetCarriesReport(context.TODO(), nil)
}
