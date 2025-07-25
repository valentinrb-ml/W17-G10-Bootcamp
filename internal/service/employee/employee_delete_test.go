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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestEmployeeService_Delete(t *testing.T) {
	testCases := []struct {
		name        string
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		inputID     int
		wantErr     bool
		wantErrCode string
	}{
		{
			name: "delete_ok",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						e.CardNumberID = "E005"
						return &e, nil
					},
					MockDelete: func(ctx context.Context, id int) error {
						return nil
					},
				}
			},
			inputID: 5,
			wantErr: false,
		},
		{
			name: "delete_non_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
					},
					MockDelete: func(ctx context.Context, id int) error {
						return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
					},
				}
			},
			inputID:     99,
			wantErr:     true,
			wantErrCode: apperrors.CodeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emRepo := tc.repoMock()
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)

			err := svc.Delete(context.Background(), tc.inputID)
			if tc.wantErr {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEmployeeService_Delete_extraCases(t *testing.T) {
	testCases := []struct {
		name        string
		inputID     int
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		wantErrCode string
		checkWrap   string
	}{
		{
			name:    "invalid id error",
			inputID: 0,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) { return nil, nil },
				}
			},
			wantErrCode: apperrors.CodeValidationError,
		},
		{
			name:    "repo.FindByID returns error",
			inputID: 2,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, context.DeadlineExceeded
					},
				}
			},
			checkWrap: "failed fetching employee by id",
		},
		{
			name:    "repo.FindByID returns nil, nil",
			inputID: 14,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, nil
					},
				}
			},
			wantErrCode: apperrors.CodeNotFound,
		},
		{
			name:    "repo.Delete returns error",
			inputID: 15,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return &models.Employee{ID: id, CardNumberID: "E00X"}, nil
					},
					MockDelete: func(ctx context.Context, id int) error {
						return context.Canceled
					},
				}
			},
			checkWrap: "failed deleting employee",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emRepo := tc.repoMock()
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)
			err := svc.Delete(context.Background(), tc.inputID)
			if tc.wantErrCode != "" {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
			} else if tc.checkWrap != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.checkWrap)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
