package mocks

import (
	"context"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionRepositoryMock struct {
	FuncFindAll  func(ctx context.Context) ([]models.Section, error)
	FuncFindById func(ctx context.Context, id int) (*models.Section, error)
	FuncDelete   func(ctx context.Context, id int) error
	FuncCreate   func(ctx context.Context, sec models.Section) (*models.Section, error)
	FuncUpdate   func(ctx context.Context, id int, sec *models.Section) (*models.Section, error)
}

func (m *SectionRepositoryMock) FindAllSections(ctx context.Context) ([]models.Section, error) {
	return m.FuncFindAll(ctx)
}

func (m *SectionRepositoryMock) FindById(ctx context.Context, id int) (*models.Section, error) {
	return m.FuncFindById(ctx, id)
}

func (m *SectionRepositoryMock) DeleteSection(ctx context.Context, id int) error {
	return m.FuncDelete(ctx, id)
}

func (m *SectionRepositoryMock) CreateSection(ctx context.Context, sec models.Section) (*models.Section, error) {
	return m.FuncCreate(ctx, sec)
}
func (m *SectionRepositoryMock) UpdateSection(ctx context.Context, id int, sec *models.Section) (*models.Section, error) {
	return m.FuncUpdate(ctx, id, sec)
}

// Mover luego a helpers
func DummySection(id int) models.Section {
	return models.Section{
		Id:                 id,
		SectionNumber:      id*10 + 1,
		CurrentCapacity:    20,
		CurrentTemperature: 5,
		MaximumCapacity:    100,
		MinimumCapacity:    10,
		MinimumTemperature: 5,
		ProductTypeId:      2,
		WarehouseId:        1,
	}
}
func DummySectionPatch(section models.Section) models.PatchSection {
	return models.PatchSection{
		SectionNumber:      &section.SectionNumber,
		CurrentCapacity:    &section.CurrentCapacity,
		CurrentTemperature: &section.CurrentTemperature,
		MaximumCapacity:    &section.MaximumCapacity,
		MinimumCapacity:    &section.MinimumCapacity,
		MinimumTemperature: &section.MinimumTemperature,
		ProductTypeId:      &section.ProductTypeId,
		WarehouseId:        &section.WarehouseId,
	}
}
