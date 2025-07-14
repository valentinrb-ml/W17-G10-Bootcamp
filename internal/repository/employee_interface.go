package repository

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) (*models.Employee, error)
	FindByCardNumberID(ctx context.Context, cardNumberID string) (*models.Employee, error)
	FindAll(ctx context.Context) ([]*models.Employee, error)
	FindByID(ctx context.Context, id int) (*models.Employee, error)
	Update(ctx context.Context, id int, e *models.Employee) (*models.Employee, error)
	Delete(ctx context.Context, id int) error
	ExistsByCardNumberID(ctx context.Context, cardNumberID string) (bool, error)
}
