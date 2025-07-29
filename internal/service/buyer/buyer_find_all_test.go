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

func TestBuyerService_FindAll(t *testing.T) {
	ctx := context.Background()

	t.Run("find all buyers successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		expectedBuyers := testhelpers.FindAllBuyersDummy()
		expectedResponse := testhelpers.FindAllBuyersResponseDummy()

		// Mock repository methods
		mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
			return expectedBuyers, nil
		}

		// Act
		result, err := buyerService.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedResponse), len(result))

		// Verify the results are sorted by ID
		for i := 1; i < len(result); i++ {
			assert.True(t, result[i-1].Id <= result[i].Id, "Results should be sorted by ID")
		}

		// Verify first buyer details
		if len(result) > 0 {
			assert.Equal(t, expectedResponse[0].Id, result[0].Id)
			assert.Equal(t, expectedResponse[0].CardNumberId, result[0].CardNumberId)
			assert.Equal(t, expectedResponse[0].FirstName, result[0].FirstName)
			assert.Equal(t, expectedResponse[0].LastName, result[0].LastName)
		}
	})

	t.Run("return empty slice when no buyers found", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		// Mock repository methods
		mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
			return []models.Buyer{}, nil
		}

		// Act
		result, err := buyerService.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("return empty slice when repository returns nil", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		// Mock repository methods
		mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
			return nil, nil
		}

		// Act
		result, err := buyerService.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		expectedError := errors.New("database connection failed")

		// Mock repository methods
		mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
			return nil, expectedError
		}

		// Act
		result, err := buyerService.FindAll(ctx)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("handle single buyer correctly", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		singleBuyer := []models.Buyer{
			{
				Id:           1,
				CardNumberId: "CARD-001",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		// Mock repository methods
		mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
			return singleBuyer, nil
		}

		// Act
		result, err := buyerService.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, singleBuyer[0].Id, result[0].Id)
		assert.Equal(t, singleBuyer[0].CardNumberId, result[0].CardNumberId)
		assert.Equal(t, singleBuyer[0].FirstName, result[0].FirstName)
		assert.Equal(t, singleBuyer[0].LastName, result[0].LastName)
	})

	t.Run("sort buyers by ID correctly", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		// Create buyers in random order
		unsortedBuyers := []models.Buyer{
			{Id: 3, CardNumberId: "CARD-003", FirstName: "Charlie", LastName: "Brown"},
			{Id: 1, CardNumberId: "CARD-001", FirstName: "Alice", LastName: "Smith"},
			{Id: 5, CardNumberId: "CARD-005", FirstName: "Eve", LastName: "Johnson"},
			{Id: 2, CardNumberId: "CARD-002", FirstName: "Bob", LastName: "Wilson"},
			{Id: 4, CardNumberId: "CARD-004", FirstName: "David", LastName: "Davis"},
		}

		// Mock repository methods
		mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
			return unsortedBuyers, nil
		}

		// Act
		result, err := buyerService.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, len(result))

		// Verify sorting
		expectedOrder := []int{1, 2, 3, 4, 5}
		for i, buyer := range result {
			assert.Equal(t, expectedOrder[i], buyer.Id)
		}
	})

	t.Run("handle buyers with duplicate IDs", func(t *testing.T) {
		// Arrange
		mockRepo := &mocks.BuyerRepositoryMocks{}
		buyerService := service.NewBuyerService(mockRepo)

		buyersWithDuplicates := []models.Buyer{
			{Id: 1, CardNumberId: "CARD-001", FirstName: "John", LastName: "Doe"},
			{Id: 1, CardNumberId: "CARD-002", FirstName: "Jane", LastName: "Smith"},
			{Id: 2, CardNumberId: "CARD-003", FirstName: "Bob", LastName: "Johnson"},
		}

		// Mock repository methods
		mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
			return buyersWithDuplicates, nil
		}

		// Act
		result, err := buyerService.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 3, len(result))

		// Should maintain stable sort order
		assert.Equal(t, 1, result[0].Id)
		assert.Equal(t, 1, result[1].Id)
		assert.Equal(t, 2, result[2].Id)
	})
}

// Benchmark test for FindAll method
func BenchmarkBuyerService_FindAll(b *testing.B) {
	ctx := context.Background()
	mockRepo := &mocks.BuyerRepositoryMocks{}
	buyerService := service.NewBuyerService(mockRepo)

	buyers := testhelpers.CreateTestBuyers(100) // Create 100 test buyers

	mockRepo.MockFindAll = func(ctx context.Context) ([]models.Buyer, error) {
		return buyers, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buyerService.FindAll(ctx)
	}
}
