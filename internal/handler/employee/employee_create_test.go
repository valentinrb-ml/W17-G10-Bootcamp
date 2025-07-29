package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Test del handler Create para Employee usando helpers de employee genérico.
func TestEmployeeHandler_Create(t *testing.T) {
	testCases := []struct {
		name         string
		payload      map[string]interface{}
		mockCreateFn func(ctx context.Context, e *models.Employee) (*models.Employee, error)
		wantStatus   int
		wantContent  string
	}{
		{
			name: "create_ok",
			// Armamos el payload (json esperado) usando un empleado dummy del helper
			payload: func() map[string]interface{} {
				emp := testhelpers.CreateTestEmployee()
				return map[string]interface{}{
					"card_number_id": emp.CardNumberID,
					"first_name":     emp.FirstName,
					"last_name":      emp.LastName,
					"warehouse_id":   emp.WarehouseID,
				}
			}(),
			// En el mock, devolvemos un empleado esperado (con ID) usando helper
			mockCreateFn: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
				emp := testhelpers.CreateExpectedEmployee(123)
				emp.CardNumberID = e.CardNumberID
				emp.FirstName = e.FirstName
				emp.LastName = e.LastName
				emp.WarehouseID = e.WarehouseID
				return emp, nil
			},
			wantStatus:  http.StatusCreated,
			wantContent: `"card_number_id":"EMP001"`,
		},
		{
			name: "create_fail",
			// Faltan campos (last_name y warehouse_id) a propósito,
			// pero otros se toman del empleado helper (consistencia de test)
			payload: func() map[string]interface{} {
				emp := testhelpers.CreateTestEmployee()
				return map[string]interface{}{
					"card_number_id": emp.CardNumberID,
					"first_name":     emp.FirstName,
					// Falta last_name y warehouse_id
				}
			}(),
			mockCreateFn: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
				// El service devuelve error de validación (omitimos el empleado aquí)
				return nil, apperrors.NewAppError(apperrors.CodeValidationError, "last_name cannot be empty")
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantContent: `"last_name cannot be empty"`,
		},
		{
			name: "create_conflict",
			// Usamos helper para la mayoría de los campos, sólo sobrescribimos card_number_id
			payload: func() map[string]interface{} {
				emp := testhelpers.CreateTestEmployee()
				return map[string]interface{}{
					"card_number_id": "E003",
					"first_name":     emp.FirstName,
					"last_name":      emp.LastName,
					"warehouse_id":   emp.WarehouseID,
				}
			}(),
			mockCreateFn: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "card_number_id already exists")
			},
			wantStatus:  http.StatusConflict,
			wantContent: `"card_number_id already exists"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock del service usando el método Create simulado por test case
			mockSvc := &employeeMocks.EmployeeServiceMock{MockCreate: tc.mockCreateFn}
			h := handler.NewEmployeeHandler(mockSvc)

			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("POST", "/employees", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Ejecutar el handler
			h.Create(w, req)
			res := w.Result()
			defer res.Body.Close()

			// Chequear código de respuesta y contenido esperado
			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantContent)
		})
	}
}
