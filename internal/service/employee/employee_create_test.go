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

// Tests de error y trayectorias no triviales en Employee Service para coverage alto
func TestEmployeeService_Create_allErrors(t *testing.T) {
	ctx := context.Background()
	whRepo := &warehouseMocks.WarehouseRepositoryMock{}

	t.Run("validation error (Empty last name)", func(t *testing.T) {
		repo := &employeeMocks.EmployeeRepositoryMock{}
		svc := service.NewEmployeeDefault(repo, whRepo)
		in := &models.Employee{CardNumberID: "C", FirstName: "X", LastName: "", WarehouseID: 1}
		res, err := svc.Create(ctx, in)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("warehouse FindById returns not found", func(t *testing.T) {
		repo := &employeeMocks.EmployeeRepositoryMock{}
		whRepo := &warehouseMocks.WarehouseRepositoryMock{
			FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "no wh")
			},
		}
		svc := service.NewEmployeeDefault(repo, whRepo)
		e := testhelpers.CreateTestEmployee()
		res, err := svc.Create(ctx, &e)
		require.Error(t, err)
		require.Nil(t, res)
		require.Contains(t, err.Error(), "warehouse_id does not exist")
	})

	t.Run("warehouse FindById returns other error", func(t *testing.T) {
		repo := &employeeMocks.EmployeeRepositoryMock{}
		whRepo := &warehouseMocks.WarehouseRepositoryMock{
			FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
				return nil, context.DeadlineExceeded
			},
		}
		svc := service.NewEmployeeDefault(repo, whRepo)
		e := testhelpers.CreateTestEmployee()
		res, err := svc.Create(ctx, &e)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed getting warehouse by id")
		require.Nil(t, res)
	})

	t.Run("warehouse FindById returns nil, no error", func(t *testing.T) {
		repo := &employeeMocks.EmployeeRepositoryMock{}
		whRepo := &warehouseMocks.WarehouseRepositoryMock{
			FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) {
				return nil, nil
			},
		}
		svc := service.NewEmployeeDefault(repo, whRepo)
		e := testhelpers.CreateTestEmployee()
		res, err := svc.Create(ctx, &e)
		require.Error(t, err)
		require.Contains(t, err.Error(), "warehouse_id does not exist")
		require.Nil(t, res)
	})

	t.Run("repo.FindByCardNumberID returns error", func(t *testing.T) {
		repo := &employeeMocks.EmployeeRepositoryMock{
			MockFindByCardNumberID: func(ctx context.Context, cardNumberID string) (*models.Employee, error) {
				return nil, context.DeadlineExceeded
			},
		}
		whRepo := &warehouseMocks.WarehouseRepositoryMock{
			FuncFindById: func(ctx context.Context, id int) (*wmodels.Warehouse, error) { return &wmodels.Warehouse{Id: 1}, nil },
		}
		svc := service.NewEmployeeDefault(repo, whRepo)
		e := testhelpers.CreateTestEmployee()
		res, err := svc.Create(ctx, &e)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed checking card_number_id uniqueness")
		require.Nil(t, res)
	})
}
