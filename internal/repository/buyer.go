package repository

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models"
)

type BuyerRepository interface {
	Create(b models.Buyer) models.Buyer
	Update(id int, b models.Buyer)
	Delete(id int)
	FindAll() []models.Buyer
	FindById(id int) (*models.Buyer, *api.ServiceError)
	IsValidCardNumberId(cardNumberId string) bool
	IsValidCardNumberIdExcludeId(cardNumberId string, excludeId int) bool
	ExistsById(id int) bool
}

type buyerRepository struct {
	db map[int]models.Buyer
}

func NewBuyerRepository(db map[int]models.Buyer) BuyerRepository {
	defaultDb := make(map[int]models.Buyer)
	if db != nil {
		defaultDb = db
	}
	return &buyerRepository{db: defaultDb}
}

func (r *buyerRepository) Create(b models.Buyer) models.Buyer {
	lastId := r.getLastId()
	b.Id = lastId + 1
	r.db[b.Id] = b
	return b
}

func (r *buyerRepository) Update(id int, b models.Buyer) {
	r.db[id] = b
}

func (r *buyerRepository) Delete(id int) {
	delete(r.db, id)
}

func (r *buyerRepository) FindAll() []models.Buyer {
	buyers := make([]models.Buyer, 0, len(r.db))
	for _, b := range r.db {
		buyers = append(buyers, b)
	}
	return buyers
}

func (r *buyerRepository) FindById(id int) (*models.Buyer, *api.ServiceError) {
	b, exists := r.db[id]
	if !exists {
		err := api.ServiceErrors[api.ErrNotFound]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The buyer you are looking for does not exist.",
		}
	}
	return &b, nil
}

func (r *buyerRepository) IsValidCardNumberId(cardNumberId string) bool {
	for _, b := range r.db {
		if b.CardNumberId == cardNumberId {
			return false
		}
	}
	return true
}

func (r *buyerRepository) IsValidCardNumberIdExcludeId(cardNumberId string, excludeId int) bool {
	for _, b := range r.db {
		if b.CardNumberId == cardNumberId && b.Id != excludeId {
			return false
		}
	}
	return true
}

func (r *buyerRepository) ExistsById(id int) bool {
	_, exists := r.db[id]
	return exists
}

func (r *buyerRepository) getLastId() int {
	maxId := 0
	for id := range r.db {
		if id > maxId {
			maxId = id
		}
	}
	return maxId
}
