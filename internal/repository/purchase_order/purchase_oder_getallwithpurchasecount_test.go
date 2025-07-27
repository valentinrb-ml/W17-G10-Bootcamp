package repository_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestPurchaseOrderRepository_GetAllWithPurchaseCount(t *testing.T) {
	query := regexp.QuoteMeta("SELECT b.id, b.id_card_number, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count FROM buyers b LEFT JOIN purchase_orders po ON b.id = po.buyer_id GROUP BY b.id")

	testCases := []struct {
		name           string
		setupMock      func(helper *testhelpers.TestPurchaseOrderHelper)
		expectedResult []models.BuyerWithPurchaseCount
		expectedError  bool
		errorCode      string
		errorMessage   string
	}{
		{
			name: "Success_WithResults",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				expectedBuyers := []models.BuyerWithPurchaseCount{
					{
						ID:                  1,
						CardNumberID:        "12345678",
						FirstName:           "John",
						LastName:            "Doe",
						PurchaseOrdersCount: 5,
					},
					{
						ID:                  2,
						CardNumberID:        "87654321",
						FirstName:           "Jane",
						LastName:            "Smith",
						PurchaseOrdersCount: 3,
					},
				}
				helper.MockGetAllWithPurchaseCountSuccess(expectedBuyers)
			},
			expectedResult: []models.BuyerWithPurchaseCount{
				{
					ID:                  1,
					CardNumberID:        "12345678",
					FirstName:           "John",
					LastName:            "Doe",
					PurchaseOrdersCount: 5,
				},
				{
					ID:                  2,
					CardNumberID:        "87654321",
					FirstName:           "Jane",
					LastName:            "Smith",
					PurchaseOrdersCount: 3,
				},
			},
			expectedError: false,
		},
		{
			name: "Success_EmptyResult",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				helper.MockGetAllWithPurchaseCountSuccess([]models.BuyerWithPurchaseCount{})
			},
			expectedResult: []models.BuyerWithPurchaseCount{},
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
			errorMessage:   "error querying all purchase counts",
		},
		{
			name: "Error_ScanError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				// Mock with wrong number of columns to cause scan error
				rows := sqlmock.NewRows([]string{"id", "id_card_number"}).
					AddRow(1, "12345678")

				helper.Mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  true,
			errorCode:      apperrors.CodeInternal,
			errorMessage:   "error scanning purchase count result",
		},
		{
			name: "Error_RowsError",
			setupMock: func(helper *testhelpers.TestPurchaseOrderHelper) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"}).
					AddRow(1, "12345678", "John", "Doe", 5).
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
			result, err := helper.Repo.GetAllWithPurchaseCount(helper.Ctx)

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

					for i, buyer := range result {
						expected := tc.expectedResult[i]
						assert.Equal(t, expected.ID, buyer.ID)
						assert.Equal(t, expected.CardNumberID, buyer.CardNumberID)
						assert.Equal(t, expected.FirstName, buyer.FirstName)
						assert.Equal(t, expected.LastName, buyer.LastName)
						assert.Equal(t, expected.PurchaseOrdersCount, buyer.PurchaseOrdersCount)
					}
				}
			}

			assert.NoError(t, helper.AssertExpectations())
		})
	}
}
