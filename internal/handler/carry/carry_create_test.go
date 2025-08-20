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

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/carry"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	carryModel "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestCarryHandler_Create(t *testing.T) {
	type arrange struct {
		mockService func() *mocks.CarryServiceMock
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
			name: "success - carry created successfully",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					mock := &mocks.CarryServiceMock{}
					expectedCarry := testhelpers.CreateExpectedCarry(1)

					mock.FuncCreate = func(ctx context.Context, c carryModel.Carry) (*carryModel.Carry, error) {
						return expectedCarry, nil
					}
					return mock
				},
				requestBody: func() *bytes.Buffer {
					carryReq := testhelpers.CreateTestCarryRequest()
					jsonData, _ := json.Marshal(carryReq)
					return bytes.NewBuffer(jsonData)
				},
			},
			output: output{
				statusCode: http.StatusCreated,
				responseBody: func() interface{} {
					return testhelpers.CreateExpectedCarry(1)
				},
				err: nil,
			},
		},
		{
			name: "error - invalid JSON",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					return &mocks.CarryServiceMock{}
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
			name: "error - validation error",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					return &mocks.CarryServiceMock{}
				},
				requestBody: func() *bytes.Buffer {
					carryReq := testhelpers.CreateTestCarryRequest()
					carryReq.CompanyName = "" // Invalid field
					jsonData, _ := json.Marshal(carryReq)
					return bytes.NewBuffer(jsonData)
				},
			},
			output: output{
				statusCode: http.StatusUnprocessableEntity,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "VALIDATION_ERROR",
						"message": "company_name is required",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - service error",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					mock := &mocks.CarryServiceMock{}

					mock.FuncCreate = func(ctx context.Context, c carryModel.Carry) (*carryModel.Carry, error) {
						return nil, apperrors.NewAppError(apperrors.CodeInternal, "database connection failed")
					}
					return mock
				},
				requestBody: func() *bytes.Buffer {
					carryReq := testhelpers.CreateTestCarryRequest()
					jsonData, _ := json.Marshal(carryReq)
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
			handler := handler.NewCarryHandler(mockService)

			// Configure router
			router := chi.NewRouter()
			router.Post("/carries", handler.Create)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/carries", tc.arrange.requestBody())
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

				// Verify that it was called with the correct carry
				actualCall := mockService.CreateCalls[0]
				require.Equal(t, "CAR001", actualCall.Carry.Cid)
				require.Equal(t, "Test Company", actualCall.Carry.CompanyName)
				require.Equal(t, "Test Address", actualCall.Carry.Address)
				require.Equal(t, "5551234567", actualCall.Carry.Telephone)
				require.Equal(t, "1", actualCall.Carry.LocalityId)
			}

			// assert - verify JSON response
			var actualResponse interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
			require.NoError(t, err)
		})
	}
}

func TestCarryHandler_Create_Success_WithLogger(t *testing.T) {
	// arrange - success case with logger
	mockService := &mocks.CarryServiceMock{}
	expectedCarry := testhelpers.CreateExpectedCarry(1)

	mockService.FuncCreate = func(ctx context.Context, c carryModel.Carry) (*carryModel.Carry, error) {
		return expectedCarry, nil
	}

	h := handler.NewCarryHandler(mockService)
	h.SetLogger(&SimpleTestLogger{})

	validRequest := testhelpers.CreateTestCarryRequest()
	requestBody, _ := json.Marshal(validRequest)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/api/v1/carries", h.Create)

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCarryHandler_Create_InvalidJSON_WithLogger(t *testing.T) {
	// arrange - invalid JSON with logger
	mockService := &mocks.CarryServiceMock{}

	h := handler.NewCarryHandler(mockService)
	h.SetLogger(&SimpleTestLogger{})

	invalidJSON := `{"invalid": json`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/api/v1/carries", h.Create)

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestCarryHandler_Create_ServiceError_WithLogger(t *testing.T) {
	// arrange - service error with logger
	mockService := &mocks.CarryServiceMock{}

	mockService.FuncCreate = func(ctx context.Context, c carryModel.Carry) (*carryModel.Carry, error) {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "cid already exists")
	}

	h := handler.NewCarryHandler(mockService)
	h.SetLogger(&SimpleTestLogger{})

	validRequest := testhelpers.CreateTestCarryRequest()
	requestBody, _ := json.Marshal(validRequest)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/api/v1/carries", h.Create)

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCarryHandler_Create_ValidationError_WithLogger(t *testing.T) {
	// arrange - validation error with logger
	mockService := &mocks.CarryServiceMock{}
	carryReq := testhelpers.CreateTestCarryRequest()
	carryReq.CompanyName = "" // Invalid field to trigger validation error
	jsonData, _ := json.Marshal(carryReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	h := handler.NewCarryHandler(mockService)
	h.SetLogger(&SimpleTestLogger{})

	router := chi.NewRouter()
	router.Post("/api/v1/carries", h.Create)

	// act
	router.ServeHTTP(recorder, req)

	// assert
	require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)

	// Just verify that we get an error response indicating validation failure
	require.Contains(t, recorder.Body.String(), "VALIDATION_ERROR")
}
