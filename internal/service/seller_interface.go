package service

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

type SellerService interface {
	Create(reqs models.RequestSeller) (*models.ResponseSeller, error)
	Update(id int, reqs models.RequestSeller) (*models.ResponseSeller, error)
	Delete(id int) error
	FindAll() ([]models.ResponseSeller, error)
	FindById(id int) (*models.ResponseSeller, error)
}

type sellerService struct {
	rp repository.SellerRepository
}

func NewSellerService(rp repository.SellerRepository) SellerService {
	return &sellerService{rp: rp}
}
