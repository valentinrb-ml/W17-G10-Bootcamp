package service

import (
	"context"
	"fmt"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionService interface {
	FindAllSections() ([]section.Section, *api.ServiceError)
	FindById(id int) (section.Section, *api.ServiceError)
	DeleteSection(id int) *api.ServiceError
	CreateSection(ctx context.Context, sec section.Section) (section.Section, error)
	UpdateSection(ctx context.Context, id int, sec section.RequestSection) (section.Section, error)
}

type SectionDefault struct {
	rp          repository.SectionRepository
	rpWareHouse repository.WarehouseRepository
}

func NewSectionServer(rp repository.SectionRepository, rpWareHouse repository.WarehouseRepository) *SectionDefault {
	return &SectionDefault{
		rp:          rp,
		rpWareHouse: rpWareHouse,
	}
}

func (s *SectionDefault) FindAllSections() ([]section.Section, *api.ServiceError) {
	sections, err := s.rp.FindAllSections()
	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (s *SectionDefault) FindById(id int) (section.Section, *api.ServiceError) {
	sec, err := s.rp.FindById(id)
	if err != nil {
		return section.Section{}, err
	}
	return sec, nil
}

func (s *SectionDefault) DeleteSection(id int) *api.ServiceError {
	err := s.rp.DeleteSection(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SectionDefault) CreateSection(ctx context.Context, sec section.Section) (section.Section, error) {
	if _, err1 := s.rpWareHouse.FindById(ctx, sec.WarehouseId); err1 != nil {
		fmt.Println("no hay tin")
		return section.Section{}, apperrors.NewAppError(apperrors.CodeNotFound, "The warehouse does not exist")
	}
	newSection, err := s.rp.CreateSection(sec)
	if err != nil {
		return section.Section{}, err
	}
	return newSection, nil

}
func (s *SectionDefault) UpdateSection(ctx context.Context, id int, sec section.RequestSection) (section.Section, error) {
	existing, err := s.rp.FindById(id)
	if err != nil {
		return section.Section{}, err
	}
	if sec.SectionNumber != nil {
		if s.rp.ExistsSectionByNumber(*sec.SectionNumber) && *sec.SectionNumber != existing.SectionNumber {
			err := api.ServiceErrors[api.ErrConflict]
			return section.Section{}, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "The section number already exists",
			}
		}
	}
	if sec.WarehouseId != nil {
		if _, err1 := s.rpWareHouse.FindById(ctx, *sec.WarehouseId); err1 != nil {
			return section.Section{}, apperrors.NewAppError(apperrors.CodeNotFound, "The warehouse does not exist")
		}
	}

	mappers.ApplySectionPatch(sec, &existing)

	secUpd, err := s.rp.UpdateSection(id, existing)

	if err != nil {
		return section.Section{}, err
	}

	return secUpd, nil

}
