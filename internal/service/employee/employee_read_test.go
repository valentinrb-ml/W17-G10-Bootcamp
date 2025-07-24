package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
	wmodels "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

// Mock warehouse repo mínimo solo para compilar el service
type warehouseRepoMock struct {
	MockCreate   func(ctx context.Context, w wmodels.Warehouse) (*wmodels.Warehouse, error)
	MockFindAll  func(ctx context.Context) ([]wmodels.Warehouse, error)
	MockFindById func(ctx context.Context, id int) (*wmodels.Warehouse, error)
	MockUpdate   func(ctx context.Context, id int, w wmodels.Warehouse) (*wmodels.Warehouse, error)
	MockDelete   func(ctx context.Context, id int) error
}

func (m *warehouseRepoMock) Create(ctx context.Context, w wmodels.Warehouse) (*wmodels.Warehouse, error) {
	if m.MockCreate != nil {
		return m.MockCreate(ctx, w)
	}
	return nil, nil
}
func (m *warehouseRepoMock) FindAll(ctx context.Context) ([]wmodels.Warehouse, error) {
	if m.MockFindAll != nil {
		return m.MockFindAll(ctx)
	}
	return nil, nil
}
func (m *warehouseRepoMock) FindById(ctx context.Context, id int) (*wmodels.Warehouse, error) {
	if m.MockFindById != nil {
		return m.MockFindById(ctx, id)
	}
	return nil, nil
}
func (m *warehouseRepoMock) Update(ctx context.Context, id int, w wmodels.Warehouse) (*wmodels.Warehouse, error) {
	if m.MockUpdate != nil {
		return m.MockUpdate(ctx, id, w)
	}
	return nil, nil
}
func (m *warehouseRepoMock) Delete(ctx context.Context, id int) error {
	if m.MockDelete != nil {
		return m.MockDelete(ctx, id)
	}
	return nil
}

func TestEmployeeService_Read(t *testing.T) {
	testCases := []struct {
		name        string
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		findAll     bool
		inputID     int // para find_by_id casos
		wantErr     bool
		wantErrCode string
		wantLen     int // Para find_all caso feliz
		wantID      int // Para find_by_id_existente
	}{
		{
			name: "find_all",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
						return []*models.Employee{
							{ID: 1, CardNumberID: "E001", FirstName: "Lucas", LastName: "Martinez", WarehouseID: 1},
							{ID: 2, CardNumberID: "E002", FirstName: "Sonia", LastName: "Lopez", WarehouseID: 2},
						}, nil
					},
				}
			},
			findAll: true,
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "find_by_id_non_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
					},
				}
			},
			findAll:     false,
			inputID:     99,
			wantErr:     true,
			wantErrCode: apperrors.CodeNotFound,
		},
		{
			name: "find_by_id_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return &models.Employee{
							ID: id, CardNumberID: "E015", FirstName: "Maria", LastName: "Perez", WarehouseID: 3,
						}, nil
					},
				}
			},
			findAll: false,
			inputID: 15,
			wantErr: false,
			wantID:  15,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emRepo := tc.repoMock()
			whRepo := &warehouseRepoMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)

			if tc.findAll {
				result, err := svc.FindAll(context.Background())
				require.False(t, tc.wantErr)
				require.NoError(t, err)
				require.Len(t, result, tc.wantLen)
			} else {
				res, err := svc.FindByID(context.Background(), tc.inputID)
				if tc.wantErr {
					require.Error(t, err)
					appErr, ok := err.(*apperrors.AppError)
					require.True(t, ok)
					require.Equal(t, tc.wantErrCode, appErr.Code)
				} else {
					require.NoError(t, err)
					require.NotNil(t, res)
					require.Equal(t, tc.wantID, res.ID)
				}
			}
		})
	}
}
