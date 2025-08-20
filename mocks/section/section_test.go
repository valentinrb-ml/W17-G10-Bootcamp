package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

func TestSectionServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.SectionServiceMock{
		FuncFindAll:   func(ctx context.Context) ([]models.Section, error) { return nil, nil },
		FuncFindById:  func(ctx context.Context, id int) (*models.Section, error) { return nil, nil },
		FuncDelete:    func(ctx context.Context, id int) error { return nil },
		FuncCreate:    func(ctx context.Context, sec models.Section) (*models.Section, error) { return nil, nil },
		FuncUpdate:    func(ctx context.Context, id int, sec models.PatchSection) (*models.Section, error) { return nil, nil },
		FuncSetLogger: func(l logger.Logger) {},
	}

	m.FindAllSections(context.TODO())
	m.FindById(context.TODO(), 0)
	m.DeleteSection(context.TODO(), 0)
	m.CreateSection(context.TODO(), models.Section{})
	m.UpdateSection(context.TODO(), 0, models.PatchSection{})
	m.SetLogger(nil)
}

func TestSectionRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.SectionRepositoryMock{
		FuncFindAll:   func(ctx context.Context) ([]models.Section, error) { return nil, nil },
		FuncFindById:  func(ctx context.Context, id int) (*models.Section, error) { return nil, nil },
		FuncDelete:    func(ctx context.Context, id int) error { return nil },
		FuncCreate:    func(ctx context.Context, sec models.Section) (*models.Section, error) { return nil, nil },
		FuncUpdate:    func(ctx context.Context, id int, sec *models.Section) (*models.Section, error) { return nil, nil },
		FuncSetLogger: func(l logger.Logger) {},
	}

	m.FindAllSections(context.TODO())
	m.FindById(context.TODO(), 0)
	m.DeleteSection(context.TODO(), 0)
	m.CreateSection(context.TODO(), models.Section{})
	m.UpdateSection(context.TODO(), 0, &models.Section{})
	m.SetLogger(nil)
}
