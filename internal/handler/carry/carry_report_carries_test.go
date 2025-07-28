package handler_test

import (
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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestCarryHandler_ReportCarries(t *testing.T) {
	type arrange struct {
		mockService        func() *mocks.CarryServiceMock
		queryParam         string
		expectedLocalityID *string
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
			name: "success - get all carries report (no query param)",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					mock := &mocks.CarryServiceMock{}

					mock.FuncGetCarriesReport = func(ctx context.Context, localityID *string) (interface{}, error) {
						return testhelpers.CreateTestCarriesReportSlice(), nil
					}
					return mock
				},
				queryParam:         "",
				expectedLocalityID: nil,
			},
			output: output{
				statusCode: http.StatusOK,
				responseBody: func() interface{} {
					return testhelpers.CreateTestCarriesReportSlice()
				},
				err: nil,
			},
		},
		{
			name: "success - get carries report by locality (with query param)",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					mock := &mocks.CarryServiceMock{}

					mock.FuncGetCarriesReport = func(ctx context.Context, localityID *string) (interface{}, error) {
						return testhelpers.CreateTestCarriesReport("1", "Test Locality 1", 5), nil
					}
					return mock
				},
				queryParam:         "?id=1",
				expectedLocalityID: testhelpers.StringPtr("1"),
			},
			output: output{
				statusCode: http.StatusOK,
				responseBody: func() interface{} {
					return testhelpers.CreateTestCarriesReport("1", "Test Locality 1", 5)
				},
				err: nil,
			},
		},
		{
			name: "error - locality not found",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					mock := &mocks.CarryServiceMock{}

					mock.FuncGetCarriesReport = func(ctx context.Context, localityID *string) (interface{}, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "locality not found")
					}
					return mock
				},
				queryParam:         "?id=999",
				expectedLocalityID: testhelpers.StringPtr("999"),
			},
			output: output{
				statusCode: http.StatusNotFound,
				responseBody: func() interface{} {
					return map[string]interface{}{
						"error":   "Not Found",
						"message": "locality not found",
					}
				},
				err: nil,
			},
		},
		{
			name: "error - internal server error",
			arrange: arrange{
				mockService: func() *mocks.CarryServiceMock {
					mock := &mocks.CarryServiceMock{}

					mock.FuncGetCarriesReport = func(ctx context.Context, localityID *string) (interface{}, error) {
						return nil, apperrors.NewAppError(apperrors.CodeInternal, "database connection failed")
					}
					return mock
				},
				queryParam:         "",
				expectedLocalityID: nil,
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
			router.Get("/carries/reportCarries", handler.ReportCarries)

			// Create request
			url := "/carries/reportCarries" + tc.arrange.queryParam
			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			// act
			router.ServeHTTP(recorder, req)

			// assert - verify status code
			require.Equal(t, tc.output.statusCode, recorder.Code)

			// assert - verify that service was called correctly
			require.Equal(t, 1, mockService.GetCarriesReportCallCount)
			require.Len(t, mockService.GetCarriesReportCalls, 1)

			// Verify the localityID parameter passed to service
			actualCall := mockService.GetCarriesReportCalls[0]
			if tc.arrange.expectedLocalityID == nil {
				require.Nil(t, actualCall.LocalityID)
			} else {
				require.NotNil(t, actualCall.LocalityID)
				require.Equal(t, *tc.arrange.expectedLocalityID, *actualCall.LocalityID)
			}

			// assert - verify JSON response
			var actualResponse interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
			require.NoError(t, err)
		})
	}
}
