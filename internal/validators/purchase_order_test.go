package validators_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func validPurchaseOrder() models.RequestPurchaseOrder {
	return models.RequestPurchaseOrder{
		OrderNumber:     "ORDER123",
		OrderDate:       time.Now().Format(time.RFC3339),
		TrackingCode:    "TRACK123",
		BuyerID:         10,
		ProductRecordID: 20,
	}
}

func TestValidatePurchaseOrderPost_Valid(t *testing.T) {
	po := validPurchaseOrder()
	err := validators.ValidatePurchaseOrderPost(po)
	assert.Nil(t, err)
}

func TestValidatePurchaseOrderPost_MissingOrderNumber(t *testing.T) {
	po := validPurchaseOrder()
	po.OrderNumber = ""
	err := validators.ValidatePurchaseOrderPost(po)
	assert.NotNil(t, err)
	appErr, ok := err.(*apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.CodeValidationError, appErr.Code)
	assert.Equal(t, "order_number is required", appErr.Message)
}

func TestValidatePurchaseOrderPost_MissingOrderDate(t *testing.T) {
	po := validPurchaseOrder()
	po.OrderDate = ""
	err := validators.ValidatePurchaseOrderPost(po)
	assert.NotNil(t, err)
	appErr, ok := err.(*apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.CodeValidationError, appErr.Code)
	assert.Equal(t, "order_date is required", appErr.Message)
}

func TestValidatePurchaseOrderPost_InvalidOrderDateFormat(t *testing.T) {
	po := validPurchaseOrder()
	po.OrderDate = "not-a-date"
	err := validators.ValidatePurchaseOrderPost(po)
	assert.NotNil(t, err)
	appErr, ok := err.(*apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.CodeBadRequest, appErr.Code)
	assert.Equal(t, "Invalid request body", appErr.Message)
}

func TestValidatePurchaseOrderPost_MissingTrackingCode(t *testing.T) {
	po := validPurchaseOrder()
	po.TrackingCode = ""
	err := validators.ValidatePurchaseOrderPost(po)
	assert.NotNil(t, err)
	appErr, ok := err.(*apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.CodeValidationError, appErr.Code)
	assert.Equal(t, "tracking_code is required", appErr.Message)
}

func TestValidatePurchaseOrderPost_BuyerIDZeroOrNegative(t *testing.T) {
	cases := []int{0, -1}
	for _, v := range cases {
		po := validPurchaseOrder()
		po.BuyerID = v
		err := validators.ValidatePurchaseOrderPost(po)
		assert.NotNil(t, err)
		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, apperrors.CodeValidationError, appErr.Code)
		assert.Equal(t, "buyer_id must be greater than 0", appErr.Message)
	}
}

func TestValidatePurchaseOrderPost_ProductRecordIDZeroOrNegative(t *testing.T) {
	cases := []int64{0, -5}
	for _, v := range cases {
		po := validPurchaseOrder()
		po.ProductRecordID = int(v)
		err := validators.ValidatePurchaseOrderPost(po)
		assert.NotNil(t, err)
		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, apperrors.CodeValidationError, appErr.Code)
		assert.Equal(t, "product_record_id must be greater than 0", appErr.Message)
	}
}
