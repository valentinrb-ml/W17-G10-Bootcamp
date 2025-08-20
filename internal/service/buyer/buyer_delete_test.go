package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/buyer"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete buyer successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)
		buyerService.SetLogger(testhelpers.NewTestLogger())

		buyerID := 1

		// Mock repository methods
		mockRepo.MockDelete = func(ctx context.Context, id int) error {
			assert.Equal(t, buyerID, id)
			return nil
		}

		// Act
		err := buyerService.Delete(ctx, buyerID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("return error when repository delete fails", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)
		buyerService.SetLogger(testhelpers.NewTestLogger())

		buyerID := 1
		expectedError := errors.New("database delete failed")

		// Mock repository methods
		mockRepo.MockDelete = func(ctx context.Context, id int) error {
			assert.Equal(t, buyerID, id)
			return expectedError
		}

		// Act
		err := buyerService.Delete(ctx, buyerID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("handle delete with zero ID", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)
		buyerService.SetLogger(testhelpers.NewTestLogger())

		buyerID := 0

		// Mock repository methods
		mockRepo.MockDelete = func(ctx context.Context, id int) error {
			assert.Equal(t, buyerID, id)
			return nil
		}

		// Act
		err := buyerService.Delete(ctx, buyerID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("handle delete with negative ID", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)
		buyerService.SetLogger(testhelpers.NewTestLogger())

		buyerID := -1

		// Mock repository methods
		mockRepo.MockDelete = func(ctx context.Context, id int) error {
			assert.Equal(t, buyerID, id)
			return nil
		}

		// Act
		err := buyerService.Delete(ctx, buyerID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("return error when buyer not found", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 999
		expectedError := errors.New("buyer not found")

		// Mock repository methods
		mockRepo.MockDelete = func(ctx context.Context, id int) error {
			assert.Equal(t, buyerID, id)
			return expectedError
		}

		// Act
		err := buyerService.Delete(ctx, buyerID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

// Benchmark test for Delete method
func BenchmarkBuyerService_Delete(b *testing.B) {
	ctx := context.Background()
	mockRepo := &mocks.BuyerRepositoryMocks{}
	buyerService := service.NewBuyerService(mockRepo)

	buyerID := 1

	mockRepo.MockDelete = func(ctx context.Context, id int) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buyerService.Delete(ctx, buyerID)
	}
}
