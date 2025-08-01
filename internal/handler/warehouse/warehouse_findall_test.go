package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/warehouse"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	warehouseModel "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestWarehouseHandler_FindAll(t *testing.T) {
	type arrange struct {
		mockService func() *mocks.WarehouseServiceMock
	}
	type output struct {
		statusCode   int
		responseBody func() interface{}
		err          error
	}
	type testCase struct {
		name    string
		arrange arrange
		output  output
	}

	// test cases
	testCases := []testCase{
		{
			name: "success - warehouses found",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}
					expectedWarehouses := []warehouseModel.Warehouse{
						*testhelpers.CreateExpectedWarehouse(1),
						*testhelpers.CreateExpectedWarehouse(2),
					}

					mock.FuncFindAll = func(ctx context.Context) ([]warehouseModel.Warehouse, error) {
						return expectedWarehouses, nil
					}
					return mock
				},
			},
			output: output{
				statusCode: http.StatusOK,
				responseBody: func() interface{} {
					return []warehouseModel.Warehouse{
						*testhelpers.CreateExpectedWarehouse(1),
						*testhelpers.CreateExpectedWarehouse(2),
					}
				},
				err: nil,
			},
		},
		{
			name: "success - empty list",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncFindAll = func(ctx context.Context) ([]warehouseModel.Warehouse, error) {
						return []warehouseModel.Warehouse{}, nil
					}
					return mock
				},
			},
			output: output{
				statusCode: http.StatusOK,
				responseBody: func() interface{} {
					return []interface{}{}
				},
				err: nil,
			},
		},
		{
			name: "error - internal server error",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncFindAll = func(ctx context.Context) ([]warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeInternal, "database connection failed")
					}
					return mock
				},
			},
			output: output{
				statusCode: http.StatusInternalServerError,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Internal Server Error",
						"message": "database connection failed",
					}
				},
				err: nil,
			},
		},
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			mockService := tc.arrange.mockService()
			handler := handler.NewWarehouseHandler(mockService)

			// Configure router
			router := chi.NewRouter()
			router.Get("/warehouses", handler.FindAll)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
			recorder := httptest.NewRecorder()

			// act
			router.ServeHTTP(recorder, req)

			// assert - verify status code
			require.Equal(t, tc.output.statusCode, recorder.Code)

			// assert - verify that service was called correctly (only in successful cases)
			if tc.output.statusCode == http.StatusOK {
				require.Equal(t, 1, mockService.FindAllCallCount)
			}

			// assert - verify JSON response
			var actualResponse interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
			require.NoError(t, err)
		})
	}
}
