package mappers

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

func RequestSectionToSection(req section.RequestSection) section.Section {
	return section.Section{
		SectionNumber:      *req.SectionNumber,
		CurrentTemperature: *req.CurrentTemperature,
		MinimumTemperature: *req.MinimumTemperature,
		CurrentCapacity:    *req.CurrentCapacity,
		MinimumCapacity:    *req.MinimumCapacity,
		MaximumCapacity:    *req.MaximumCapacity,
		WarehouseId:        *req.WarehouseId,
		ProductId:          req.ProductId,
	}
}

func SectionToResponseSection(s section.Section) section.ResponseSection {
	return section.ResponseSection{
		Id:                 s.Id,
		SectionNumber:      s.SectionNumber,
		CurrentTemperature: s.CurrentTemperature,
		MinimumTemperature: s.MinimumTemperature,
		CurrentCapacity:    s.CurrentCapacity,
		MinimumCapacity:    s.MinimumCapacity,
		MaximumCapacity:    s.MaximumCapacity,
		WarehouseId:        s.WarehouseId,
		ProductId:          s.ProductId,
	}
}
