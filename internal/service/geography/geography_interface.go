package service

import (
	"context"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

type GeographyService interface {
	Create(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error)
	CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error)
	CountSellersGroupedByLocality(ctx context.Context) ([]models.ResponseLocalitySellers, error)

	// SetLogger allows injecting the logger after creation
	SetLogger(l logger.Logger)
}

type geographyService struct {
	rp     repository.GeographyRepository
	logger logger.Logger
}

func NewGeographyService(rp repository.GeographyRepository) GeographyService {
	return &geographyService{rp: rp}
}

// SetLogger allows you to inject the logger after creation
func (s *geographyService) SetLogger(l logger.Logger) {
	s.logger = l
}
