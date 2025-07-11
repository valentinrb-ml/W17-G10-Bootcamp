package repository

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

// SectionRepository is an interface that represents a section repository
type SectionRepository interface {
	FindAllSections() ([]section.Section, *api.ServiceError)
	FindById(id int) (section.Section, *api.ServiceError)
	ExistsSectionById(id int) bool
	DeleteSection(id int) *api.ServiceError
	CreateSection(sec section.Section) (section.Section, *api.ServiceError)
	ExistsSectionByNumber(secNum int) bool
	UpdateSection(id int, sec section.Section) (section.Section, *api.ServiceError)
}
