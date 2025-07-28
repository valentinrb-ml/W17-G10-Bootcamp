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
