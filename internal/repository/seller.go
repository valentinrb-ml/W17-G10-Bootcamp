package repository

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

type SellerRepository interface {
	Create(s models.Seller) models.Seller
	Update(id int, s models.Seller)
	Delete(id int)
	FindAll() []models.Seller
	FindById(id int) (*models.Seller, *api.ServiceError)

	IsValidCid(cid int) bool
	IsValidCidExcludeId(cid int, excludeId int) bool
	ExistsById(id int) bool
}

type sellerRepository struct {
	db map[int]models.Seller
}

func NewSellerRepository(db map[int]models.Seller) SellerRepository {
	defaultDb := make(map[int]models.Seller)
	if db != nil {
		defaultDb = db
	}
	return &sellerRepository{db: defaultDb}
}

func (r *sellerRepository) Create(s models.Seller) models.Seller {
	lastId := r.getLastId()
	s.Id = lastId + 1

	r.db[s.Id] = s

	return s
}

func (r *sellerRepository) Update(id int, s models.Seller) {
	r.db[id] = s
}

func (r *sellerRepository) Delete(id int) {
	delete(r.db, id)
}

func (r *sellerRepository) FindAll() []models.Seller {
	sellers := make([]models.Seller, 0, len(r.db))
	for _, s := range r.db {
		sellers = append(sellers, s)
	}
	return sellers
}

func (r *sellerRepository) FindById(id int) (*models.Seller, *api.ServiceError) {
	s, exists := r.db[id]
	if !exists {
		err := api.ServiceErrors[api.ErrNotFound]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The seller you are looking for does not exist.",
		}
	}

	return &s, nil
}

func (r *sellerRepository) IsValidCid(cid int) bool {
	for _, s := range r.db {
		if s.Cid == cid {
			return false
		}
	}

	return true
}

func (r *sellerRepository) IsValidCidExcludeId(cid int, excludeId int) bool {
	for _, s := range r.db {
		if s.Cid == cid && s.Id != excludeId {
			return false
		}
	}
	return true
}

func (r *sellerRepository) ExistsById(id int) bool {
	_, exists := r.db[id]

	return exists
}

func (r *sellerRepository) getLastId() int {
	maxId := 0
	for id := range r.db {
		if id > maxId {
			maxId = id
		}
	}
	return maxId
}
