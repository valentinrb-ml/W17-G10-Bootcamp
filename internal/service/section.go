package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionDefault struct {
	rp          repository.SectionRepository
	rpWareHouse repository.WarehouseRepository
}

func NewSectionServer(rp repository.SectionRepository) *SectionDefault {
	return &SectionDefault{
		rp: rp,
	}
}

func (s *SectionDefault) FindAllSections(ctx context.Context) ([]section.Section, error) {
	sections, err := s.rp.FindAllSections(ctx)
	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (s *SectionDefault) FindById(ctx context.Context, id int) (*section.Section, error) {
	sec, err := s.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return sec, nil
}

func (s *SectionDefault) DeleteSection(ctx context.Context, id int) error {
	err := s.rp.DeleteSection(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SectionDefault) CreateSection(ctx context.Context, sec section.Section) (*section.Section, error) {
	newSection, err := s.rp.CreateSection(ctx, sec)
	if err != nil {
		return nil, err
	}
	return newSection, nil

}

func (s *SectionDefault) UpdateSection(ctx context.Context, id int, sec section.PatchSection) (*section.Section, error) {
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
