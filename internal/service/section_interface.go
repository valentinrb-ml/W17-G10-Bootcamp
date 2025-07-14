package service

import (
	"context"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

type SectionService interface {
	FindAllSections(ctx context.Context) ([]section.Section, error)
	FindById(ctx context.Context, id int) (*section.Section, error)
	DeleteSection(ctx context.Context, id int) error
	CreateSection(ctx context.Context, sec section.Section) (*section.Section, error)
	UpdateSection(ctx context.Context, id int, sec section.PatchSection) (*section.Section, error)
}
