package mocks

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionServiceMock struct {
	FuncFindAll   func(ctx context.Context) ([]models.Section, error)
	FuncFindById  func(ctx context.Context, id int) (*models.Section, error)
	FuncDelete    func(ctx context.Context, id int) error
	FuncCreate    func(ctx context.Context, sec models.Section) (*models.Section, error)
	FuncUpdate    func(ctx context.Context, id int, sec models.PatchSection) (*models.Section, error)
	FuncSetLogger func(l logger.Logger)
}

func (m *SectionServiceMock) FindAllSections(ctx context.Context) ([]models.Section, error) {
	return m.FuncFindAll(ctx)
}

func (m *SectionServiceMock) FindById(ctx context.Context, id int) (*models.Section, error) {
	return m.FuncFindById(ctx, id)
}

func (m *SectionServiceMock) DeleteSection(ctx context.Context, id int) error {
	return m.FuncDelete(ctx, id)
}

func (m *SectionServiceMock) CreateSection(ctx context.Context, sec models.Section) (*models.Section, error) {
	return m.FuncCreate(ctx, sec)
}
func (m *SectionServiceMock) UpdateSection(ctx context.Context, id int, sec models.PatchSection) (*models.Section, error) {
	return m.FuncUpdate(ctx, id, sec)
}

func (m *SectionServiceMock) SetLogger(l logger.Logger) {
	m.FuncSetLogger(l)
}
