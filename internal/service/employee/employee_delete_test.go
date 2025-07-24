package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
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
						return &models.Employee{ID: id, CardNumberID: "E005"}, nil
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
			whRepo := &warehouseRepoMock{}
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
