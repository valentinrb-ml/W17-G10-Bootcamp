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

// Test para la lectura de empleados en el service (FindAll y FindByID) usando helpers centralizados
func TestEmployeeService_Read(t *testing.T) {
	testCases := []struct {
		name        string
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		findAll     bool
		inputID     int
		wantErr     bool
		wantErrCode string
		wantLen     int // cantidad esperada para FindAll
		wantID      int // id esperado para FindByID
	}{
		{
			name: "find_all",
			// Simula el repo devolviendo empleados usando el helper (DRY)
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				emps := testhelpers.CreateTestEmployees() // Slice de empleados dummy
				var empsPtrs []*models.Employee
				for i := range emps {
					empsPtrs = append(empsPtrs, &emps[i]) // Pasa a []*models.Employee
				}
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
						return empsPtrs, nil // Devuelve todos los empleados "de la BD"
					},
				}
			},
			findAll: true,
			wantErr: false,
			wantLen: 2, // Según los dummy del helper
		},
		{
			name: "find_by_id_non_existent",
			// Simula repo que devuelve error not found
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
					},
				}
			},
			findAll:     false,
			inputID:     99, // id que no existe
			wantErr:     true,
			wantErrCode: apperrors.CodeNotFound,
		},
		{
			name: "find_by_id_existent",
			// Usa el helper para definir el empleado encontrado
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return testhelpers.CreateExpectedEmployee(id), nil
					},
				}
			},
			findAll: false,
			inputID: 15, // id esperado
			wantErr: false,
			wantID:  15, // id esperado en la respuesta
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Instancia los mocks necesarios para service
			emRepo := tc.repoMock()
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)

			if tc.findAll {
				// Test para FindAll: checa tamaño y no error
				result, err := svc.FindAll(context.Background())
				require.False(t, tc.wantErr)
				require.NoError(t, err)
				require.Len(t, result, tc.wantLen)
			} else {
				// Test para FindByID: checa error si corresponde, y valor si corresponde
				res, err := svc.FindByID(context.Background(), tc.inputID)
				if tc.wantErr {
					require.Error(t, err)
					appErr, ok := err.(*apperrors.AppError)
					require.True(t, ok)
					require.Equal(t, tc.wantErrCode, appErr.Code)
					require.Nil(t, res) // Para el caso de not found, el resultado debe ser nil
				} else {
					require.NoError(t, err)
					require.NotNil(t, res)
					require.Equal(t, tc.wantID, res.ID)
				}
			}
		})
	}
}
func TestEmployeeService_FindAll_error(t *testing.T) {
	repo := &employeeMocks.EmployeeRepositoryMock{
		MockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
			return nil, context.DeadlineExceeded
		},
	}
	whRepo := &warehouseMocks.WarehouseRepositoryMock{}
	svc := service.NewEmployeeDefault(repo, whRepo)
	res, err := svc.FindAll(context.Background())
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "failed fetching all employees")
}
func TestEmployeeService_FindByID_errors(t *testing.T) {
	repo := &employeeMocks.EmployeeRepositoryMock{}
	whRepo := &warehouseMocks.WarehouseRepositoryMock{}
	svc := service.NewEmployeeDefault(repo, whRepo)
	svc.SetLogger(testhelpers.NewTestLogger())

	// ID inválido (<=0)
	res, err := svc.FindByID(context.Background(), 0)
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "id must be positive")

	// repo.FindByID devuelve error
	repo2 := &employeeMocks.EmployeeRepositoryMock{
		MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
			return nil, context.DeadlineExceeded
		},
	}
	svc2 := service.NewEmployeeDefault(repo2, whRepo)
	res, err = svc2.FindByID(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed fetching employee by id")
	require.Nil(t, res)

	// repo.FindByID no encuentra (nil,nil)
	repo3 := &employeeMocks.EmployeeRepositoryMock{
		MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
			return nil, nil
		},
	}
	svc3 := service.NewEmployeeDefault(repo3, whRepo)
	res, err = svc3.FindByID(context.Background(), 9)
	require.Error(t, err)
	require.Contains(t, err.Error(), "employee not found")
	require.Nil(t, res)
}
