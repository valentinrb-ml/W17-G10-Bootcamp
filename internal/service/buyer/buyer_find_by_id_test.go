package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/buyer"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/buyer"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestBuyerService_FindById(t *testing.T) {
	ctx := context.Background()

	t.Run("find buyer by ID successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		expectedBuyer := testhelpers.CreateTestBuyerWithID(buyerID)
		expectedResponse := testhelpers.DummyResponseBuyer()

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return expectedBuyer, nil
		}

		// Act
		result, err := buyerService.FindById(ctx, buyerID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResponse.Id, result.Id)
		assert.Equal(t, expectedResponse.CardNumberId, result.CardNumberId)
		assert.Equal(t, expectedResponse.FirstName, result.FirstName)
		assert.Equal(t, expectedResponse.LastName, result.LastName)
	})

	t.Run("eturn error when buyer not found", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 999
		expectedError := errors.New("buyer not found")

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return nil, expectedError
		}

		// Act
		result, err := buyerService.FindById(ctx, buyerID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1
		expectedError := errors.New("database connection failed")

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return nil, expectedError
		}

		// Act
		result, err := buyerService.FindById(ctx, buyerID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("handle zero ID", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 0
		expectedError := errors.New("invalid ID")

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return nil, expectedError
		}

		// Act
		result, err := buyerService.FindById(ctx, buyerID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("handle negative ID", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := -1
		expectedError := errors.New("invalid ID")

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return nil, expectedError
		}

		// Act
		result, err := buyerService.FindById(ctx, buyerID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("find buyer with all fields populated", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 123
		expectedBuyer := &models.Buyer{
			Id:           buyerID,
			CardNumberId: "CARD-SPECIAL",
			FirstName:    "Special",
			LastName:     "User",
		}

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return expectedBuyer, nil
		}

		// Act
		result, err := buyerService.FindById(ctx, buyerID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBuyer.Id, result.Id)
		assert.Equal(t, expectedBuyer.CardNumberId, result.CardNumberId)
		assert.Equal(t, expectedBuyer.FirstName, result.FirstName)
		assert.Equal(t, expectedBuyer.LastName, result.LastName)
	})

	t.Run("panic when repository returns nil buyer without error", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyerID := 1

		// Mock repository methods - returns nil buyer but no error
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			assert.Equal(t, buyerID, id)
			return nil, nil
		}

		// Act & Assert
		// This demonstrates a bug in the service implementation - it should handle nil pointer
		assert.Panics(t, func() {
			buyerService.FindById(ctx, buyerID)
		}, "Service should handle nil buyer gracefully but currently panics")
	})

	t.Run("find multiple different buyers by different IDs", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		testCases := []struct {
			id           int
			expectedName string
		}{
			{1, "John"},
			{2, "Jane"},
			{3, "Bob"},
		}

		// Mock repository methods
		mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
			switch id {
			case 1:
				return &models.Buyer{Id: 1, CardNumberId: "CARD-001", FirstName: "John", LastName: "Doe"}, nil
			case 2:
				return &models.Buyer{Id: 2, CardNumberId: "CARD-002", FirstName: "Jane", LastName: "Smith"}, nil
			case 3:
				return &models.Buyer{Id: 3, CardNumberId: "CARD-003", FirstName: "Bob", LastName: "Johnson"}, nil
			default:
				return nil, errors.New("buyer not found")
			}
		}

		// Act & Assert
		for _, tc := range testCases {
			result, err := buyerService.FindById(ctx, tc.id)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.id, result.Id)
			assert.Equal(t, tc.expectedName, result.FirstName)
		}
	})
}

// Benchmark test for FindById method
func BenchmarkBuyerService_FindById(b *testing.B) {
	ctx := context.Background()
	mockRepo := &mocks.BuyerRepositoryMocks{}
	buyerService := service.NewBuyerService(mockRepo)

	buyerID := 1
	expectedBuyer := testhelpers.CreateTestBuyerWithID(buyerID)

	mockRepo.MockFindById = func(ctx context.Context, id int) (*models.Buyer, error) {
		return expectedBuyer, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buyerService.FindById(ctx, buyerID)
	}
}
