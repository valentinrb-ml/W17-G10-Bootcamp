package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/inbound_order"
	inboundOrderMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Test para GET /api/v1/employees/reportInboundOrders con distintos escenarios usando helpers/mocks
func TestInboundOrderHandler_Report(t *testing.T) {
	testCases := []struct {
		name        string
		queryID     *int
		mockReport  func(ctx context.Context, id *int) (interface{}, error)
		wantStatus  int
		wantContent string
	}{
		{
			name:    "report_ok",
			queryID: func() *int { i := 1; return &i }(), // Simula llamado con ?id=1
			mockReport: func(ctx context.Context, id *int) (interface{}, error) {
				// Devuelve el slice de reportes como interface{}
				return testhelpers.CreateInboundOrderReports(), nil
			},
			wantStatus:  http.StatusOK,
			wantContent: `"inbound_orders_count":5`, // Espera campo correcto en JSON
		},
		{
			name:    "report_empty",
			queryID: nil, // Simula llamado SIN id param (?id=)
			// Devuelve un slice vacío (indicando ningún reporte)
			mockReport: func(ctx context.Context, id *int) (interface{}, error) {
				return []models.InboundOrderReport{}, nil
			},
			wantStatus:  http.StatusOK,
			wantContent: `"data":[]`, // JSON esperado para lista vacía
		},
		{
			name:    "report_not_found",
			queryID: func() *int { i := 999; return &i }(), // id que no existe
			mockReport: func(ctx context.Context, id *int) (interface{}, error) {
				// Simula error not found
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
			},
			wantStatus:  http.StatusNotFound,
			wantContent: `"not found"`,
		},
		{
			name:    "report_all_no_id",
			queryID: nil, // Sin parámetro id
			mockReport: func(ctx context.Context, id *int) (interface{}, error) {
				// Simula que devuelve varios reportes al no pasar id
				return testhelpers.CreateInboundOrderReports(), nil
			},
			wantStatus:  http.StatusOK,
			wantContent: `"inbound_orders_count":5`, // cualquier campo representativo de los datos
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock de servicio usando el closure de la tabla de test
			mockSvc := &inboundOrderMocks.InboundOrderServiceMock{MockReport: tc.mockReport}
			h := handler.NewInboundOrderHandler(mockSvc)
			h.SetLogger(testhelpers.NewTestLogger())
			// Construye la URL con o sin query param según el caso
			url := "/api/v1/employees/reportInboundOrders"
			if tc.queryID != nil {
				url += "?id=" + strconv.Itoa(*tc.queryID)
			}
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			// Ejecuta el handler directamente
			h.Report(w, req)
			res := w.Result()
			defer res.Body.Close()
			// Validamos status y parte de la respuesta deseada
			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantContent)
		})
	}
}
