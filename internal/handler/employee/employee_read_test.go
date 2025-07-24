package employee

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/require"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Test for GET /employees (find_all)
func TestEmployeeHandler_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockFindAll func(ctx context.Context) ([]*models.Employee, error)
		wantStatus  int
		wantBodyHas string
	}{
		{
			name: "find_all",
			mockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
				return []*models.Employee{
					{ID: 1, CardNumberID: "E001", FirstName: "Lucas", LastName: "Martinez", WarehouseID: 1},
					{ID: 2, CardNumberID: "E002", FirstName: "Paola", LastName: "Gomez", WarehouseID: 1},
				}, nil
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"card_number_id":"E001"`,
		},
		{
			name: "find_all_empty",
			mockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
				return []*models.Employee{}, nil
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `[]`, // respuesta vacía
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &employeeMocks.EmployeeServiceMock{MockFindAll: tc.mockFindAll}
			h := NewEmployeeHandler(mockSvc)

			req := httptest.NewRequest("GET", "/employees", nil)
			w := httptest.NewRecorder()

			h.GetAll(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantBodyHas)
		})
	}
}

// Test for GET /employees/{id}
func TestEmployeeHandler_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		mockFindByID func(ctx context.Context, id int) (*models.Employee, error)
		wantStatus   int
		wantBodyHas  string
	}{
		{
			name: "find_by_id_existent",
			id:   1,
			mockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
				return &models.Employee{ID: 1, CardNumberID: "E001", FirstName: "Lucas", LastName: "Martinez", WarehouseID: 1}, nil
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"card_number_id":"E001"`,
		},
		{
			name: "find_by_id_non_existent",
			id:   99,
			mockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
			},
			wantStatus:  http.StatusNotFound,
			wantBodyHas: `"employee not found"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &employeeMocks.EmployeeServiceMock{MockFindByID: tc.mockFindByID}
			h := NewEmployeeHandler(mockSvc)

			req := httptest.NewRequest("GET", "/employees/"+strconv.Itoa(tc.id), nil)
			w := httptest.NewRecorder()

			// Para emular el router de chi, seteamos chi context param
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			h.GetByID(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantBodyHas)
		})
	}
}
