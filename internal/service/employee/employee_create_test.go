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

type warehouseRepoMockInline struct {
	MockFindById func(ctx context.Context, id int) (*wmodels.Warehouse, error)
}

func (m *warehouseRepoMockInline) FindById(ctx context.Context, id int) (*wmodels.Warehouse, error) {
	return m.MockFindById(ctx, id)
}
func (m *warehouseRepoMockInline) Create(ctx context.Context, w wmodels.Warehouse) (*wmodels.Warehouse, error) {
	return nil, nil
}
func (m *warehouseRepoMockInline) FindAll(ctx context.Context) ([]wmodels.Warehouse, error) {
	return nil, nil
}
func (m *warehouseRepoMockInline) Update(ctx context.Context, id int, w wmodels.Warehouse) (*wmodels.Warehouse, error) {
	return nil, nil
}
func (m *warehouseRepoMockInline) Delete(ctx context.Context, id int) error {
	return nil
}

func TestEmployeeService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		repoMock      func() *employeeMocks.EmployeeRepositoryMock
		warehouseMock func() *warehouseRepoMockInline
		input         *models.Employee
		wantErrCode   string
	}{
		{
			name: "create_ok",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByCardNumberID: func(ctx context.Context, cardNumberID string) (*models.Employee, error) {
						return nil, nil
					},
					MockCreate: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
						e.ID = 1
						return e, nil
					},
				}
			},
			warehouseMock: func() *warehouseRepoMockInline {
				return &warehouseRepoMockInline{
					MockFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return &wmodels.Warehouse{Id: id}, nil // Warehouse existe
					},
				}
			},
			input: &models.Employee{
				CardNumberID: "E001",
				FirstName:    "Paola",
				LastName:     "Lopez",
				WarehouseID:  1,
			},
			wantErrCode: "",
		},
		{
			name: "create_conflict",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByCardNumberID: func(ctx context.Context, cardNumberID string) (*models.Employee, error) {
						return &models.Employee{ID: 99, CardNumberID: cardNumberID}, nil
					},
					MockCreate: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
						return nil, apperrors.NewAppError(apperrors.CodeConflict, "card_number_id already exists")
					},
				}
			},
			warehouseMock: func() *warehouseRepoMockInline {
				return &warehouseRepoMockInline{
					MockFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return &wmodels.Warehouse{Id: id}, nil // Warehouse existe
					},
				}
			},
			input: &models.Employee{
				CardNumberID: "E001",
				FirstName:    "Lucia",
				LastName:     "Soler",
				WarehouseID:  1,
			},
			wantErrCode: apperrors.CodeConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emRepo := tc.repoMock()
			whRepo := tc.warehouseMock()
			svc := service.NewEmployeeDefault(emRepo, whRepo)

			res, err := svc.Create(context.Background(), tc.input)

			if tc.wantErrCode == "" {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, tc.input.CardNumberID, res.CardNumberID)
			} else {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
			}
		})
	}
}
