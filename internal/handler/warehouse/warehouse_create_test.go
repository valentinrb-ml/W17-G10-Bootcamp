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

func TestWarehouseHandler_Create(t *testing.T) {
	type arrange struct {
		mockService func() *mocks.WarehouseServiceMock
		requestBody func() *bytes.Buffer
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
			name: "success - warehouse created successfully",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}
					expectedWarehouse := testhelpers.CreateExpectedWarehouse(1)

					mock.FuncCreate = func(ctx context.Context, w warehouseModel.Warehouse) (*warehouseModel.Warehouse, error) {
						return expectedWarehouse, nil
					}
					return mock
				},
				requestBody: func() *bytes.Buffer {
					warehouseReq := testhelpers.CreateTestWarehouseRequest()
					jsonData, _ := json.Marshal(warehouseReq)
					return bytes.NewBuffer(jsonData)
				},
			},
			output: output{
				statusCode: http.StatusCreated,
				responseBody: func() interface{} {
					return testhelpers.CreateExpectedWarehouse(1)
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
				requestBody: func() *bytes.Buffer {
					return bytes.NewBuffer([]byte(`{"invalid": json}`))
				},
			},
			output: output{
				statusCode: http.StatusInternalServerError,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Internal Server Error",
						"message": "Invalid JSON format",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - create failed",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncCreate = func(ctx context.Context, w warehouseModel.Warehouse) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeValidationError, "invalid request body")
					}
					return mock
				},
				requestBody: func() *bytes.Buffer {
					return bytes.NewBuffer([]byte(`{"address": "Av. Maipú 1234, CABA, Buenos Aires"}`))
				},
			},
			output: output{
				statusCode: http.StatusUnprocessableEntity,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "VALIDATION_ERROR",
						"message": "invalid request body",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - warehouse_code already exists",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncCreate = func(ctx context.Context, w warehouseModel.Warehouse) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists")
					}
					return mock
				},
				requestBody: func() *bytes.Buffer {
					warehouseReq := testhelpers.CreateTestWarehouseRequest()
					jsonData, _ := json.Marshal(warehouseReq)
					return bytes.NewBuffer(jsonData)
				},
			},
			output: output{
				statusCode: http.StatusConflict,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Conflict",
						"message": "warehouse_code already exists",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - validation error",
			arrange: arrange{
				mockService: func() *mocks.WarehouseServiceMock {
					mock := &mocks.WarehouseServiceMock{}

					mock.FuncCreate = func(ctx context.Context, w warehouseModel.Warehouse) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeValidationError, "warehouse_code is required")
					}
					return mock
				},
				requestBody: func() *bytes.Buffer {
					warehouseReq := testhelpers.CreateTestWarehouseRequest()
					warehouseReq.WarehouseCode = "" // Invalid field
					jsonData, _ := json.Marshal(warehouseReq)
					return bytes.NewBuffer(jsonData)
				},
			},
			output: output{
				statusCode: 422,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Validation Error",
						"message": "invalid request body",
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

					mock.FuncCreate = func(ctx context.Context, w warehouseModel.Warehouse) (*warehouseModel.Warehouse, error) {
						return nil, apperrors.NewAppError(apperrors.CodeInternal, "database connection failed")
					}
					return mock
				},
				requestBody: func() *bytes.Buffer {
					warehouseReq := testhelpers.CreateTestWarehouseRequest()
					jsonData, _ := json.Marshal(warehouseReq)
					return bytes.NewBuffer(jsonData)
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
			router.Post("/warehouses", handler.Create)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/warehouses", tc.arrange.requestBody())
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// act
			router.ServeHTTP(recorder, req)

			// assert - verify status code
			require.Equal(t, tc.output.statusCode, recorder.Code)

			// assert - verify that service was called correctly (only in successful cases)
			if tc.output.statusCode == http.StatusCreated {
				require.Equal(t, 1, mockService.CreateCallCount)
				require.Len(t, mockService.CreateCalls, 1)

				// Verify that it was called with the correct warehouse
				actualCall := mockService.CreateCalls[0]
				require.Equal(t, "WH001", actualCall.Warehouse.WarehouseCode)
				require.Equal(t, "123 Main St", actualCall.Warehouse.Address)
				require.Equal(t, 10.5, actualCall.Warehouse.MinimumTemperature)
				require.Equal(t, 1000, actualCall.Warehouse.MinimumCapacity)
				require.Equal(t, "5551234567", actualCall.Warehouse.Telephone)
				require.Equal(t, "LOC001", actualCall.Warehouse.LocalityId)
			}

			// assert - verify JSON response
			var actualResponse interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
			require.NoError(t, err)
		})
	}
}
