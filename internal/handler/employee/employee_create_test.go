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
)

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
			payload: map[string]interface{}{
				"card_number_id": "E001",
				"first_name":     "Paola",
				"last_name":      "Lopez",
				"warehouse_id":   1,
			},
			mockCreateFn: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
				return &models.Employee{
					ID:           123,
					CardNumberID: e.CardNumberID,
					FirstName:    e.FirstName,
					LastName:     e.LastName,
					WarehouseID:  e.WarehouseID,
				}, nil
			},
			wantStatus:  http.StatusCreated,
			wantContent: `"card_number_id":"E001"`,
		},
		{
			name: "create_fail",
			payload: map[string]interface{}{
				"card_number_id": "E002",
				"first_name":     "FailSinLastname",
				// last_name y warehouse_id faltantes
			},
			mockCreateFn: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
				return nil, apperrors.NewAppError(apperrors.CodeValidationError, "last_name cannot be empty")
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantContent: `"last_name cannot be empty"`,
		},
		{
			name: "create_conflict",
			payload: map[string]interface{}{
				"card_number_id": "E003",
				"first_name":     "Lucas",
				"last_name":      "Martinez",
				"warehouse_id":   1,
			},
			mockCreateFn: func(ctx context.Context, e *models.Employee) (*models.Employee, error) {
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "card_number_id already exists")
			},
			wantStatus:  http.StatusConflict,
			wantContent: `"card_number_id already exists"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &employeeMocks.EmployeeServiceMock{MockCreate: tc.mockCreateFn}
			h := handler.NewEmployeeHandler(mockSvc)

			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("POST", "/employees", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Create(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantContent)
		})
	}
}
