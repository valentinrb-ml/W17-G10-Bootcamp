package mocks

import (
    "context"
    
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

type CarryRepositoryMock struct {
    FuncCreate                           func(ctx context.Context, c carry.Carry) (*carry.Carry, error)
    FuncGetCarriesCountByAllLocalities   func(ctx context.Context) ([]carry.CarriesReport, error)
    FuncGetCarriesCountByLocalityID      func(ctx context.Context, localityID string) (*carry.CarriesReport, error)
}

func (m *CarryRepositoryMock) Create(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
    return m.FuncCreate(ctx, c)
}

func (m *CarryRepositoryMock) GetCarriesCountByAllLocalities(ctx context.Context) ([]carry.CarriesReport, error) {
    return m.FuncGetCarriesCountByAllLocalities(ctx)
}

func (m *CarryRepositoryMock) GetCarriesCountByLocalityID(ctx context.Context, localityID string) (*carry.CarriesReport, error) {
    return m.FuncGetCarriesCountByLocalityID(ctx, localityID)
}