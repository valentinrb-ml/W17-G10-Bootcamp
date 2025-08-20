package service

import (
	"context"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
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
	rp     repository.SectionRepository
	logger logger.Logger
}

// SetLogger allows injecting the logger after creation
func (s *SectionDefault) SetLogger(l logger.Logger) {
	s.logger = l
}

// NewSectionServer creates a new SectionDefault service with the given SectionRepository.
func NewSectionService(rp repository.SectionRepository) *SectionDefault {
	return &SectionDefault{
		rp: rp,
	}
}
