package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestPurchaseOrderHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.PurchaseOrderServiceMock)
		expectedStatus int
		expectedError  *apperrors.AppError
	}{
		{
			name: "success - create purchase order",
			requestBody: `{
				"data": {
					"order_number": "PO-001",
					"order_date": "2023-01-01T00:00:00Z",
					"tracking_code": "TRACK001",
					"buyer_id": 101,
					"product_record_id": 201
				}
			}`,
			mockSetup: func(m *mocks.PurchaseOrderServiceMock) {
				m.CreateFn = func(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
					return &models.ResponsePurchaseOrder{
						ID:              1,
						OrderNumber:     req.OrderNumber,
						OrderDate:       req.OrderDate,
						TrackingCode:    req.TrackingCode,
						BuyerID:         req.BuyerID,
						ProductRecordID: req.ProductRecordID,
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  nil,
		},
		{
			name: "error - invalid json",
			requestBody: `{
				"data": {
					"order_number": "PO-001",
					"order_date": "invalid-date",
					"tracking_code": "TRACK001",
					"buyer_id": 101,
					"product_record_id": 201
				}
			}`,
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"),
		},
		{
			name:           "error - invalid json syntax (JSON mal cerrado)",
			requestBody:    `{ "data": { "order_number": "PO-xxx" `, // JSON inválido
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"),
		},
		{
			name: "error - service error",
			requestBody: `{
				"data": {
					"order_number": "PO-002",
					"order_date": "2023-01-01T00:00:00Z",
					"tracking_code": "TRACK002",
					"buyer_id": 102,
					"product_record_id": 202
				}
			}`,
			mockSetup: func(m *mocks.PurchaseOrderServiceMock) {
				m.CreateFn = func(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
					return nil, apperrors.NewAppError(apperrors.CodeConflict, "order already exists")
				}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  apperrors.NewAppError(apperrors.CodeConflict, "order already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.PurchaseOrderServiceMock{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			} else {
				mockService.CreateFn = func(context.Context, models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
					return nil, nil
				}
			}

			h := handler.NewPurchaseOrderHandler(mockService)
			h.SetLogger(testhelpers.NewTestLogger())

			req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Create(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != nil {
				var errWrap struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				}
				err := json.NewDecoder(resp.Body).Decode(&errWrap)
				require.NoError(t, err)
				require.Equal(t, tt.expectedError.Code, errWrap.Error.Code)
				require.Equal(t, tt.expectedError.Message, errWrap.Error.Message)
			}
		})
	}
}
