package mappers

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

func RequestSectionToSection(req models.PostSection) models.Section {
	return models.Section{
		SectionNumber:      req.SectionNumber,
		CurrentTemperature: *req.CurrentTemperature,
		MinimumTemperature: *req.MinimumTemperature,
		CurrentCapacity:    req.CurrentCapacity,
		MinimumCapacity:    req.MinimumCapacity,
		MaximumCapacity:    req.MaximumCapacity,
		WarehouseId:        req.WarehouseId,
		ProductTypeId:      req.ProductTypeId,
	}
}

func SectionToResponseSection(s models.Section) models.ResponseSection {
	return models.ResponseSection{
		Id:                 s.Id,
		SectionNumber:      s.SectionNumber,
		CurrentTemperature: s.CurrentTemperature,
		MinimumTemperature: s.MinimumTemperature,
		CurrentCapacity:    s.CurrentCapacity,
		MinimumCapacity:    s.MinimumCapacity,
		MaximumCapacity:    s.MaximumCapacity,
		WarehouseId:        s.WarehouseId,
		ProductTypeId:      s.ProductTypeId,
	}
}

func ApplySectionPatch(sec models.PatchSection, existing *models.Section) {
	if sec.ProductTypeId != nil {
		existing.ProductTypeId = *sec.ProductTypeId

	}
	if sec.SectionNumber != nil {
		existing.SectionNumber = *sec.SectionNumber
	}
	if sec.CurrentTemperature != nil {
		existing.CurrentTemperature = *sec.CurrentTemperature
	}
	if sec.MinimumTemperature != nil {
		existing.MinimumTemperature = *sec.MinimumTemperature
	}
	if sec.CurrentCapacity != nil {
		existing.CurrentCapacity = *sec.CurrentCapacity
	}
	if sec.MinimumCapacity != nil {
		existing.MinimumCapacity = *sec.MinimumCapacity
	}
	if sec.MaximumCapacity != nil {
		existing.MaximumCapacity = *sec.MaximumCapacity
	}
	if sec.WarehouseId != nil {
		existing.WarehouseId = *sec.WarehouseId
	}
}
