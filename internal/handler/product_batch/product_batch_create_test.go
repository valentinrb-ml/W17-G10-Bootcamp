package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_batch"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_batch"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductBatchHandler_CreateProductBatches(t *testing.T) {
	type args struct {
		requestBody any
	}
	tests := []struct {
		name             string
		args             args
		mockService      func() *mocks.ProductBatchServiceMock
		wantStatus       int
		wantResponseBody any
		wantErrorCode    string
		wantErrorMsgSub  string
	}{
		{
			name: "success",
			args: args{
				requestBody: testhelpers.DummyPostProductBatch(1),
			},
			mockService: func() *mocks.ProductBatchServiceMock {
				mock := &mocks.ProductBatchServiceMock{}
				mock.FuncCreate = func(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
					dummy := testhelpers.DummyProductBatch(1)
					return &dummy, nil
				}
				return mock
			},
			wantStatus:       http.StatusCreated,
			wantResponseBody: testhelpers.DummyResponseProductBatch(1),
		},
		{
			name: "error - invalid request payload",
			args: args{
				requestBody: `{not-valid-json}`,
			},
			mockService: func() *mocks.ProductBatchServiceMock {
				return &mocks.ProductBatchServiceMock{}
			},
			wantStatus:      http.StatusBadRequest,
			wantErrorCode:   apperrors.CodeBadRequest,
			wantErrorMsgSub: "invalid JSON",
		},
		{
			name: "error - validation error",
			args: args{
				requestBody: models.PostProductBatches{},
			},
			mockService: func() *mocks.ProductBatchServiceMock {
				return &mocks.ProductBatchServiceMock{}
			},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorCode:   apperrors.CodeValidationError,
			wantErrorMsgSub: "required",
		},
		{
			name: "error - service error",
			args: args{
				requestBody: testhelpers.DummyPostProductBatch(1),
			},
			mockService: func() *mocks.ProductBatchServiceMock {
				mock := &mocks.ProductBatchServiceMock{}
				mock.FuncCreate = func(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
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

			req, err := http.NewRequest(http.MethodPost, "/api/v1/product-batches", bytes.NewReader(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h := handler.NewProductBatchesHandler(tt.mockService())

			h.CreateProductBatches(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusCreated {
				var envelope struct {
					Data models.ProductBatchesResponse `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &envelope)
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseBody, envelope.Data)
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
