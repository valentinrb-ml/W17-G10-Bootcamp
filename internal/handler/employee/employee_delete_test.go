package employee

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
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
			id:   1,
			mockDeleteFn: func(ctx context.Context, id int) error {
				return nil // Eliminación exitosa
			},
			wantStatus:  http.StatusNoContent,
			wantContent: "",
		},
		{
			name: "delete_non_existent",
			id:   99,
			mockDeleteFn: func(ctx context.Context, id int) error {
				return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
			},
			wantStatus:  http.StatusNotFound,
			wantContent: `"employee not found"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mocks.EmployeeServiceMock{
				MockDelete: tc.mockDeleteFn,
				// Los demás métodos pueden quedar nil.
			}
			h := NewEmployeeHandler(mockSvc)

			req := httptest.NewRequest("DELETE", "/employees/"+strconv.Itoa(tc.id), nil)
			w := httptest.NewRecorder()

			// Simula seteo de chi param (como haría el router)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			h.Delete(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tc.wantStatus, res.StatusCode)
			if tc.wantContent != "" {
				require.Contains(t, w.Body.String(), tc.wantContent)
			}
		})
	}
}
