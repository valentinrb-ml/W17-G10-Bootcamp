package repository

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionMap struct {
	db map[int]section.Section
}

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

// NewSectionMap is a function that returns a new instance of SectionMap
func NewSectionMap(db map[int]section.Section) *SectionMap {
	defaultDb := make(map[int]section.Section)
	if db != nil {
		defaultDb = db
	}

	return &SectionMap{db: defaultDb}
}

func (r *SectionMap) FindAllSections() ([]section.Section, *api.ServiceError) {
	sections := make([]section.Section, 0, len(r.db))
	for _, sec := range r.db {
		sections = append(sections, sec)
	}
	return sections, nil
}

func (r *SectionMap) FindById(id int) (section.Section, *api.ServiceError) {

	sec, ok := r.db[id]

	if !ok {
		err := api.ServiceErrors[api.ErrNotFound]
		return section.Section{}, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The section you are looking for does not exist.",
		}
	}

	return sec, nil
}

func (r *SectionMap) ExistsSectionById(id int) bool {
	_, ok := r.db[id]
	return ok
}

func (r *SectionMap) DeleteSection(id int) *api.ServiceError {
	if !r.ExistsSectionById(id) {
		err := api.ServiceErrors[api.ErrNotFound]
		return &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The section you are trying to delete does not exist.",
		}
	}
	delete(r.db, id)
	return nil
}

func (r *SectionMap) getLastId() int {
	maxId := 0
	for id := range r.db {
		if id > maxId {
			maxId = id
		}
	}
	return maxId
}

func (r *SectionMap) ExistsSectionByNumber(secNum int) bool {

	for _, sec := range r.db {
		if sec.SectionNumber == secNum {
			return true
		}
	}
	return false

}

func (r *SectionMap) CreateSection(sec section.Section) (section.Section, *api.ServiceError) {
	if r.ExistsSectionByNumber(sec.SectionNumber) {
		err := api.ServiceErrors[api.ErrConflict]
		return section.Section{}, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The section number already exists",
		}
	}
	id := r.getLastId() + 1
	sec.Id = id
	r.db[sec.Id] = sec
	return sec, nil
}

func (r *SectionMap) UpdateSection(id int, sec section.Section) (section.Section, *api.ServiceError) {
	if !r.ExistsSectionById(id) {
		err := api.ServiceErrors[api.ErrNotFound]
		return section.Section{}, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "No sections found."}
	}

	r.db[id] = sec

	return sec, nil

}
