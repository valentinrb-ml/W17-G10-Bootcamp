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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func strPtr(s string) *string { return &s }

func TestEmployeeService_Update(t *testing.T) {
	testCases := []struct {
		name          string
		repoMock      func() *employeeMocks.EmployeeRepositoryMock
		patch         *models.EmployeePatch
		inputID       int
		wantErr       bool
		wantErrCode   string
		wantFirstName string
	}{
		{
			name: "update_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				// Usa helper:
				base := testhelpers.CreateTestEmployee()
				updatedFirstName := base.FirstName
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						emp := base
						emp.ID = id
						emp.FirstName = updatedFirstName
						return &emp, nil
					},
					MockUpdate: func(ctx context.Context, id int, e *models.Employee) error {
						updatedFirstName = e.FirstName // Simula persistir el cambio
						return nil
					},
				}
			},
			inputID:       1,
			patch:         &models.EmployeePatch{FirstName: strPtr("After")},
			wantErr:       false,
			wantFirstName: "After",
		},
		{
			name: "update_non_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
					},
					MockUpdate: func(ctx context.Context, id int, e *models.Employee) error {
						return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
					},
				}
			},
			inputID:     99,
			patch:       &models.EmployeePatch{FirstName: strPtr("Nada")},
			wantErr:     true,
			wantErrCode: apperrors.CodeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emRepo := tc.repoMock()
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)

			res, err := svc.Update(context.Background(), tc.inputID, tc.patch)
			if tc.wantErr {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				if tc.wantFirstName != "" && res.FirstName != tc.wantFirstName {
					t.Errorf("expected FirstName %s got %s", tc.wantFirstName, res.FirstName)
				}
			}
		})
	}
}

func TestEmployeeService_Update_extraCases(t *testing.T) {
	val := 8
	testCases := []struct {
		name        string
		id          int
		patch       *models.EmployeePatch
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		whMock      func() *warehouseMocks.WarehouseRepositoryMock
		wantErrCode string // usa "" y checkWrap cuando se espera error no AppError
		checkWrap   string
	}{
		{
			name:  "patch validation error",
			id:    1,
			patch: &models.EmployeePatch{},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) { return nil, nil },
				}
			},
			whMock:      func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			wantErrCode: apperrors.CodeValidationError,
		},
		{
			name:  "id invalid (<=0)",
			id:    0,
			patch: &models.EmployeePatch{FirstName: strPtr("Test")},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) { return nil, nil },
				}
			},
			whMock:      func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			wantErrCode: apperrors.CodeValidationError,
		},
		{
			name:  "FindByID returns error",
			id:    2,
			patch: &models.EmployeePatch{FirstName: strPtr("Q")},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, context.DeadlineExceeded
					},
				}
			},
			whMock:    func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			checkWrap: "failed fetching employee by id",
		},
		{
			name:  "CardNumberID in use by another (conflict)",
			id:    1,
			patch: &models.EmployeePatch{CardNumberID: strPtr("NEWID")},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = 1
						e.CardNumberID = "C"
						return &e, nil
					},
					MockFindByCardNumberID: func(ctx context.Context, cardNumberID string) (*models.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = 99
						e.CardNumberID = cardNumberID
						return &e, nil
					},
					MockUpdate: func(ctx context.Context, id int, e *models.Employee) error { return nil },
				}
			},
			whMock:      func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			wantErrCode: apperrors.CodeConflict,
		},
		{
			name:  "warehouse find returns generic error",
			id:    3,
			patch: &models.EmployeePatch{WarehouseID: &val},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						e.CardNumberID = "A"
						return &e, nil
					},
					MockUpdate: func(ctx context.Context, id int, emp *models.Employee) error { return nil },
				}
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
			name:  "warehouse not found (not found code)",
			id:    4,
			patch: &models.EmployeePatch{WarehouseID: &val},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						e.CardNumberID = "A"
						return &e, nil
					},
					MockUpdate: func(ctx context.Context, id int, emp *models.Employee) error { return nil },
				}
			},
			whMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "no wh")
					},
				}
			},
			wantErrCode: apperrors.CodeBadRequest,
		},
		{
			name:  "warehouse nil, no error",
			id:    5,
			patch: &models.EmployeePatch{WarehouseID: &val},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						e.CardNumberID = "A"
						return &e, nil
					},
					MockUpdate: func(ctx context.Context, id int, emp *models.Employee) error { return nil },
				}
			},
			whMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) { return nil, nil },
				}
			},
			wantErrCode: apperrors.CodeBadRequest,
		},
		{
			name:  "repo.Update returns error",
			id:    6,
			patch: &models.EmployeePatch{FirstName: strPtr("T")},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						e.CardNumberID = "X"
						return &e, nil
					},
					MockUpdate: func(ctx context.Context, id int, emp *models.Employee) error {
						return context.Canceled
					},
				}
			},
			whMock:    func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			checkWrap: "failed updating employee",
		},
		{
			name:  "repo.FindByID after update returns error",
			id:    7,
			patch: &models.EmployeePatch{FirstName: strPtr("T")},
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				step := 0
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						step++
						if step == 2 {
							return nil, context.DeadlineExceeded
						}
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						e.CardNumberID = "X"
						return &e, nil
					},
					MockUpdate: func(ctx context.Context, id int, emp *models.Employee) error { return nil },
				}
			},
			whMock:    func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			checkWrap: "failed fetching employee after update",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewEmployeeDefault(tc.repoMock(), tc.whMock())
			res, err := svc.Update(context.Background(), tc.id, tc.patch)
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
			}
		})
	}
}
