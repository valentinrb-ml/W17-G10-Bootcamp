package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

// Create creates a new carrier by delegating to the repository layer
// Returns the created carrier or an error if the operation fails
func (s *CarryDefault) Create(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
	ca, err := s.rp.Create(ctx, c)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

// GetCarriesReport retrieves carrier statistics based on locality filtering
// If localityID is nil, returns statistics for all localities
// If localityID is provided, returns statistics for that specific locality
// Returns either a slice of CarriesReport (all localities) or a single CarriesReport (specific locality)
func (s *CarryDefault) GetCarriesReport(ctx context.Context, localityID *string) (interface{}, error) {
	if localityID == nil {
		return s.rp.GetCarriesCountByAllLocalities(ctx)
	}
	return s.rp.GetCarriesCountByLocalityID(ctx, *localityID)
}
