package service

import (
	"context"

	buyerRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type BuyerService interface {
	Create(ctx context.Context, reqs models.RequestBuyer) (*models.ResponseBuyer, error)
	Update(ctx context.Context, id int, reqs models.RequestBuyer) (*models.ResponseBuyer, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]models.ResponseBuyer, error)
	FindById(ctx context.Context, id int) (*models.ResponseBuyer, error)

	// SetLogger allows injecting the logger after creation
	SetLogger(l logger.Logger)
}

type buyerService struct {
	rp     buyerRepo.BuyerRepository
	logger logger.Logger
}

func NewBuyerService(rp buyerRepo.BuyerRepository) BuyerService {
	return &buyerService{rp: rp}
}

// SetLogger allows you to inject the logger after creation
func (s *buyerService) SetLogger(l logger.Logger) {
	s.logger = l
	s.rp.SetLogger(l) // También inyectar el logger al repository
}
