package service

import (
	"context"

	carryRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/carry"
	geographyRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

type CarryService interface {
	Create(ctx context.Context, c carry.Carry) (*carry.Carry, error)
	GetCarriesReport(ctx context.Context, localityID *string) (interface{}, error)
}

type CarryDefault struct {
	rp    carryRepo.CarryRepository
	rpGeo geographyRepo.GeographyRepository
}

func NewCarryService(rp carryRepo.CarryRepository, rpGeo geographyRepo.GeographyRepository) *CarryDefault {
	return &CarryDefault{
		rp:    rp,
		rpGeo: rpGeo,
	}
}
