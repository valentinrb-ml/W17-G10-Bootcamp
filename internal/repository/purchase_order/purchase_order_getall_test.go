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

func TestPurchaseOrderRepository_GetAll(t *testing.T) {
	query := regexp.QuoteMeta("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders")

	testCases := []struct {
		name           string
		setupMock      func(helper *testhelpers.TestPurchaseOrderHelper)
		expectedResult []models.PurchaseOrder
		expectedError  bool
		errorCode      string
		errorMessage   string
	}{
		{
			name: "Success_WithResults",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				expectedPOs := []models.PurchaseOrder{
					{
						ID:              1,
						OrderNumber:     "ORD-001",
						OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
						TrackingCode:    "TRK-001",
						BuyerID:         1,
						ProductRecordID: 1,
					},
					{
						ID:              2,
						OrderNumber:     "ORD-002",
						OrderDate:       time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC),
						TrackingCode:    "TRK-002",
						BuyerID:         2,
						ProductRecordID: 2,
					},
				}
				helper.MockGetAllPurchaseOrdersSuccess(expectedPOs)
			},
			expectedResult: []models.PurchaseOrder{
				{
					ID:              1,
					OrderNumber:     "ORD-001",
					OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
					TrackingCode:    "TRK-001",
					BuyerID:         1,
					ProductRecordID: 1,
				},
				{
					ID:              2,
					OrderNumber:     "ORD-002",
					OrderDate:       time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC),
					TrackingCode:    "TRK-002",
					BuyerID:         2,
					ProductRecordID: 2,
				},
			},
			expectedError: false,
		},
		{
			name: "Success_EmptyResult",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				helper.MockGetAllPurchaseOrdersSuccess([]models.PurchaseOrder{})
			},
			expectedResult: []models.PurchaseOrder{},
			expectedError:  false,
		},
		{
			name: "Error_QueryError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				helper.Mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeInternal,
			errorMessage:   "error querying all purchase orders",
		},
		{
			name: "Error_ScanError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				// Mock with wrong number of columns to cause scan error
				rows := sqlmock.NewRows([]string{"id", "order_number"}).
					AddRow(1, "ORD-001")
				helper.Mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeInternal,
			errorMessage:   "error scanning purchase order",
		},
		{
			name: "Error_DateParseError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				// Mock with invalid date format
				rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
					AddRow(1, "ORD-001", "invalid-date", "TRK-001", 1, 1)
				helper.Mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeInternal,
			errorMessage:   "error parsing order date",
		},
		{
			name: "Error_RowsError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
					AddRow(1, "ORD-001", "2024-01-15 10:30:00", "TRK-001", 1, 1).
					RowError(0, sql.ErrConnDone)
				helper.Mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeInternal,
			errorMessage:   "error after iterating rows",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			helper, err := testhelpers.NewTestPurchaseOrderHelper()
			assert.NoError(t, err)
			defer helper.Close()

			tc.setupMock(helper)

			// Act
			result, err := helper.Repo.GetAll(helper.Ctx)

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

				if len(tc.expectedResult) == 0 {
					// Para resultados vacíos, permitir tanto nil como slice vacío
					if result == nil {
						assert.Nil(t, result)
					} else {
						assert.NotNil(t, result)
						assert.Empty(t, result)
						assert.Len(t, result, 0)
					}
				} else {
					assert.NotNil(t, result)
					assert.Len(t, result, len(tc.expectedResult))

					for i, po := range result {
						assert.Equal(t, tc.expectedResult[i].ID, po.ID)
						assert.Equal(t, tc.expectedResult[i].OrderNumber, po.OrderNumber)
						assert.Equal(t, tc.expectedResult[i].OrderDate, po.OrderDate)
						assert.Equal(t, tc.expectedResult[i].TrackingCode, po.TrackingCode)
						assert.Equal(t, tc.expectedResult[i].BuyerID, po.BuyerID)
						assert.Equal(t, tc.expectedResult[i].ProductRecordID, po.ProductRecordID)
					}
				}
			}

			assert.NoError(t, helper.AssertExpectations())
		})
	}
}
