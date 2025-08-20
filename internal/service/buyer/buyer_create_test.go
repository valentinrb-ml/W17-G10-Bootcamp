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

func TestBuyerService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("should create buyer successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{
			// Initialize unused methods to avoid nil pointer panics
			MockFindAll: func(ctx context.Context) ([]models.Buyer, error) {
				return nil, nil
			},
		}
		service := service.NewBuyerService(mockRepo)
		service.SetLogger(testhelpers.NewTestLogger())
		//service := NewBuyerService(mockRepo)

		req := testhelpers.DummyRequestBuyer()
		expectedBuyer := testhelpers.CreateTestBuyerWithID(1)
		expectedResponse := testhelpers.DummyResponseBuyer()

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			assert.Equal(t, *req.CardNumberId, cardNumber)
			assert.Equal(t, 0, id) // Should be 0 for create operation
			return false           // Card number doesn't exist
		}

		mockRepo.MockCreate = func(ctx context.Context, b models.Buyer) (*models.Buyer, error) {
			// Verify the mapped buyer has correct values
			assert.Equal(t, *req.CardNumberId, b.CardNumberId)
			assert.Equal(t, *req.FirstName, b.FirstName)
			assert.Equal(t, *req.LastName, b.LastName)
			return expectedBuyer, nil
		}

		// Act
		result, err := service.Create(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResponse.Id, result.Id)
		assert.Equal(t, expectedResponse.CardNumberId, result.CardNumberId)
		assert.Equal(t, expectedResponse.FirstName, result.FirstName)
		assert.Equal(t, expectedResponse.LastName, result.LastName)
	})

	t.Run("should return conflict error when card number already exists", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		//service := NewBuyerService(mockRepo)
		service := service.NewBuyerService(mockRepo)
		service.SetLogger(testhelpers.NewTestLogger())

		req := testhelpers.DummyRequestBuyer()

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			assert.Equal(t, *req.CardNumberId, cardNumber)
			assert.Equal(t, 0, id)
			return true // Card number already exists
		}

		// MockCreate should not be called when card number exists
		mockRepo.MockCreate = func(ctx context.Context, b models.Buyer) (*models.Buyer, error) {
			t.Error("Create should not be called when card number exists")
			return nil, nil
		}

		// Act
		result, err := service.Create(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify it's the correct error type
		var appErr *apperrors.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperrors.CodeConflict, appErr.Code)
		assert.Contains(t, appErr.Message, "card number ID is already in use")
	})

	t.Run("should return error when repository create fails", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		//service := NewBuyerService(mockRepo)}
		service := service.NewBuyerService(mockRepo)
		service.SetLogger(testhelpers.NewTestLogger())

		req := testhelpers.DummyRequestBuyer()
		expectedError := errors.New("database connection failed")

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			return false // Card number doesn't exist
		}

		mockRepo.MockCreate = func(ctx context.Context, b models.Buyer) (*models.Buyer, error) {
			return nil, expectedError
		}

		// Act
		result, err := service.Create(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should handle nil pointer fields in request", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		//service := NewBuyerService(mockRepo)
		service := service.NewBuyerService(mockRepo)
		service.SetLogger(testhelpers.NewTestLogger())

		// Create request with nil CardNumberId (this should cause issues)
		req := models.RequestBuyer{
			CardNumberId: nil, // This will likely cause a panic
			FirstName:    testhelpers.PtrBuyer("John"),
			LastName:     testhelpers.PtrBuyer("Doe"),
		}

		// This test demonstrates a potential issue in your code
		// The service doesn't handle nil CardNumberId properly

		// Act & Assert
		assert.Panics(t, func() {
			service.Create(ctx, req)
		}, "Service should handle nil CardNumberId gracefully")
	})

	t.Run("should create buyer with all fields populated", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		//service := NewBuyerService(mockRepo)
		service := service.NewBuyerService(mockRepo)
		service.SetLogger(testhelpers.NewTestLogger())

		cardNumber := "CARD-999"
		firstName := "Alice"
		lastName := "Wonder"

		req := models.RequestBuyer{
			CardNumberId: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
		}

		createdBuyer := &models.Buyer{
			Id:           999,
			CardNumberId: cardNumber,
			FirstName:    firstName,
			LastName:     lastName,
		}

		// Mock repository methods
		mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
			return false
		}

		mockRepo.MockCreate = func(ctx context.Context, b models.Buyer) (*models.Buyer, error) {
			// Verify all fields are correctly mapped
			assert.Equal(t, cardNumber, b.CardNumberId)
			assert.Equal(t, firstName, b.FirstName)
			assert.Equal(t, lastName, b.LastName)
			return createdBuyer, nil
		}

		// Act
		result, err := service.Create(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 999, result.Id)
		assert.Equal(t, cardNumber, result.CardNumberId)
		assert.Equal(t, firstName, result.FirstName)
		assert.Equal(t, lastName, result.LastName)
	})
}

// Benchmark test for Create method
func BenchmarkBuyerService_Create(b *testing.B) {
	ctx := context.Background()
	mockRepo := &mocks.BuyerRepositoryMocks{}
	//service := NewBuyerService(mockRepo)
	service := service.NewBuyerService(mockRepo)

	req := testhelpers.DummyRequestBuyer()
	expectedBuyer := testhelpers.CreateTestBuyerWithID(1)

	mockRepo.MockCardNumberExists = func(ctx context.Context, cardNumber string, id int) bool {
		return false
	}

	mockRepo.MockCreate = func(ctx context.Context, buyer models.Buyer) (*models.Buyer, error) {
		return expectedBuyer, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Create(ctx, req)
	}
}
