package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

type SellerService interface {
	Create(ctx context.Context, reqs models.RequestSeller) (*models.ResponseSeller, error)
	Update(ctx context.Context, id int, reqs models.RequestSeller) (*models.ResponseSeller, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]models.ResponseSeller, error)
	FindById(ctx context.Context, id int) (*models.ResponseSeller, error)
}

type sellerService struct {
	rp repository.SellerRepository
}

func NewSellerService(rp repository.SellerRepository) SellerService {
	return &sellerService{rp: rp}
}
