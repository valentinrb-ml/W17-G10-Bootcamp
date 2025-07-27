package repository_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderRepository_GetByID(t *testing.T) {
	query := regexp.QuoteMeta("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders WHERE id = ?")

	testCases := []struct {
		name           string
		orderID        int
		setupMock      func(helper *testhelpers.TestPurchaseOrderHelper, orderID int)
		expectedResult *models.PurchaseOrder
		expectedError  bool
		errorCode      string
		errorMessage   string
	}{
		{
			name:    "Success",
			orderID: 1,
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, orderID int) {
				expectedPO := models.PurchaseOrder{
					ID:              1,
					OrderNumber:     "ORD-001",
					OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
					TrackingCode:    "TRK-001",
					BuyerID:         1,
					ProductRecordID: 1,
				}
				helper.MockGetPurchaseOrderByIDSuccess(expectedPO)
			},
			expectedResult: &models.PurchaseOrder{
				ID:              1,
				OrderNumber:     "ORD-001",
				OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				TrackingCode:    "TRK-001",
				BuyerID:         1,
				ProductRecordID: 1,
			},
			expectedError: false,
		},
		{
			name:    "Error_NotFound",
			orderID: 999,
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, orderID int) {
				helper.MockGetPurchaseOrderByIDNotFound(orderID)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeNotFound,
			errorMessage:   "purchase order not found",
		},
		{
			name:    "Error_QueryError",
			orderID: 1,
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, orderID int) {
				helper.Mock.ExpectQuery(query).
					WithArgs(orderID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeInternal,
			errorMessage:   "error querying purchase order by id",
		},
		{
			name:    "Error_DateParseError",
			orderID: 1,
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, orderID int) {
				// Mock with invalid date format
				rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
					AddRow(1, "ORD-001", "invalid-date", "TRK-001", 1, 1)

				helper.Mock.ExpectQuery(query).
					WithArgs(orderID).
					WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeInternal,
			errorMessage:   "error parsing order date",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			helper, err := testhelpers.NewTestPurchaseOrderHelper()
			assert.NoError(t, err)
			defer helper.Close()

			tc.setupMock(helper, tc.orderID)

			// Act
			result, err := helper.Repo.GetByID(helper.Ctx, tc.orderID)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)

				appErr, ok := err.(*apperrors.AppError)
				assert.True(t, ok)
				assert.Equal(t, tc.errorCode, appErr.Code)
				assert.Equal(t, tc.errorMessage, appErr.Message)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult.ID, result.ID)
				assert.Equal(t, tc.expectedResult.OrderNumber, result.OrderNumber)
				assert.Equal(t, tc.expectedResult.OrderDate, result.OrderDate)
				assert.Equal(t, tc.expectedResult.TrackingCode, result.TrackingCode)
				assert.Equal(t, tc.expectedResult.BuyerID, result.BuyerID)
				assert.Equal(t, tc.expectedResult.ProductRecordID, result.ProductRecordID)
			}

			assert.NoError(t, helper.AssertExpectations())
		})
	}
}
