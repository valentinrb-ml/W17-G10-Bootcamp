package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/geography"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestGeographyHandler_CountSellersByLocality(t *testing.T) {
	dummyResp := testhelpers.DummyResponseLocalitySellers()
	dummySliceResp := testhelpers.DummySliceResponseLocalitySellers()

	tests := []struct {
		name            string
		query           string
		mockService     func() *mocks.GeographyServiceMock
		wantStatus      int
		wantIsSlice     bool
		wantErrorCode   string
		wantErrorMsgSub string
	}{
		{
			name:  "success - by id",
			query: "?id=2000",
			mockService: func() *mocks.GeographyServiceMock {
				mock := &mocks.GeographyServiceMock{}
				mock.CountSellersByLocalityFn = func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
					require.Equal(t, "2000", id)
					return &dummyResp, nil
				}
				return mock
			},
			wantStatus:  http.StatusOK,
			wantIsSlice: false,
		},
		{
			name:  "success - grouped",
			query: "",
			mockService: func() *mocks.GeographyServiceMock {
				mock := &mocks.GeographyServiceMock{}
				mock.CountSellersGroupedByLocalityFn = func(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
					return dummySliceResp, nil
				}
				return mock
			},
			wantStatus:  http.StatusOK,
			wantIsSlice: true,
		},
		{
			name:  "error - id not found",
			query: "?id=9999",
			mockService: func() *mocks.GeographyServiceMock {
				mock := &mocks.GeographyServiceMock{}
				mock.CountSellersByLocalityFn = func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The locality you are looking for does not exist.")
				}
				return mock
			},
			wantStatus:      http.StatusNotFound,
			wantErrorCode:   apperrors.CodeNotFound,
			wantErrorMsgSub: "does not exist",
		},
		{
			name:  "error - grouped DB/internal error",
			query: "",
			mockService: func() *mocks.GeographyServiceMock {
				mock := &mocks.GeographyServiceMock{}
				mock.CountSellersGroupedByLocalityFn = func(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "grouped error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "grouped error",
		},
		{
			name:  "error - unexpected panic",
			query: "?id=1999",
			mockService: func() *mocks.GeographyServiceMock {
				mock := &mocks.GeographyServiceMock{}
				mock.CountSellersByLocalityFn = func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
					return nil, errors.New("internal server error")
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

			req := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportSellers"+tt.query, nil)
			rec := httptest.NewRecorder()
			h := handler.NewGeographyHandler(tt.mockService())
			h.SetLogger(testhelpers.NewTestLogger())

			h.CountSellersByLocality(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				if tt.wantIsSlice {
					var envelope struct {
						Data []models.ResponseLocalitySellers `json:"data"`
					}
					err := json.Unmarshal(rec.Body.Bytes(), &envelope)
					require.NoError(t, err)
					require.Equal(t, dummySliceResp, envelope.Data)
				} else {
					var envelope struct {
						Data models.ResponseLocalitySellers `json:"data"`
					}
					err := json.Unmarshal(rec.Body.Bytes(), &envelope)
					require.NoError(t, err)
					require.Equal(t, dummyResp, envelope.Data)
				}
			} else {
				var body struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
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
