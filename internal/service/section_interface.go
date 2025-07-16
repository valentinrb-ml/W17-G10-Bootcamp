package service

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

// SectionService defines the business logic interface for sections.
type SectionService interface {
	FindAllSections(ctx context.Context) ([]models.Section, error)
	FindById(ctx context.Context, id int) (*models.Section, error)
	DeleteSection(ctx context.Context, id int) error
	CreateSection(ctx context.Context, sec models.Section) (*models.Section, error)
	UpdateSection(ctx context.Context, id int, sec models.PatchSection) (*models.Section, error)
}

// SectionDefault is the default implementation of SectionService.
type SectionDefault struct {
	rp          repository.SectionRepository
	rpWareHouse repository.WarehouseRepository
}

// NewSectionServer creates a new SectionDefault service with the given SectionRepository.
func NewSectionServer(rp repository.SectionRepository) *SectionDefault {
	return &SectionDefault{
		rp: rp,
	}
}
