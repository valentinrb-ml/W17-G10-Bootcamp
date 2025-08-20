package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerHandler_FindAll(t *testing.T) {
	tests := []struct {
		name         string
		mockService  func() *mocks.BuyerServiceMock
		wantStatus   int
		response     []models.ResponseBuyer
		errorCode    string
		errorMessage string
	}{
		{
			name: "success - many buyers",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseBuyer, error) {
					return testhelpers.FindAllBuyersResponseDummy(), nil
				}
				return mock
			},
			wantStatus: http.StatusOK,
			response:   testhelpers.FindAllBuyersResponseDummy(),
		},
		{
			name: "success - empty list",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseBuyer, error) {
					return []models.ResponseBuyer{}, nil
				}
				return mock
			},
			wantStatus: http.StatusOK,
			response:   []models.ResponseBuyer{},
		},
		{
			name: "error - db/internal error",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseBuyer, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "internal server error")
				}
				return mock
			},
			wantStatus:   http.StatusInternalServerError,
			errorCode:    apperrors.CodeInternal,
			errorMessage: "internal server error",
		},
		{
			name: "error - unknown error type (fallback)",
			mockService: func() *mocks.BuyerServiceMock {
				mock := &mocks.BuyerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseBuyer, error) {
					return nil, errors.New("some unknown error")
				}
				return mock
			},
			wantStatus:   http.StatusInternalServerError,
			errorCode:    apperrors.CodeInternal,
			errorMessage: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/buyers", nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			h := handler.NewBuyerHandler(tt.mockService())
			h.SetLogger(testhelpers.NewTestLogger())

			h.FindAll(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var envelope struct {
					Data []models.ResponseBuyer `json:"data"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &envelope)
				require.NoError(t, err)
				require.Equal(t, tt.response, envelope.Data)
			} else {
				var body struct {
					Error struct {
						Code    string         `json:"code"`
						Message string         `json:"message"`
						Details map[string]any `json:"details"`
					} `json:"error"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.errorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.errorMessage)
			}
		})
	}
}
