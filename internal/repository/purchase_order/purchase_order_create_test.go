package repository_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderRepository_Create(t *testing.T) {
	insertQuery := regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES (?, ?, ?, ?, ?)")

	testCases := []struct {
		name           string
		setupMock      func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder)
		expectedError  bool
		errorCode      string
		errorMessage   string
		validateResult func(t *testing.T, result *models.PurchaseOrder, original models.PurchaseOrder)
	}{
		{
			name: "Success",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder) {
				expectedID := int64(1)
				helper.MockCreatePurchaseOrderSuccess(po, expectedID)
			},
			expectedError: false,
			validateResult: func(t *testing.T, result *models.PurchaseOrder, original models.PurchaseOrder) {
				assert.NotNil(t, result)
				assert.Equal(t, 1, result.ID)
				assert.Equal(t, original.OrderNumber, result.OrderNumber)
				assert.Equal(t, original.OrderDate, result.OrderDate)
				assert.Equal(t, original.TrackingCode, result.TrackingCode)
				assert.Equal(t, original.BuyerID, result.BuyerID)
				assert.Equal(t, original.ProductRecordID, result.ProductRecordID)
			},
		},
		{
			name: "Error_BuyerNotFound",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder) {
				helper.MockBuyerExists(po.BuyerID, false)
			},
			expectedError: true,
			errorCode:     apperrors.CodeNotFound,
			errorMessage:  "buyer with id 1 does not exist",
		},
		{
			name: "Error_ProductRecordNotFound",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder) {
				helper.MockBuyerExists(po.BuyerID, true)
				helper.MockProductRecordExists(po.ProductRecordID, false)
			},
			expectedError: true,
			errorCode:     apperrors.CodeNotFound,
			errorMessage:  "product record with id 1 does not exist",
		},
		{
			name: "Error_OrderNumberAlreadyExists",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder) {
				helper.MockBuyerExists(po.BuyerID, true)
				helper.MockProductRecordExists(po.ProductRecordID, true)
				helper.MockOrderNumberExists(po.OrderNumber, true)
			},
			expectedError: true,
			errorCode:     apperrors.CodeConflict,
			errorMessage:  "order_number already exists",
		},
		{
			name: "Error_DatabaseError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder) {
				helper.MockBuyerExists(po.BuyerID, true)
				helper.MockProductRecordExists(po.ProductRecordID, true)
				helper.MockOrderNumberExists(po.OrderNumber, false)

				helper.Mock.ExpectExec(insertQuery).
					WithArgs(po.OrderNumber, po.OrderDate, po.TrackingCode, po.BuyerID, po.ProductRecordID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorCode:     apperrors.CodeInternal,
			errorMessage:  "error creating purchase order",
		},
		{
			name: "Error_MySQLError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder) {
				helper.MockBuyerExists(po.BuyerID, true)
				helper.MockProductRecordExists(po.ProductRecordID, true)
				helper.MockOrderNumberExists(po.OrderNumber, false)

				mysqlErr := &mysql.MySQLError{
					Number:  1062,
					Message: "Duplicate entry",
				}

				helper.Mock.ExpectExec(insertQuery).
					WithArgs(po.OrderNumber, po.OrderDate, po.TrackingCode, po.BuyerID, po.ProductRecordID).
					WillReturnError(mysqlErr)
			},
			expectedError: true,
			errorCode:     apperrors.CodeInternal,
			errorMessage:  "database error: Duplicate entry",
		},
		{
			name: "Error_LastInsertIdError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper, po models.PurchaseOrder) {
				helper.MockBuyerExists(po.BuyerID, true)
				helper.MockProductRecordExists(po.ProductRecordID, true)
				helper.MockOrderNumberExists(po.OrderNumber, false)

				result := sqlmock.NewErrorResult(sql.ErrConnDone)
				helper.Mock.ExpectExec(insertQuery).
					WithArgs(po.OrderNumber, po.OrderDate, po.TrackingCode, po.BuyerID, po.ProductRecordID).
					WillReturnResult(result)
			},
			expectedError: true,
			errorCode:     apperrors.CodeInternal,
			errorMessage:  "error getting last insert id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			helper, err := testhelpers.NewTestPurchaseOrderHelper()
			assert.NoError(t, err)
			defer helper.Close()

			po := helper.CreateValidPurchaseOrder()
			tc.setupMock(helper, po)

			// Act
			result, err := helper.Repo.Create(helper.Ctx, po)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)

				appErr, ok := err.(*apperrors.AppError)
				assert.True(t, ok)
				assert.Equal(t, tc.errorCode, appErr.Code)
				assert.Contains(t, appErr.Message, tc.errorMessage)
			} else {
				assert.NoError(t, err)
				if tc.validateResult != nil {
					tc.validateResult(t, result, po)
				}
			}

			assert.NoError(t, helper.AssertExpectations())
		})
	}
}
