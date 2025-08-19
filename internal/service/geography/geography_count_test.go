package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/geography"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func TestGeographyService_CountSellersGroupedByLocality(t *testing.T) {
	tests := []struct {
		name     string
		mockRepo func() *mocks.GeographyRepositoryMock
		wantErr  bool
		wantMsg  string
		wantResp []models.ResponseLocalitySellers
	}{
		{
			name: "success",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncCountSellersGroupedByLocality = func(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
					return []models.ResponseLocalitySellers{
						{LocalityId: "101", LocalityName: "CABA", SellersCount: 10},
						{LocalityId: "102", LocalityName: "Córdoba", SellersCount: 5},
					}, nil
				}
				return mock
			},
			wantErr: false,
			wantResp: []models.ResponseLocalitySellers{
				{LocalityId: "101", LocalityName: "CABA", SellersCount: 10},
				{LocalityId: "102", LocalityName: "Córdoba", SellersCount: 5},
			},
		},
		{
			name: "repo error",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncCountSellersGroupedByLocality = func(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
					return nil, errors.New("repo failure")
				}
				return mock
			},
			wantErr: true,
			wantMsg: "repo failure",
		},
		{
			name: "repo returns nil slice, nil error",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncCountSellersGroupedByLocality = func(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
					return nil, nil
				}
				return mock
			},
			wantErr:  false,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			srv := service.NewGeographyService(repo)
			resp, err := srv.CountSellersGroupedByLocality(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantMsg != "" {
					require.Contains(t, err.Error(), tt.wantMsg)
				}
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestGeographyService_CountSellersByLocality(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		mockRepo func() *mocks.GeographyRepositoryMock
		wantErr  bool
		wantMsg  string
		wantResp *models.ResponseLocalitySellers
	}{
		{
			name: "success",
			id:   "101",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncCountSellersByLocality = func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
					return &models.ResponseLocalitySellers{
						LocalityId: "101", LocalityName: "CABA", SellersCount: 20,
					}, nil
				}
				return mock
			},
			wantErr: false,
			wantResp: &models.ResponseLocalitySellers{
				LocalityId: "101", LocalityName: "CABA", SellersCount: 20,
			},
		},
		{
			name: "repo error",
			id:   "666",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncCountSellersByLocality = func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
					return nil, errors.New("no sellers")
				}
				return mock
			},
			wantErr: true,
			wantMsg: "no sellers",
		},
		{
			name: "repo returns nil,nil",
			id:   "999",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncCountSellersByLocality = func(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
					return nil, nil
				}
				return mock
			},
			wantErr:  false,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			srv := service.NewGeographyService(repo)
			resp, err := srv.CountSellersByLocality(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantMsg != "" {
					require.Contains(t, err.Error(), tt.wantMsg)
				}
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
