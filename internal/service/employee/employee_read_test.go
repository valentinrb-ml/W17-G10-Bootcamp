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

func TestEmployeeService_Read(t *testing.T) {
	testCases := []struct {
		name        string
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		findAll     bool
		inputID     int
		wantErr     bool
		wantErrCode string
		wantLen     int
		wantID      int
	}{
		{
			name: "find_all",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				emps := testhelpers.CreateTestEmployees()
				var empsPtrs []*models.Employee
				for i := range emps {
					empsPtrs = append(empsPtrs, &emps[i])
				}
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
						return empsPtrs, nil
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
						return testhelpers.CreateExpectedEmployee(id), nil
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
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
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
					require.Nil(t, res)
				} else {
					require.NoError(t, err)
					require.NotNil(t, res)
					require.Equal(t, tc.wantID, res.ID)
				}
			}
		})
	}
}
