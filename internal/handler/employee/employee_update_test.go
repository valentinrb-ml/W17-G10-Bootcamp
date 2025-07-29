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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
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
			id: func() int {
				// Usamos el ID del helper para mantener consistencia
				return testhelpers.CreateTestEmployee().ID // normalmente 1
			}(),
			patch: map[string]interface{}{
				"first_name": "NombreActualizado",
			},
			// Devuelve un empleado modificado usando el helper como base
			mockUpdateFn: func(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
				e := testhelpers.CreateExpectedEmployee(id)
				e.FirstName = "NombreActualizado"
				return e, nil
			},
			wantStatus:  http.StatusOK,
			wantContent: `"first_name":"NombreActualizado"`,
		},
		{
			name: "update_non_existent",
			id:   99, // Un id que no existe en los helpers
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
			// Configuramos el mock del service con el helper/func que queremos
			mockSvc := &employeeMocks.EmployeeServiceMock{
				MockUpdate: tc.mockUpdateFn,
			}
			h := handler.NewEmployeeHandler(mockSvc)

			// Genera el body del PATCH como JSON usando lo que se va a actualizar
			patchBody, _ := json.Marshal(tc.patch)
			url := "/employees/" + strconv.Itoa(tc.id)
			req := httptest.NewRequest("PATCH", url, bytes.NewReader(patchBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Simula cómo chi pasaría los parámetros de ruta (id)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Ejecuta el handler PATCH
			h.Update(w, req)
			res := w.Result()
			defer res.Body.Close()

			// Chequea respuesta HTTP y contenido esperado
			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantContent)
		})
	}
}
