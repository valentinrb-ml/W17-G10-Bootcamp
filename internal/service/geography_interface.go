package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

type GeographyService interface {
	Create(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error)
	CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error)
}

type geographyService struct {
	rp repository.GeographyRepository
}

func NewGeographyService(rp repository.GeographyRepository) GeographyService {
	return &geographyService{rp: rp}
}
