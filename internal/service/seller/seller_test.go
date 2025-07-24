package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/seller"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/seller"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func dummyRequestSeller() models.RequestSeller {
	cid := 101
	companyName := "Frutas del Sur"
	address := "Calle 1"
	telephone := "221-111"
	localityId := "1900"

	return models.RequestSeller{
		Cid:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		LocalityId:  &localityId,
	}
}

func dummySeller() models.Seller {
	return models.Seller{
		Id:          1,
		Cid:         101,
		CompanyName: "Frutas del Sur",
		Address:     "Calle 1",
		Telephone:   "221-111",
		LocalityId:  "1900",
	}
}

func dummyResponseSeller() models.ResponseSeller {
	return models.ResponseSeller{
		Id:          1,
		Cid:         101,
		CompanyName: "Frutas del Sur",
		Address:     "Calle 1",
		Telephone:   "221-111",
		LocalityId:  "1900",
	}
}

func TestSellerService_Create_Success(t *testing.T) {
	mockRepo := &mocks.SellerRepositoryMock{
		CreateFn: func(ctx context.Context, s models.Seller) (*models.Seller, error) {
			result := dummySeller()
			return &result, nil
		},
	}
	svc := service.NewSellerService(mockRepo, nil)

	req := dummyRequestSeller()
	expected := dummyResponseSeller()

	resp, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, &expected, resp)
}

func TestSellerService_Update_Success(t *testing.T) {
	mockRepo := &mocks.SellerRepositoryMock{
		FindByIdFn: func(ctx context.Context, id int) (*models.Seller, error) {
			result := dummySeller()
			return &result, nil
		},
		UpdateFn: func(ctx context.Context, id int, s models.Seller) error {
			return nil
		},
	}
	svc := service.NewSellerService(mockRepo, nil)

	req := dummyRequestSeller()
	expected := dummyResponseSeller()

	resp, err := svc.Update(context.Background(), 1, req)
	require.NoError(t, err)
	require.Equal(t, &expected, resp)
}

func TestSellerService_Delete_Success(t *testing.T) {
	mockRepo := &mocks.SellerRepositoryMock{
		DeleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	svc := service.NewSellerService(mockRepo, nil)

	err := svc.Delete(context.Background(), 1)
	require.NoError(t, err)
}

func TestSellerService_FindAll_Success(t *testing.T) {
	mockRepo := &mocks.SellerRepositoryMock{
		FindAllFn: func(ctx context.Context) ([]models.Seller, error) {
			return []models.Seller{dummySeller()}, nil
		},
	}
	svc := service.NewSellerService(mockRepo, nil)

	expected := []models.ResponseSeller{dummyResponseSeller()}

	resp, err := svc.FindAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestSellerService_FindById_Success(t *testing.T) {
	mockRepo := &mocks.SellerRepositoryMock{
		FindByIdFn: func(ctx context.Context, id int) (*models.Seller, error) {
			result := dummySeller()
			return &result, nil
		},
	}
	svc := service.NewSellerService(mockRepo, nil)

	expected := dummyResponseSeller()
	resp, err := svc.FindById(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, &expected, resp)
}
