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
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Test para GET /employees (find_all)
func TestEmployeeHandler_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockFindAll func(ctx context.Context) ([]*models.Employee, error)
		wantStatus  int
		wantBodyHas string
	}{
		{
			name: "find_all",
			// Usamos el helper para poblar la lista de empleados dummy
			mockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
				emps := testhelpers.CreateTestEmployees()
				var empsPtrs []*models.Employee
				for i := range emps {
					empsPtrs = append(empsPtrs, &emps[i])
				}
				return empsPtrs, nil
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"card_number_id":"EMP001"`, // checa que JWT de testhelpers esté presente
		},
		{
			name: "find_all_empty",
			mockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
				return []*models.Employee{}, nil // caso vacío
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `[]`, // respuesta vacía
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Crea el mock del service usando la función del test
			mockSvc := &employeeMocks.EmployeeServiceMock{MockFindAll: tc.mockFindAll}
			h := handler.NewEmployeeHandler(mockSvc)

			// Ejecuta el GET como lo haría el router
			req := httptest.NewRequest("GET", "/employees", nil)
			w := httptest.NewRecorder()

			h.GetAll(w, req)
			res := w.Result()
			defer res.Body.Close()

			// Checa código HTTP y que en el body está el dato esperado
			require.Equal(t, tc.wantStatus, res.StatusCode)
			require.Contains(t, w.Body.String(), tc.wantBodyHas)
		})
	}
}

// Test para GET /employees/{id}
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
			id:   func() int { return testhelpers.CreateTestEmployee().ID }(), // Consistencia
			mockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
				return testhelpers.CreateExpectedEmployee(id), nil // Devuelve un empleado dummy usando el helper
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"card_number_id":"EMP001"`, // o cambia según helper
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
			h := handler.NewEmployeeHandler(mockSvc)

			req := httptest.NewRequest("GET", "/employees/"+strconv.Itoa(tc.id), nil)
			w := httptest.NewRecorder()

			// Emula el router de chi (param id)
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
