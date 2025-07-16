package repository

import (
	"context"
	"database/sql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

// SectionRepository is an interface that represents a section repository
type SectionRepository interface {
	FindAllSections(ctx context.Context) ([]models.Section, error)
	FindById(ctx context.Context, id int) (*models.Section, error)
	DeleteSection(ctx context.Context, id int) error
	CreateSection(ctx context.Context, sec models.Section) (*models.Section, error)
	UpdateSection(ctx context.Context, id int, sec *models.Section) (*models.Section, error)
}

// sectionRepository implements SectionRepository using MySQL as the data source.
type sectionRepository struct {
	mysql *sql.DB
}

// NewSectionMap is a function that returns a new instance of SectionMap
func NewSectionRepository(db *sql.DB) SectionRepository {
	return &sectionRepository{db}
}
