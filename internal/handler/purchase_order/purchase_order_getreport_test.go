package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/purchase_order"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderHandler_GetReport(t *testing.T) {
	type args struct {
		queryID string
	}
	buyerData := []models.BuyerWithPurchaseCount{
		{ID: 101, CardNumberID: "CARD101", FirstName: "John", LastName: "Doe", PurchaseOrdersCount: 3},
	}

	tests := []struct {
		name            string
		args            args
		mockSetup       func(*mocks.PurchaseOrderServiceMock)
		expectedStatus  int
		expectedError   *apperrors.AppError
		expectedResults []models.BuyerWithPurchaseCount
	}{
		{
			name: "ok - with buyerID param",
			args: args{queryID: "101"},
			mockSetup: func(m *mocks.PurchaseOrderServiceMock) {
				m.GetReportByBuyerFn = func(_ context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) {
					require.NotNil(t, buyerID)
					require.Equal(t, 101, *buyerID)
					return buyerData, nil
				}
			},
			expectedStatus:  http.StatusOK,
			expectedResults: buyerData,
		},
		{
			name: "ok - no buyerID param (global)",
			args: args{queryID: ""},
			mockSetup: func(m *mocks.PurchaseOrderServiceMock) {
				m.GetReportByBuyerFn = func(_ context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) {
					require.Nil(t, buyerID)
					return buyerData, nil
				}
			},
			expectedStatus:  http.StatusOK,
			expectedResults: buyerData,
		},
		{
			name:           "error - invalid buyerID param",
			args:           args{queryID: "badval"},
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid buyer ID parameter"),
		},
		{
			name: "error - service error",
			args: args{queryID: "101"},
			mockSetup: func(m *mocks.PurchaseOrderServiceMock) {
				m.GetReportByBuyerFn = func(_ context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  apperrors.NewAppError(apperrors.CodeNotFound, "not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSrv := &mocks.PurchaseOrderServiceMock{}
			// Para nil/mockSetup, garantiza que GetReportByBuyerFn no cause panic
			if tt.mockSetup != nil {
				tt.mockSetup(mockSrv)
			} else {
				mockSrv.GetReportByBuyerFn = func(_ context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) {
					return nil, errors.New("unexpected call in this test")
				}
			}
			h := handler.NewPurchaseOrderHandler(mockSrv)
			h.SetLogger(testhelpers.NewTestLogger())

			req := httptest.NewRequest(http.MethodGet, "/purchase-orders/report", nil)
			if tt.args.queryID != "" {
				q := req.URL.Query()
				q.Set("id", tt.args.queryID)
				req.URL.RawQuery = q.Encode()
			}
			w := httptest.NewRecorder()
			h.GetReport(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != nil {
				// Decodifica {"error": {...}}
				var errResp struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				}
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				require.Equal(t, tt.expectedError.Code, errResp.Error.Code)
				require.Equal(t, tt.expectedError.Message, errResp.Error.Message)
			}
			if tt.expectedResults != nil {
				var body struct {
					Data []models.BuyerWithPurchaseCount `json:"data"`
				}
				_ = json.NewDecoder(resp.Body).Decode(&body)
				require.Equal(t, tt.expectedResults, body.Data)
			}
		})
	}
}
