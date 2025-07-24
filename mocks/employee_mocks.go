package mocks

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

type EmployeeServiceMock struct {
	MockCreate   func(ctx context.Context, e *models.Employee) (*models.Employee, error)
	MockFindAll  func(ctx context.Context) ([]*models.Employee, error)
	MockFindByID func(ctx context.Context, id int) (*models.Employee, error)
	MockUpdate   func(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error)
	MockDelete   func(ctx context.Context, id int) error
}

func (m *EmployeeServiceMock) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	return m.MockCreate(ctx, e)
}
func (m *EmployeeServiceMock) FindAll(ctx context.Context) ([]*models.Employee, error) {
	if m.MockFindAll == nil {
		return nil, nil
	}
	return m.MockFindAll(ctx)
}
func (m *EmployeeServiceMock) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	if m.MockFindByID == nil {
		return nil, nil
	}
	return m.MockFindByID(ctx, id)
}
func (m *EmployeeServiceMock) Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
	if m.MockUpdate == nil {
		return nil, nil
	}
	return m.MockUpdate(ctx, id, patch)
}
func (m *EmployeeServiceMock) Delete(ctx context.Context, id int) error {
	if m.MockDelete == nil {
		return nil
	}
	return m.MockDelete(ctx, id)
}
