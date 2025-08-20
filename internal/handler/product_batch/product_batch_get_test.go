package handler_test

import (
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

func TestProductBatchesHandler_GetReportProduct(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		mockService     func() *mocks.ProductBatchServiceMock
		wantStatus      int
		wantResponse    any
		wantErrorCode   string
		wantErrorSubMsg string
	}{
		{
			name:  "success - returns all product batch reports",
			query: "",
			mockService: func() *mocks.ProductBatchServiceMock {
				mock := &mocks.ProductBatchServiceMock{}
				mock.FuncGetReport = func(ctx context.Context) ([]models.ReportProduct, error) {
					return testhelpers.DummyReportProductsList(), nil
				}
				return mock
			},
			wantStatus:   http.StatusOK,
			wantResponse: testhelpers.DummyReportProductsList(),
		},
		{
			name:  "success - returns product batch report by id",
			query: "id=33",
			mockService: func() *mocks.ProductBatchServiceMock {
				mock := &mocks.ProductBatchServiceMock{}
				mock.FuncGetReportById = func(ctx context.Context, id int) (*models.ReportProduct, error) {
					rep := testhelpers.DummyReportProduct()
					return &rep, nil
				}
				return mock
			},
			wantStatus:   http.StatusOK,
			wantResponse: testhelpers.DummyReportProduct(),
		},
		{
			name:  "error - invalid id query parameter",
			query: "id=pepe",
			mockService: func() *mocks.ProductBatchServiceMock {
				return &mocks.ProductBatchServiceMock{}
			},
			wantStatus:      http.StatusBadRequest,
			wantErrorCode:   apperrors.CodeBadRequest,
			wantErrorSubMsg: "invalid integer",
		},
		{
			name:  "error - service error when fetching all",
			query: "",
			mockService: func() *mocks.ProductBatchServiceMock {
				mock := &mocks.ProductBatchServiceMock{}
				mock.FuncGetReport = func(ctx context.Context) ([]models.ReportProduct, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "db fail")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorSubMsg: "db fail",
		},
		{
			name:  "error - service error when fetching by id",
			query: "id=33",
			mockService: func() *mocks.ProductBatchServiceMock {
				mock := &mocks.ProductBatchServiceMock{}
				mock.FuncGetReportById = func(ctx context.Context, id int) (*models.ReportProduct, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				return mock
			},
			wantStatus:      http.StatusNotFound,
			wantErrorCode:   apperrors.CodeNotFound,
			wantErrorSubMsg: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/product-batches/report"
			if tt.query != "" {
				url = url + "?" + tt.query
			}
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			h := handler.NewProductBatchesHandler(tt.mockService())
			// inject test logger to cover logger != nil branches
			h.SetLogger(testhelpers.NewTestLogger())

			h.GetReportProduct(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				if tt.query == "" {
					var envelope struct {
						Data []models.ReportProduct `json:"data"`
					}
					err = json.Unmarshal(rec.Body.Bytes(), &envelope)
					require.NoError(t, err)
					require.Equal(t, tt.wantResponse, envelope.Data)
				} else {
					var envelope struct {
						Data models.ReportProduct `json:"data"`
					}
					err = json.Unmarshal(rec.Body.Bytes(), &envelope)
					require.NoError(t, err)
					require.Equal(t, tt.wantResponse, envelope.Data)
				}
			} else {
				var body struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.wantErrorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrorSubMsg)
			}
		})
	}
}
