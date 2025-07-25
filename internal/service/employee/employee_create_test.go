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
				// Mock del repositorio: simula que el cardNumberID NO existe y luego crea el empleado
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByCardNumberID: func(ctx context.Context, cardNumberID string) (*models.Employee, error) {
						return nil, nil
					},
					MockCreate: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
						e.ID = 1 // simula el autoincremento de la db
						return e, nil
					},
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				// Simula que el warehouse existe
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
						return &wmodels.Warehouse{Id: id}, nil // Warehouse existe
					},
				}
			},
			// Usamos el helper, solo sobrescribiendo los valores relevantes para este test
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				// Setea los valores deseados para el test:
				e.CardNumberID = "E001"
				e.FirstName = "Paola"
				e.LastName = "Lopez"
				e.WarehouseID = 1
				return &e
			}(),
			wantErrCode: "",
		},
		{
			name: "create_conflict",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				// Simula que el cardNumberID YA existe, así que Create nunca es invocado
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
			// Usamos de nuevo el helper, cambiando solo los valores relevantes para este caso
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				e.CardNumberID = "E001"
				e.FirstName = "Lucia"
				e.LastName = "Soler"
				e.WarehouseID = 1
				return &e
			}(),
			wantErrCode: apperrors.CodeConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Instancia los mocks para repo y warehouse según el caso
			emRepo := tc.repoMock()
			whRepo := tc.warehouseMock()
			svc := service.NewEmployeeDefault(emRepo, whRepo)
			// Ejecuta el método Create del service con el input armado con helper
			res, err := svc.Create(context.Background(), tc.input)

			if tc.wantErrCode == "" {
				// Caso de éxito: no hay error, el empleado creado debe coincidir con el input
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, tc.input.CardNumberID, res.CardNumberID)
			} else {
				// Caso de conflicto: debe retornar error del tipo esperado
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
			}
		})
	}
}
