package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/inbound_order"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	inboundOrderMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/inbound_order"
	warehouseMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Test unitario para el método Report del servicio de inbound orders,
// cubriéndose los caminos de branch global (todos los empleados) y por empleado específico.
func TestInboundOrderService_Report(t *testing.T) {
	testCases := []struct {
		name       string
		repoMock   func() *inboundOrderMocks.InboundOrderRepositoryMock
		employeeID *int // nil para ReportAll, con valor para ReportByID
		wantEmpty  bool // espera reporte vacío
		wantErr    bool
	}{
		{
			name: "report_all_ok",
			// Cuando employeeID es nil, se espera un reporte "global".
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					// Usa el helper para un listado dummy
					MockReportAll: func(ctx context.Context) ([]models.InboundOrderReport, error) {
						return testhelpers.CreateInboundOrderReports(), nil
					},
				}
			},
			employeeID: nil,   // Simula llamada general
			wantEmpty:  false, // Debe haber datos
			wantErr:    false,
		},
		{
			name: "report_by_id_ok",
			// Cuando employeeID está presente, se espera un solo reporte.
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockReportByID: func(ctx context.Context, id int) (*models.InboundOrderReport, error) {
						r := testhelpers.CreateInboundOrderReport(id)
						return &r, nil // Simula encontrado
					},
				}
			},
			employeeID: func() *int { v := 1; return &v }(), // busca un id específico
			wantEmpty:  false,
			wantErr:    false,
		},
		{
			name: "report_by_id_not_found",
			// Simula el camino donde no hay empleado
			repoMock: func() *inboundOrderMocks.InboundOrderRepositoryMock {
				return &inboundOrderMocks.InboundOrderRepositoryMock{
					MockReportByID: func(ctx context.Context, id int) (*models.InboundOrderReport, error) {
						return nil, nil // Simula "no encontrado"
					},
				}
			},
			employeeID: func() *int { v := 999; return &v }(),
			wantEmpty:  true,  // Espera que no haya datos (nil)
			wantErr:    false, // Aquí no es error, solo no hay datos
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Crea los mocks de Repos necesarios (los otros no relevantes pueden estar vacíos/dummy)
			repo := tc.repoMock()
			empRepo := &employeeMocks.EmployeeRepositoryMock{}
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
			svc := service.NewInboundOrderService(repo, empRepo, whRepo)
			svc.SetLogger(testhelpers.NewTestLogger())

			// Ejecuta el método Report
			result, err := svc.Report(context.Background(), tc.employeeID)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.employeeID == nil {
					// Caso branch global: debe ser un slice (list)
					reports, ok := result.([]models.InboundOrderReport)
					require.True(t, ok)
					if tc.wantEmpty {
						require.Empty(t, reports)
					} else {
						require.NotEmpty(t, reports)
					}
				} else {
					// Caso por empleado: debe ser pointer o nil
					rep, ok := result.(*models.InboundOrderReport)
					require.True(t, ok)
					if tc.wantEmpty {
						require.Nil(t, rep)
					} else {
						require.NotNil(t, rep)
					}
				}
			}
		})
	}
}
