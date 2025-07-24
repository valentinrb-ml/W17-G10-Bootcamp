package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

func TestEmployeeHandler_Update(t *testing.T) {
	testCases := []struct {
		name         string
		id           int
		patch        map[string]interface{}
		mockUpdateFn func(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error)
		wantStatus   int
		wantContent  string
	}{
		{
			name: "update_ok",
			id:   1,
			patch: map[string]interface{}{
				"first_name": "NombreActualizado",
			},
			mockUpdateFn: func(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
				return &models.Employee{
					ID:           1,
					CardNumberID: "E001",
					FirstName:    "NombreActualizado",
					LastName:     "Martinez",
					WarehouseID:  1,
				}, nil
			},
			wantStatus:  http.StatusOK,
			wantContent: `"first_name":"NombreActualizado"`,
		},
		{
			name: "update_non_existent",
			id:   99,
			patch: map[string]interface{}{
				"first_name": "NuevoNombre",
			},
			mockUpdateFn: func(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
			},
			wantStatus:  http.StatusNotFound,
			wantContent: `"employee not found"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &employeeMocks.EmployeeServiceMock{
				MockUpdate: tc.mockUpdateFn,
			}
			h := handler.NewEmployeeHandler(mockSvc)

			patchBody, _ := json.Marshal(tc.patch)
			url := "/employees/" + strconv.Itoa(tc.id)
			req := httptest.NewRequest("PATCH", url, bytes.NewReader(patchBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup chi context param para el id
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			h.Update(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantContent)
		})
	}
}
