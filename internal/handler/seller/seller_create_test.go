package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func dummyRequestSeller() models.RequestSeller {
	cid := 101
	companyName := "Frutas del Sur"
	address := "Calle 1"
	telephone := "221-111"
	localityId := "1900"

	return models.RequestSeller{
		Cid:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		LocalityId:  &localityId,
	}
}

func dummyResponseSeller() models.ResponseSeller {
	return models.ResponseSeller{
		Id:          1,
		Cid:         101,
		CompanyName: "Frutas del Sur",
		Address:     "Calle 1",
		Telephone:   "221-111",
		LocalityId:  "1900",
	}
}

func TestSellerHandler_Create(t *testing.T) {
	type args struct {
		requestBody any
	}

	tests := []struct {
		name             string
		args             args
		mockService      func() *mocks.SellerServiceMock
		wantStatus       int
		wantResponseBody any
		wantErrorCode    string
		wantErrorMsgSub  string
	}{
		{
			name: "success",
			args: args{
				requestBody: dummyRequestSeller(),
			},
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.CreateFn = func(ctx context.Context, req models.RequestSeller) (*models.ResponseSeller, error) {
					expected := dummyResponseSeller()
					return &expected, nil
				}
				return mock
			},
			wantStatus:       http.StatusCreated,
			wantResponseBody: dummyResponseSeller(),
		},
		{
			name: "error - invalid request payload",
			args: args{
				requestBody: `{invalid json}`,
			},
			mockService: func() *mocks.SellerServiceMock {
				return &mocks.SellerServiceMock{}
			},
			wantStatus:      http.StatusBadRequest,
			wantErrorCode:   apperrors.CodeBadRequest,
			wantErrorMsgSub: "invalid JSON format",
		},
		{
			name: "error - invalid Seller validation",
			args: args{
				requestBody: models.RequestSeller{},
			},
			mockService: func() *mocks.SellerServiceMock {
				return &mocks.SellerServiceMock{}
			},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorCode:   apperrors.CodeValidationError,
			wantErrorMsgSub: "required",
		},
		{
			name: "error - service layer returns error",
			args: args{
				requestBody: dummyRequestSeller(),
			},
			mockService: func() *mocks.SellerServiceMock {
				mock := &mocks.SellerServiceMock{}
				mock.CreateFn = func(ctx context.Context, req models.RequestSeller) (*models.ResponseSeller, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "unexpected error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "unexpected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBodyBytes []byte
			switch val := tt.args.requestBody.(type) {
			case string:
				requestBodyBytes = []byte(val)
			default:
				b, err := json.Marshal(tt.args.requestBody)
				require.NoError(t, err)
				requestBodyBytes = b
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewReader(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h := handler.NewSellerHandler(tt.mockService())

			// Call handler
			h.Create(rec, req)

			// Check status
			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusCreated {
				var responseEnvelope struct {
					Data models.ResponseSeller `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &responseEnvelope)
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseBody, responseEnvelope.Data)
			} else {
				var body struct {
					Error struct {
						Code    string         `json:"code"`
						Message string         `json:"message"`
						Details map[string]any `json:"details"`
					} `json:"error"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.wantErrorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrorMsgSub)
			}
		})
	}
}
