package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

// Create creates a new carrier by delegating to the repository layer
// Returns the created carrier or an error if the operation fails
func (s *CarryDefault) Create(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
	return s.rp.Create(ctx, c)
}

// GetCarriesReport retrieves carrier statistics based on locality filtering
// If localityID is nil, returns statistics for all localities
// If localityID is provided, returns statistics for that specific locality
// Returns either a slice of CarriesReport (all localities) or a single CarriesReport (specific locality)
func (s *CarryDefault) GetCarriesReport(ctx context.Context, localityID *string) (interface{}, error) {
	if localityID == nil {
		return s.rp.GetCarriesCountByAllLocalities(ctx)
	}

	l, _ := s.rpGeo.FindLocalityById(ctx, *localityID)
	if l == nil {
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "locality not found")
	}
	return s.rp.GetCarriesCountByLocalityID(ctx, *localityID)
}
