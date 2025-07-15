package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

type CarryService interface {
	Create(ctx context.Context, c carry.Carry) (*carry.Carry, error)
}

type CarryDefault struct {
	rp repository.CarryRepository
}

func NewCarryService(rp repository.CarryRepository) *CarryDefault {
	return &CarryDefault{rp: rp}
}