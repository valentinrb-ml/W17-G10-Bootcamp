package mocks

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

// GeographyServiceMock implements GeographyService, but allows customizing behavior in tests.
type GeographyServiceMock struct {
	CreateFn                        func(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error)
	CountSellersByLocalityFn        func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error)
	CountSellersGroupedByLocalityFn func(ctx context.Context) ([]models.ResponseLocalitySellers, error)
	SetLoggerFn                     func(l logger.Logger)
}

func (g *GeographyServiceMock) Create(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error) {
	return g.CreateFn(ctx, gr)
}

func (g *GeographyServiceMock) CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
	return g.CountSellersByLocalityFn(ctx, id)
}

func (g *GeographyServiceMock) CountSellersGroupedByLocality(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
	return g.CountSellersGroupedByLocalityFn(ctx)
}

func (m *GeographyServiceMock) SetLogger(l logger.Logger) {
	if m.SetLoggerFn != nil {
		m.SetLoggerFn(l)
	}
}
