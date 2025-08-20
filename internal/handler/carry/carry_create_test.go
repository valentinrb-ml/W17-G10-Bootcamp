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
			handler.SetLogger(testhelpers.NewTestLogger())

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
