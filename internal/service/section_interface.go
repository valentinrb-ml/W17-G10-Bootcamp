package service

import (
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
