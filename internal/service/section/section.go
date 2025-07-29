package service

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

// FindAllSections fetches and returns all sections from the repository.
func (s *SectionDefault) FindAllSections(ctx context.Context) ([]models.Section, error) {
	sections, err := s.rp.FindAllSections(ctx)
	if err != nil {
		return nil, err
	}
	return sections, nil
}

// FindById retrieves a section by its ID using the repository.
func (s *SectionDefault) FindById(ctx context.Context, id int) (*models.Section, error) {
	sec, err := s.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return sec, nil
}

// DeleteSection removes a section by ID using the repository.
func (s *SectionDefault) DeleteSection(ctx context.Context, id int) error {
	err := s.rp.DeleteSection(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// CreateSection creates a new section using the repository.
func (s *SectionDefault) CreateSection(ctx context.Context, sec models.Section) (*models.Section, error) {
	newSection, err := s.rp.CreateSection(ctx, sec)
	if err != nil {
		return nil, err
	}
	return newSection, nil

}

// UpdateSection partially updates an existing section by applying a patch and persisting changes.
// Uses a mapper to apply only the changed fields.
func (s *SectionDefault) UpdateSection(ctx context.Context, id int, sec models.PatchSection) (*models.Section, error) {
	existing, err := s.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	mappers.ApplySectionPatch(sec, existing)

	secUpd, err := s.rp.UpdateSection(ctx, id, existing)

	if err != nil {
		return nil, err
	}

	return secUpd, nil

}
