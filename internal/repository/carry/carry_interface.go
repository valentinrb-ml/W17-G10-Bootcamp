package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

type CarryRepository interface {
	Create(ctx context.Context, c carry.Carry) (*carry.Carry, error)
	GetCarriesCountByAllLocalities(ctx context.Context) ([]carry.CarriesReport, error)
	GetCarriesCountByLocalityID(ctx context.Context, localityID string) (*carry.CarriesReport, error)
}

type CarryMySQL struct {
	db *sql.DB
}

func NewCarryRepository(db *sql.DB) *CarryMySQL {
	return &CarryMySQL{db}
}