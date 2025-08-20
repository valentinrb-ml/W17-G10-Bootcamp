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
)

func TestWarehouseHandler_Delete(t *testing.T) {
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
			name: "success - warehouse deleted",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncDelete = func(ctx context.Context, id int) error {
						return nil
					}
					return mock
				},
				urlID: "1",
			},
			output: output{
				statusCode: http.StatusNoContent,
				responseBody: func() interface{} {
					return nil // No content expected for successful delete
				},
				err: nil,
			},
		},
		{
			name: "error - warehouse not found",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncDelete = func(ctx context.Context, id int) error {
						return apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
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
			name: "error - foreign key constraint violation",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncDelete = func(ctx context.Context, id int) error {
						return apperrors.NewAppError(apperrors.CodeConflict, "warehouse is referenced by other entities")
					}
					return mock
				},
				urlID: "1",
			},
			output: output{
				statusCode: http.StatusConflict,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Conflict",
						"message": "warehouse is referenced by other entities",
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

					mock.FuncDelete = func(ctx context.Context, id int) error {
						return apperrors.NewAppError(apperrors.CodeInternal, "database connection failed")
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

			// Configure router with ID parameter
			router := chi.NewRouter()
			router.Delete("/warehouses/{id}", handler.Delete)

			// Create request with ID in URL
			url := "/warehouses/" + tc.arrange.urlID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			recorder := httptest.NewRecorder()

			// act
			router.ServeHTTP(recorder, req)

			// assert - verify status code
			require.Equal(t, tc.output.statusCode, recorder.Code)

			// assert - verify that service was called correctly (only in successful cases)
			if tc.output.statusCode == http.StatusNoContent {
				require.Equal(t, 1, mockService.DeleteCallCount)
				require.Len(t, mockService.DeleteCalls, 1)

				// Verify that it was called with the correct ID
				actualCall := mockService.DeleteCalls[0]
				require.Equal(t, 1, actualCall.Id)
			}

			// assert - verify response body (only if content is expected)
			if tc.output.statusCode != http.StatusNoContent {
				var actualResponse interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
				require.NoError(t, err)
			} else {
				// For 204 No Content, body should be empty
				require.Empty(t, recorder.Body.String())
			}
		})
	}
}

func TestWarehouseHandler_Delete_WithNilLogger(t *testing.T) {
	// arrange
	mock := &mocks.WarehouseServiceMock{}

	mock.FuncDelete = func(ctx context.Context, id int) error {
		return nil
	}

	handler := handler.NewWarehouseHandler(mock)
	// Don't set logger, so it remains nil

	// Configure router
	router := chi.NewRouter()
	router.Delete("/warehouses/{id}", handler.Delete)

	// Create request
	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
	recorder := httptest.NewRecorder()

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Empty(t, recorder.Body.String())
}

func TestWarehouseHandler_Delete_Success_WithLogger(t *testing.T) {
	// arrange - success case with logger
	mockService := &mocks.WarehouseServiceMock{}

	mockService.FuncDelete = func(ctx context.Context, id int) error {
		return nil
	}

	warehouseHandler := handler.NewWarehouseHandler(mockService)
	warehouseHandler.SetLogger(&SimpleTestLogger{})

	router := chi.NewRouter()
	router.Delete("/warehouses/{id}", warehouseHandler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
	recorder := httptest.NewRecorder()

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestWarehouseHandler_Delete_InvalidID_WithLogger(t *testing.T) {
	// arrange - invalid ID with logger
	mockService := &mocks.WarehouseServiceMock{}

	warehouseHandler := handler.NewWarehouseHandler(mockService)
	warehouseHandler.SetLogger(&SimpleTestLogger{})

	router := chi.NewRouter()
	router.Delete("/warehouses/{id}", warehouseHandler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/invalid", nil)
	recorder := httptest.NewRecorder()

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestWarehouseHandler_Delete_ServiceError_WithLogger(t *testing.T) {
	// arrange - service error with logger
	mockService := &mocks.WarehouseServiceMock{}

	mockService.FuncDelete = func(ctx context.Context, id int) error {
		return apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
	}

	warehouseHandler := handler.NewWarehouseHandler(mockService)
	warehouseHandler.SetLogger(&SimpleTestLogger{})

	router := chi.NewRouter()
	router.Delete("/warehouses/{id}", warehouseHandler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
	recorder := httptest.NewRecorder()

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusNotFound, recorder.Code)
}
