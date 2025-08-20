package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/inbound_order"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	inboundOrderMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/inbound_order"
	warehouseMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	employeeModels "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
	warehouseModels "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestInboundOrderService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		repoMock      func() *inboundOrderMocks.InboundOrderRepositoryMock
		employeeMock  func() *employeeMocks.EmployeeRepositoryMock
		warehouseMock func() *warehouseMocks.WarehouseRepositoryMock
		input         *models.InboundOrder
		wantErrCode   string
	}{
		{
			name: "create_ok",
			// Todos los pasos de validación pasan, el create es exitoso.
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) {
						return false, nil
					},
					MockCreate: func(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
						o.ID = 1
						return o, nil
					},
				}
			},
			employeeMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*employeeModels.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						return &e, nil
					},
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*warehouseModels.Warehouse, error) {
						return &warehouseModels.Warehouse{Id: id}, nil
					},
				}
			},
			input:       testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode: "", // Éxito
		},
		{
			name: "order_number_conflict",
			// El número de orden YA existe, el service debe devolver conflicto.
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) {
						return true, nil
					},
					MockCreate: func(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
						return nil, apperrors.NewAppError(apperrors.CodeConflict, "order_number already exists")
					},
				}
			},
			employeeMock:  func() *employeeMocks.EmployeeRepositoryMock { return &employeeMocks.EmployeeRepositoryMock{} },
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			input:         testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode:   apperrors.CodeConflict,
		},
		{
			name: "validation_error",
			// Falla la validación de campos requeridos
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{}
			},
			employeeMock:  func() *employeeMocks.EmployeeRepositoryMock { return &employeeMocks.EmployeeRepositoryMock{} },
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			input:         &models.InboundOrder{}, // Input vacío
			wantErrCode:   "VALIDATION_ERROR",
		},
		{
			name: "exists check error",
			// Error al consultar repositorio para unicidad
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) {
						return false, apperrors.NewAppError("INTERNAL", "db down")
					},
				}
			},
			employeeMock:  func() *employeeMocks.EmployeeRepositoryMock { return &employeeMocks.EmployeeRepositoryMock{} },
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock { return &warehouseMocks.WarehouseRepositoryMock{} },
			input:         testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode:   "INTERNAL",
		},
		{
			name: "employee_repo error",
			// Falla el repositorio de empleados al buscar empleado
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) { return false, nil },
				}
			},
			employeeMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*employeeModels.Employee, error) {
						return nil, apperrors.NewAppError("INTERNAL", "fail")
					},
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{}
			},
			input:       testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode: "INTERNAL",
		},
		{
			name: "employee not found",
			// El empleado referenciado no existe
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) { return false, nil },
				}
			},
			employeeMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*employeeModels.Employee, error) { return nil, nil },
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{}
			},
			input:       testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode: "CONFLICT",
		},
		{
			name: "warehouse find error",
			// Falla la consulta a bodega (warehouse) desde el repo de warehouse (error sql)
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) { return false, nil },
				}
			},
			employeeMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*employeeModels.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						return &e, nil
					},
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*warehouseModels.Warehouse, error) {
						return nil, apperrors.NewAppError("INTERNAL", "fail wh")
					},
				}
			},
			input:       testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode: "INTERNAL",
		},
		{
			name: "warehouse not found",
			// El warehouse no existe (devuelve error not found)
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) { return false, nil },
				}
			},
			employeeMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*employeeModels.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						return &e, nil
					},
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*warehouseModels.Warehouse, error) {
						return nil, apperrors.NewAppError("NOT_FOUND", "warehouse not found")
					},
				}
			},
			input:       testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode: "CONFLICT",
		},
		{
			name: "warehouse nil",
			// Referencia a warehouse no encontrada (devuelve nil)
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) { return false, nil },
				}
			},
			employeeMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*employeeModels.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						return &e, nil
					},
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*warehouseModels.Warehouse, error) {
						return nil, nil
					},
				}
			},
			input:       testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode: "CONFLICT",
		},
		{
			name: "create repo error",
			// Falla el insert en el repo (DB error)
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockExistsByOrderNumber: func(ctx context.Context, orderNumber string) (bool, error) { return false, nil },
					MockCreate: func(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
						return nil, apperrors.NewAppError("INTERNAL", "fail saving")
					},
				}
			},
			employeeMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*employeeModels.Employee, error) {
						e := testhelpers.CreateTestEmployee()
						e.ID = id
						return &e, nil
					},
				}
			},
			warehouseMock: func() *warehouseMocks.WarehouseRepositoryMock {
				return &warehouseMocks.WarehouseRepositoryMock{
					FuncFindById: func(ctx context.Context, id int) (*warehouseModels.Warehouse, error) {
						return &warehouseModels.Warehouse{Id: id}, nil
					},
				}
			},
			input:       testhelpers.CreateExpectedInboundOrder(0),
			wantErrCode: "INTERNAL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repoMock()
			empRepo := tc.employeeMock()
			whRepo := tc.warehouseMock()
			svc := service.NewInboundOrderService(repo, empRepo, whRepo)
			svc.SetLogger(testhelpers.NewTestLogger())

			res, err := svc.Create(context.Background(), tc.input)
			if tc.wantErrCode == "" {
				require.NoError(t, err)
				require.NotNil(t, res)
			} else {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Equal(t, tc.wantErrCode, appErr.Code)
				require.Nil(t, res)
			}
		})
	}
}
