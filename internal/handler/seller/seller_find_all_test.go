package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSellerHandler_FindAll(t *testing.T) {
	tests := []struct {
		name            string
		mockService     func() *mocks.SellerServiceMock
		wantStatus      int
		wantResponse    []models.ResponseSeller
		wantErrorCode   string
		wantErrorMsgSub string
	}{
		{
			name: "success - many sellers",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseSeller, error) {
					return testhelpers.FindAllSellersResponseDummy(), nil
				}
				return mock
			},
			wantStatus:   http.StatusOK,
			wantResponse: testhelpers.FindAllSellersResponseDummy(),
		},
		{
			name: "success - empty list",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseSeller, error) {
					return []models.ResponseSeller{}, nil
				}
				return mock
			},
			wantStatus:   http.StatusOK,
			wantResponse: []models.ResponseSeller{},
		},
		{
			name: "error - db/internal error",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseSeller, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "internal server error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "internal server error",
		},
		{
			name: "error - unknown error type (fallback)",
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.FindAllFn = func(ctx context.Context) ([]models.ResponseSeller, error) {
					return nil, errors.New("some unknown error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers", nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			h := handler.NewSellerHandler(tt.mockService())

			h.FindAll(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var envelope struct {
					Data []models.ResponseSeller `json:"data"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &envelope)
				require.NoError(t, err)
				require.Equal(t, tt.wantResponse, envelope.Data)
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
				require.Equal(t, tt.wantErrorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrorMsgSub)
			}
		})
	}
}
