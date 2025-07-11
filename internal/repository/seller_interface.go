package repository

import (
	"database/sql"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

type SellerRepository interface {
	Create(s models.Seller) (*models.Seller, error)
	Update(id int, s models.Seller) error
	Delete(id int) error
	FindAll() ([]models.Seller, error)
	FindById(id int) (*models.Seller, error)

	CIDExists(cid int, id int) bool
}

type sellerRepository struct {
	mysql *sql.DB
}

func NewSellerRepository(mysql *sql.DB) SellerRepository {
	return &sellerRepository{
		mysql: mysql,
	}
}
