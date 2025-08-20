package handler_test

import (
	"bytes"
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

func TestWarehouseHandler_Update(t *testing.T) {
	type arrange struct {
		mockService func() *mocks.WarehouseServiceMock
		urlID       string
		requestBody interface{}
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
			name: "success - warehouse updated",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}
					expectedWarehouse := testhelpers.CreateExpectedWarehouse(1)

					mock.FuncUpdate = func(ctx context.Context, id int, patch warehouseModel.WarehousePatchDTO) (*warehouseModel.Warehouse, error) {
						return expectedWarehouse, nil
					}
					return mock
				},
				urlID: "1",
				requestBody: warehouseModel.WarehousePatchDTO{
					WarehouseCode: testhelpers.StringPtr("UPD-001"),
					Address:       testhelpers.StringPtr("Updated Address 123"),
					Telephone:     testhelpers.StringPtr("9876543210"),
				},
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

					mock.FuncUpdate = func(ctx context.Context, id int, patch warehouseModel.WarehousePatchDTO) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
					}
					return mock
				},
				urlID: "99",
				requestBody: warehouseModel.WarehousePatchDTO{
					WarehouseCode: testhelpers.StringPtr("UPD-001"),
				},
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
				requestBody: warehouseModel.WarehousePatchDTO{
					WarehouseCode: testhelpers.StringPtr("UPD-001"),
				},
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
			name: "error - invalid JSON",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					return &mocks.WarehouseServiceMock{}
				},
				urlID:       "1",
				requestBody: `{"invalid": json}`,
			},
			output: output{
				statusCode: http.StatusInternalServerError,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Internal Server Error",
						"message": "invalid character 'j' looking for beginning of value",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - validation failed",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncUpdate = func(ctx context.Context, id int, patch warehouseModel.WarehousePatchDTO) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeValidationError, "invalid phone format")
					}
					return mock
				},
				urlID: "1",
				requestBody: warehouseModel.WarehousePatchDTO{
					Telephone: testhelpers.StringPtr("invalid-phone"),
				},
			},
			output: output{
				statusCode: http.StatusUnprocessableEntity,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Unprocessable Entity",
						"message": "invalid phone format",
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

					mock.FuncUpdate = func(ctx context.Context, id int, patch warehouseModel.WarehousePatchDTO) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeInternal, "database connection failed")
					}
					return mock
				},
				urlID: "1",
				requestBody: warehouseModel.WarehousePatchDTO{
					WarehouseCode: testhelpers.StringPtr("UPD-001"),
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
			handler.SetLogger(testhelpers.NewTestLogger())

			// Configure router with ID parameter
			router := chi.NewRouter()
			router.Patch("/warehouses/{id}", handler.Update)

			// Create request body
			var requestBody []byte
			var err error
			if str, ok := tc.arrange.requestBody.(string); ok {
				// Handle invalid JSON string case
				requestBody = []byte(str)
			} else {
				// Handle normal struct marshaling
				requestBody, err = json.Marshal(tc.arrange.requestBody)
				require.NoError(t, err)
			}

			// Create request with ID in URL and JSON body
			url := "/warehouses/" + tc.arrange.urlID
			req := httptest.NewRequest(http.MethodPatch, url, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// act
			router.ServeHTTP(recorder, req)

			// assert - verify status code
			require.Equal(t, tc.output.statusCode, recorder.Code)

			// assert - verify that service was called correctly (only in successful cases)
			if tc.output.statusCode == http.StatusOK {
				require.Equal(t, 1, mockService.UpdateCallCount)
				require.Len(t, mockService.UpdateCalls, 1)

				// Verify that it was called with the correct ID
				actualCall := mockService.UpdateCalls[0]
				require.Equal(t, 1, actualCall.Id)
			}

			// assert - verify JSON response
			var actualResponse interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
			require.NoError(t, err)
		})
	}
}
