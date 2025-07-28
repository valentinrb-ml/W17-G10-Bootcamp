package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/buyer"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("update buyer successfully with all fields", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		newCardNumber := "CARD-UPDATED"
		newFirstName := "UpdatedJohn"
		newLastName := "UpdatedDoe"

		req := models.RequestBuyer{
			CardNumberId: &newCardNumber,
			FirstName:    &newFirstName,
			LastName:     &newLastName,
		}

		existingBuyer := testhelpers.CreateTestBuyerWithID(buyerID)

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			assert.Equal(t, newCardNumber, cardNumber)
			assert.Equal(t, buyerID, id)
			return false // Card number is available
		}

		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return existingBuyer, nil
		}

		mockRepo.MockUpdate = func(ctx context.Context, id int, b models.Buyer) error {
			assert.Equal(t, buyerID, id)
			// Verify the buyer was updated with new values
			assert.Equal(t, newCardNumber, b.CardNumberId)
			assert.Equal(t, newFirstName, b.FirstName)
			assert.Equal(t, newLastName, b.LastName)
			return nil
		}

		// Act
		result, err := buyerService.Update(ctx, buyerID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, buyerID, result.Id)
		assert.Equal(t, newCardNumber, result.CardNumberId)
		assert.Equal(t, newFirstName, result.FirstName)
		assert.Equal(t, newLastName, result.LastName)
	})

	t.Run("update buyer successfully with partial fields", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		newFirstName := "PartialUpdate"

		req := models.RequestBuyer{
			CardNumberId: nil, // No update to card number
			FirstName:    &newFirstName,
			LastName:     nil, // No update to last name
		}

		existingBuyer := &models.Buyer{
			Id:           buyerID,
			CardNumberId: "CARD-001",
			FirstName:    "John",
			LastName:     "Doe",
		}

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return existingBuyer, nil
		}

		mockRepo.MockUpdate = func(ctx context.Context, id int, b models.Buyer) error {
			assert.Equal(t, buyerID, id)
			// Verify only first name was updated
			assert.Equal(t, "CARD-001", b.CardNumberId) // Should remain unchanged
			assert.Equal(t, newFirstName, b.FirstName)  // Should be updated
			assert.Equal(t, "Doe", b.LastName)          // Should remain unchanged
			return nil
		}

		// Act
		result, err := buyerService.Update(ctx, buyerID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, buyerID, result.Id)
		assert.Equal(t, "CARD-001", result.CardNumberId)
		assert.Equal(t, newFirstName, result.FirstName)
		assert.Equal(t, "Doe", result.LastName)
	})

	t.Run("return conflict error when card number already exists for another buyer", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		conflictingCardNumber := "CARD-EXISTING"

		req := models.RequestBuyer{
			CardNumberId: &conflictingCardNumber,
			FirstName:    testhelpers.PtrBuyer("John"),
			LastName:     testhelpers.PtrBuyer("Doe"),
		}

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			assert.Equal(t, conflictingCardNumber, cardNumber)
			assert.Equal(t, buyerID, id)
			return true // Card number exists for another buyer
		}

		// FindById and Update should not be called when card number conflicts
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			t.Error("FindById should not be called when card number conflicts")
			return nil, nil
		}

		mockRepo.MockUpdate = func(ctx context.Context, id int, b models.Buyer) error {
			t.Error("Update should not be called when card number conflicts")
			return nil
		}

		// Act
		result, err := buyerService.Update(ctx, buyerID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		var appErr *apperrors.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperrors.CodeConflict, appErr.Code)
		assert.Contains(t, appErr.Message, "already in use by another buyer")
	})

	t.Run("return error when buyer not found", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 999
		req := testhelpers.DummyRequestBuyer()
		expectedError := errors.New("buyer not found")

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			return false // Card number is available
		}

		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return nil, expectedError
		}

		mockRepo.MockUpdate = func(ctx context.Context, id int, b models.Buyer) error {
			t.Error("Update should not be called when buyer not found")
			return nil
		}

		// Act
		result, err := buyerService.Update(ctx, buyerID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("return error when repository update fails", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		req := testhelpers.DummyRequestBuyer()
		existingBuyer := testhelpers.CreateTestBuyerWithID(buyerID)
		expectedError := errors.New("database update failed")

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			return false
		}

		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			return existingBuyer, nil
		}

		mockRepo.MockUpdate = func(ctx context.Context, id int, b models.Buyer) error {
			return expectedError
		}

		// Act
		result, err := buyerService.Update(ctx, buyerID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("update buyer without card number validation when card number is not provided", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		newFirstName := "OnlyFirstName"

		req := models.RequestBuyer{
			CardNumberId: nil, // No card number update
			FirstName:    &newFirstName,
			LastName:     nil,
		}

		existingBuyer := testhelpers.CreateTestBuyerWithID(buyerID)

		// Mock repository methods - CardNumberExists should NOT be called
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			t.Error("CardNumberExists should not be called when CardNumberId is nil")
			return false
		}

		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			return existingBuyer, nil
		}

		mockRepo.MockUpdate = func(ctx context.Context, id int, b models.Buyer) error {
			assert.Equal(t, newFirstName, b.FirstName)
			return nil
		}

		// Act
		result, err := buyerService.Update(ctx, buyerID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newFirstName, result.FirstName)
	})

	t.Run("handle empty request gracefully", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		req := models.RequestBuyer{
			CardNumberId: nil,
			FirstName:    nil,
			LastName:     nil,
		}

		existingBuyer := testhelpers.CreateTestBuyerWithID(buyerID)

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			return existingBuyer, nil
		}

		mockRepo.MockUpdate = func(ctx context.Context, id int, b models.Buyer) error {
			// Verify no fields were changed
			assert.Equal(t, existingBuyer.CardNumberId, b.CardNumberId)
			assert.Equal(t, existingBuyer.FirstName, b.FirstName)
			assert.Equal(t, existingBuyer.LastName, b.LastName)
			return nil
		}

		// Act
		result, err := buyerService.Update(ctx, buyerID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		// All fields should remain the same
		assert.Equal(t, existingBuyer.CardNumberId, result.CardNumberId)
		assert.Equal(t, existingBuyer.FirstName, result.FirstName)
		assert.Equal(t, existingBuyer.LastName, result.LastName)
	})
}

// Benchmark test for Update method
func BenchmarkBuyerService_Update(b *testing.B) {
	ctx := context.Background()
	mockRepo := &mocks.BuyerRepositoryMocks{}
	buyerService := service.NewBuyerService(mockRepo)

	buyerID := 1
	req := testhelpers.DummyRequestBuyer()
	existingBuyer := testhelpers.CreateTestBuyerWithID(buyerID)

	mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
		return false
	}

	mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
		return existingBuyer, nil
	}

	mockRepo.MockUpdate = func(ctx context.Context, id int, buyer models.Buyer) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buyerService.Update(ctx, buyerID, req)
	}
}
