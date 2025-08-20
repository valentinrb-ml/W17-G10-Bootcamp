package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

// Create creates a new carrier by delegating to the repository layer
// Returns the created carrier or an error if the operation fails
func (s *CarryDefault) Create(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
	result, err := s.rp.Create(ctx, c)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

// GetCarriesReport retrieves carrier statistics based on locality filtering
// If localityID is nil, returns statistics for all localities
// If localityID is provided, returns statistics for that specific locality
// Returns either a slice of CarriesReport (all localities) or a single CarriesReport (specific locality)
func (s *CarryDefault) GetCarriesReport(ctx context.Context, localityID *string) (interface{}, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "carry-service", "Starting carries report generation", map[string]interface{}{
			"locality_id": localityID,
		})
	}

	if localityID == nil {
		if s.logger != nil {
			s.logger.Debug(ctx, "carry-service", "Generating report for all localities")
		}

		result, err := s.rp.GetCarriesCountByAllLocalities(ctx)
		if err != nil {
			if s.logger != nil {
				s.logger.Error(ctx, "carry-service", "Failed to get carries count for all localities", err)
			}
			return nil, err
		}

		if s.logger != nil {
			s.logger.Info(ctx, "carry-service", "Successfully generated report for all localities")
		}
		return result, nil
	}

	if s.logger != nil {
		s.logger.Debug(ctx, "carry-service", "Generating report for specific locality", map[string]interface{}{
			"locality_id": *localityID,
		})
	}

	// Validate that locality exists
	l, err := s.rpGeo.FindLocalityById(ctx, *localityID)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "carry-service", "Error while validating locality", err, map[string]interface{}{
				"locality_id": *localityID,
			})
		}
		return nil, err
	}

	if l == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "carry-service", "Locality not found", map[string]interface{}{
				"locality_id": *localityID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "locality not found")
	}

	if s.logger != nil {
		s.logger.Debug(ctx, "carry-service", "Locality validation successful", map[string]interface{}{
			"locality_id":   *localityID,
			"locality_name": l.Name,
		})
	}

	result, err := s.rp.GetCarriesCountByLocalityID(ctx, *localityID)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "carry-service", "Failed to get carries count for specific locality", err, map[string]interface{}{
				"locality_id": *localityID,
			})
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.Info(ctx, "carry-service", "Successfully generated report for specific locality", map[string]interface{}{
			"locality_id": *localityID,
		})
	}
	return result, nil
}
