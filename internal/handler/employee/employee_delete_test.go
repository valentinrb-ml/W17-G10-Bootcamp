package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestEmployeeHandler_Delete(t *testing.T) {
	testCases := []struct {
		name         string
		id           int
		mockDeleteFn func(ctx context.Context, id int) error
		wantStatus   int
		wantContent  string
	}{
		{
			name: "delete_ok",
			// Usamos el helper para asegurar que el id existe como dummy
			id: func() int {
				emp := testhelpers.CreateTestEmployee()
				return emp.ID // normalmente 1 según CreateTestEmployee()
			}(),
			// El mock simula eliminación exitosa para ese id
			mockDeleteFn: func(ctx context.Context, id int) error {
				return nil
			},
			wantStatus:  http.StatusNoContent,
			wantContent: "", // 204 no devuelve nada
		},
		{
			name: "delete_non_existent",
			// Usamos un id que no existe en los testhelpers para simular "inexistente"
			id: 99,
			// El mock simula una eliminación fallida (NOT_FOUND)
			mockDeleteFn: func(ctx context.Context, id int) error {
				return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
			},
			wantStatus:  http.StatusNotFound,
			wantContent: `"employee not found"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Creamos un service mock que usará el handler
			mockSvc := &employeeMocks.EmployeeServiceMock{MockDelete: tc.mockDeleteFn}
			h := handler.NewEmployeeHandler(mockSvc)

			// Armamos la request con el id, y el seteo de param chi (simula routing real)
			req := httptest.NewRequest("DELETE", "/employees/"+strconv.Itoa(tc.id), nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Ejecutamos el handler (como lo haría el router)
			h.Delete(w, req)
			res := w.Result()
			defer res.Body.Close()

			// Checamos el código recibido y el contenido en caso de error
			require.Equal(t, tc.wantStatus, res.StatusCode)
			if tc.wantContent != "" {
				require.Contains(t, w.Body.String(), tc.wantContent)
			}
		})
	}
}
