package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/inbound_order"
	inboundOrderMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestInboundOrderHandler_Create(t *testing.T) {
	testCases := []struct {
		name         string
		payload      map[string]interface{}
		mockCreateFn func(ctx context.Context, in *models.InboundOrder) (*models.InboundOrder, error)
		wantStatus   int
		wantContent  string
	}{
		{
			name: "create_ok",
			// Armamos el payload usando el helper para DRY y centralización de datos base
			payload: func() map[string]interface{} {
				ord := testhelpers.CreateTestInboundOrder()
				return map[string]interface{}{
					"order_number":     ord.OrderNumber,
					"order_date":       ord.OrderDate,
					"employee_id":      ord.EmployeeID,
					"product_batch_id": ord.ProductBatchID,
					"warehouse_id":     ord.WarehouseID,
				}
			}(),
			// El mock simula creación exitosa y devuelve un order con ID esperando en respuesta
			mockCreateFn: func(ctx context.Context, in *models.InboundOrder) (*models.InboundOrder, error) {
				o := testhelpers.CreateExpectedInboundOrder(22)
				return o, nil
			},
			wantStatus:  http.StatusCreated,
			wantContent: `"order_number":"INV001"`,
		},
		{
			name: "create_validation_error",
			// Payload incompleto: simulando JSON inválido (muestra cómo sería válido el helper)
			payload: map[string]interface{}{
				// Faltan campos importantes
				"order_number": "",
			},
			mockCreateFn: func(ctx context.Context, in *models.InboundOrder) (*models.InboundOrder, error) {
				return nil, apperrors.NewAppError(apperrors.CodeValidationError, "invalid inbound_order")
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantContent: `"invalid inbound_order"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Inyectamos el mock del servicio en el handler
			mockSvc := &inboundOrderMocks.InboundOrderServiceMock{MockCreate: tc.mockCreateFn}
			h := handler.NewInboundOrderHandler(mockSvc)
			h.SetLogger(testhelpers.NewTestLogger())
			// Marshal el payload en el formato recibido por el handler
			body, _ := json.Marshal(map[string]interface{}{"data": tc.payload})
			req := httptest.NewRequest("POST", "/api/v1/inboundOrders", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			// Ejecutamos el handler
			h.Create(w, req)
			res := w.Result()
			defer res.Body.Close()
			// Verifica status y presencia de contenido esperado en la respuesta
			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantContent)
		})
	}
}
