package repository

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type BuyerRepository interface {
	Create(ctx context.Context, s models.Buyer) (*models.Buyer, error)
	Update(ctx context.Context, id int, b models.Buyer) error
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]models.Buyer, error)
	FindById(ctx context.Context, id int) (*models.Buyer, error)
	CardNumberExists(ctx context.Context, cardNumber string, excludeId int) bool

	// SetLogger allows injecting the logger after creation
	SetLogger(l logger.Logger)
}

type buyerRepository struct {
	mysql  *sql.DB
	logger logger.Logger
}

func NewBuyerRepository(mysql *sql.DB) BuyerRepository {
	return &buyerRepository{
		mysql: mysql,
	}
}

// SetLogger allows you to inject the logger after creation
func (r *buyerRepository) SetLogger(l logger.Logger) {
	r.logger = l
}
