package repository

import (
	"context"
	"database/sql"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

type SellerRepository interface {
	Create(ctx context.Context, s models.Seller) (*models.Seller, error)
	Update(ctx context.Context, id int, s models.Seller) error
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]models.Seller, error)
	FindById(ctx context.Context, id int) (*models.Seller, error)
}

type sellerRepository struct {
	mysql *sql.DB
}

func NewSellerRepository(mysql *sql.DB) SellerRepository {
	return &sellerRepository{
		mysql: mysql,
	}
}
