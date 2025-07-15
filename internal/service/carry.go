package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

func (s *CarryDefault) Create(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
	ca, err := s.rp.Create(ctx, c)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (s *CarryDefault) GetCarriesReport(ctx context.Context, localityID *int) (interface{}, error) {
	if localityID == nil {
		return s.rp.GetCarriesCountByAllLocalities(ctx)
	}
	return s.rp.GetCarriesCountByLocalityID(ctx, *localityID)
}