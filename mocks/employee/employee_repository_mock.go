package mocks

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Mock que implementa la interfaz EmployeeRepository
type EmployeeRepositoryMock struct {
	MockCreate             func(ctx context.Context, e *models.Employee) (*models.Employee, error)
	MockFindByCardNumberID func(ctx context.Context, cardNumberID string) (*models.Employee, error)
	MockFindAll            func(ctx context.Context) ([]*models.Employee, error)
	MockFindByID           func(ctx context.Context, id int) (*models.Employee, error)
	MockUpdate             func(ctx context.Context, id int, e *models.Employee) error
	MockDelete             func(ctx context.Context, id int) error
}

func (m *EmployeeRepositoryMock) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	return m.MockCreate(ctx, e)
}
func (m *EmployeeRepositoryMock) FindByCardNumberID(ctx context.Context, cardNumberID string) (*models.Employee, error) {
	return m.MockFindByCardNumberID(ctx, cardNumberID)
}
func (m *EmployeeRepositoryMock) FindAll(ctx context.Context) ([]*models.Employee, error) {
	return m.MockFindAll(ctx)
}
func (m *EmployeeRepositoryMock) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	return m.MockFindByID(ctx, id)
}
func (m *EmployeeRepositoryMock) Update(ctx context.Context, id int, e *models.Employee) error {
	return m.MockUpdate(ctx, id, e)
}
func (m *EmployeeRepositoryMock) Delete(ctx context.Context, id int) error {
	return m.MockDelete(ctx, id)
}
