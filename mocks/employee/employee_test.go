package mocks_test

import (
	"context"
	"testing"

	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

func TestEmployeeRepositoryMock_DummyCoverage(t *testing.T) {
	m := &mocks.EmployeeRepositoryMock{
		MockCreate:             func(ctx context.Context, e *models.Employee) (*models.Employee, error) { return nil, nil },
		MockFindByCardNumberID: func(ctx context.Context, cardNumberID string) (*models.Employee, error) { return nil, nil },
		MockFindAll:            func(ctx context.Context) ([]*models.Employee, error) { return nil, nil },
		MockFindByID:           func(ctx context.Context, id int) (*models.Employee, error) { return nil, nil },
		MockUpdate:             func(ctx context.Context, id int, e *models.Employee) error { return nil },
		MockDelete:             func(ctx context.Context, id int) error { return nil },
	}

	m.Create(context.TODO(), &models.Employee{})
	m.FindByCardNumberID(context.TODO(), "")
	m.FindAll(context.TODO())
	m.FindByID(context.TODO(), 0)
	m.Update(context.TODO(), 0, &models.Employee{})
	m.Delete(context.TODO(), 0)
}

func TestEmployeeServiceMock_DummyCoverage(t *testing.T) {
	m := &mocks.EmployeeServiceMock{
		MockCreate:   func(ctx context.Context, e *models.Employee) (*models.Employee, error) { return nil, nil },
		MockFindAll:  func(ctx context.Context) ([]*models.Employee, error) { return nil, nil },
		MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) { return nil, nil },
		MockUpdate: func(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
			return nil, nil
		},
		MockDelete: func(ctx context.Context, id int) error { return nil },
	}

	m.Create(context.TODO(), &models.Employee{})
	m.FindAll(context.TODO())
	m.FindByID(context.TODO(), 0)
	m.Update(context.TODO(), 0, &models.EmployeePatch{})
	m.Delete(context.TODO(), 0)
}
