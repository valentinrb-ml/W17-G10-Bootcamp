package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

type CarryService interface {
	Create(ctx context.Context, c carry.Carry) (*carry.Carry, error)
	GetCarriesReport(ctx context.Context, localityID *string) (interface{}, error)
}

type CarryDefault struct {
	rp repository.CarryRepository
	rpGeo repository.GeographyRepository
}

func NewCarryService(rp repository.CarryRepository, rpGeo repository.GeographyRepository) *CarryDefault {
	return &CarryDefault{
		rp:    rp,
		rpGeo: rpGeo,
	}
}