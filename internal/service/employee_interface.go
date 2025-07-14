package service

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

type EmployeeService interface {
	Create(ctx context.Context, e *models.Employee) (*models.Employee, error)
	FindAll(ctx context.Context) ([]*models.Employee, error)
	FindByID(ctx context.Context, id int) (*models.Employee, error)
	Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error)
	Delete(ctx context.Context, id int) error
}
