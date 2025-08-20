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

func TestWarehouseHandler_FindById(t *testing.T) {
	type arrange struct {
		mockService func() *mocks.WarehouseServiceMock
		urlID       string
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
			name: "success - warehouse found",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}
					expectedWarehouse := testhelpers.CreateExpectedWarehouse(1)

					mock.FuncFindById = func(ctx context.Context, id int) (*warehouseModel.Warehouse, error) {
						return expectedWarehouse, nil
					}
					return mock
				},
				urlID: "1",
			},
			output: output{
				statusCode: http.StatusOK,
				responseBody: func() interface{} {
					return testhelpers.CreateExpectedWarehouse(1)
				},
				err: nil,
			},
		},
		{
			name: "error - warehouse not found",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncFindById = func(ctx context.Context, id int) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
					}
					return mock
				},
				urlID: "99",
			},
			output: output{
				statusCode: http.StatusNotFound,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Not Found",
						"message": "warehouse not found",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - invalid ID format",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					return &mocks.WarehouseServiceMock{}
				},
				urlID: "invalid",
			},
			output: output{
				statusCode: http.StatusBadRequest,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Bad Request",
						"message": "Invalid id",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - internal server error",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncFindById = func(ctx context.Context, id int) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeInternal, "database connection failed")
					}
					return mock
				},
				urlID: "1",
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
			handler.SetLogger(testhelpers.NewTestLogger())

			// Configure router with ID parameter
			router := chi.NewRouter()
			router.Get("/warehouses/{id}", handler.FindById)

			// Create request with ID in URL
			url := "/warehouses/" + tc.arrange.urlID
			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			// act
			router.ServeHTTP(recorder, req)

			// assert - verify status code
			require.Equal(t, tc.output.statusCode, recorder.Code)

			// assert - verify that service was called correctly (only in successful cases)
			if tc.output.statusCode == http.StatusOK {
				require.Equal(t, 1, mockService.FindByIdCallCount)
				require.Len(t, mockService.FindByIdCalls, 1)

				// Verify that it was called with the correct ID
				actualCall := mockService.FindByIdCalls[0]
				require.Equal(t, 1, actualCall.Id)
			}

			// assert - verify JSON response
			var actualResponse interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
			require.NoError(t, err)
		})
	}
}
