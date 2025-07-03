package service

import (
	"fmt"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionService interface {
	FindAllSections() ([]section.Section, *api.ServiceError)
	FindById(id int) (section.Section, *api.ServiceError)
	DeleteSection(id int) *api.ServiceError
	CreateSection(sec section.Section) (section.Section, *api.ServiceError)
	UpdateSection(id int, sec section.RequestSection) (section.Section, *api.ServiceError)
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

func (s *SectionDefault) CreateSection(sec section.Section) (section.Section, *api.ServiceError) {
	if _, err1 := s.rpWareHouse.FindById(sec.WarehouseId); err1 != nil {
		fmt.Println("no hay tin")
		return section.Section{}, &api.ServiceError{
			Code:         err1.Code,
			ResponseCode: http.StatusBadRequest,
			Message:      "The warehouse does not exist",
		}
	}
	newSection, err := s.rp.CreateSection(sec)
	if err != nil {
		return section.Section{}, err
	}
	return newSection, nil

}
func (s *SectionDefault) UpdateSection(id int, sec section.RequestSection) (section.Section, *api.ServiceError) {
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
		if _, err1 := s.rpWareHouse.FindById(*sec.WarehouseId); err1 != nil {
			return section.Section{}, &api.ServiceError{
				Code:         err1.Code,
				ResponseCode: err1.ResponseCode,
				Message:      "The warehouse does not exist",
			}
		}
	}

	mappers.ApplySectionPatch(sec, &existing)

	secUpd, err := s.rp.UpdateSection(id, existing)

	if err != nil {
		return section.Section{}, err
	}

	return secUpd, nil

}
