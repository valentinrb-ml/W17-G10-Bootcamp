package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	warehouseMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
	wmodels "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestEmployeeService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		repoMock      func() *employeeMocks.EmployeeRepositoryMock
		warehouseMock func() *warehouseMocks.WarehouseRepositoryMock
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
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
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
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
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

func TestEmployeeService_Create_extraCases(t *testing.T) {
	testCases := []struct {
		name        string
		input       *models.Employee
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		whMock      func() *warehouseMocks.WarehouseRepositoryMock
		wantErrCode string // "" si se espera success
		checkWrap   string // mensaje si se espera error envuelto
	}{
		{
			name:  "validation_error",
			input: &models.Employee{},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{}
			},
			whMock: func() *warehouseMocks.WarehouseRepositoryMock { // won't be called
				return &warehouseMocks.WarehouseRepositoryMock{}
			},
			wantErrCode: apperrors.CodeValidationError,
		},
		{
			name: "warehouse find generic error",
			input: &models.Employee{
				CardNumberID: "X", FirstName: "A", LastName: "B", WarehouseID: 1,
			},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{}
			},
			whMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return nil, context.DeadlineExceeded
					},
				}
			},
			checkWrap: "failed getting warehouse by id",
		},
		{
			name: "warehouse not exists (not found code)",
			input: &models.Employee{
				CardNumberID: "X", FirstName: "A", LastName: "B", WarehouseID: 1,
			},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{}
			},
			whMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "no warehouse")
					},
				}
			},
			wantErrCode: apperrors.CodeBadRequest,
		},
		{
			name: "warehouse nil, no error",
			input: &models.Employee{
				CardNumberID: "X", FirstName: "A", LastName: "B", WarehouseID: 1,
			},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{}
			},
			whMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return nil, nil
					},
				}
			},
			wantErrCode: apperrors.CodeBadRequest,
		},
		{
			name: "findByCardNumberID returns err (generic)",
			input: &models.Employee{
				CardNumberID: "X", FirstName: "A", LastName: "B", WarehouseID: 1,
			},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByCardNumberID: func(ctx context.Context, cardNumberID string) (*models.Employee, error) {
						return nil, context.DeadlineExceeded
					},
				}
			},
			whMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return &wmodels.Warehouse{Id: 1}, nil
					},
				}
			},
			checkWrap: "failed checking card_number_id",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewEmployeeDefault(tc.repoMock(), tc.whMock())
			res, err := svc.Create(context.Background(), tc.input)
			if tc.wantErrCode != "" {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
				require.Nil(t, res)
			} else if tc.checkWrap != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.checkWrap)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
			}
		})
	}
}
