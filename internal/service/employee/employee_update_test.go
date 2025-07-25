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

func TestEmployeeService_Update(t *testing.T) {
	testCases := []struct {
		name          string
		repoMock      func() *employeeMocks.EmployeeRepositoryMock
		patch         *models.EmployeePatch
		inputID       int
		wantErr       bool
		wantErrCode   string
		wantFirstName string // Para el caso exitoso
	}{
		{
			name: "update_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				updatedFirstName := "Before"
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return &models.Employee{
							ID:           id,
							CardNumberID: "E001",
							FirstName:    updatedFirstName, // Devuelve el valor actualizado
							LastName:     "Test",
							WarehouseID:  1,
						}, nil
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
			whRepo := &warehouseRepoMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)

			res, err := svc.Update(context.Background(), tc.inputID, tc.patch)
			if tc.wantErr {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
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

func strPtr(s string) *string { return &s }
