package service

import (
	"slices"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models"
)

type BuyerService interface {
	Create(req models.RequestBuyer) (*models.ResponseBuyer, *api.ServiceError)
	Update(id int, req models.RequestBuyer) (*models.ResponseBuyer, *api.ServiceError)
	Delete(id int) *api.ServiceError
	FindAll() []models.ResponseBuyer
	FindById(id int) (*models.ResponseBuyer, *api.ServiceError)
}

type buyerService struct {
	rp repository.BuyerRepository
}

func NewBuyerService(rp repository.BuyerRepository) BuyerService {
	return &buyerService{rp: rp}
}

func (sv *buyerService) Create(req models.RequestBuyer) (*models.ResponseBuyer, *api.ServiceError) {
	if !sv.rp.IsValidCardNumberId(*req.CardNumberId) {
		err := api.ServiceErrors[api.ErrConflict]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The card number ID is already in use, please use a different one",
		}
	}

	mb := mappers.RequestBuyerToBuyer(req)
	b := sv.rp.Create(mb)
	resp := models.ResponseBuyer(b)
	return &resp, nil
}

func (sv *buyerService) Update(id int, req models.RequestBuyer) (*models.ResponseBuyer, *api.ServiceError) {
	existing, err := sv.rp.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.CardNumberId != nil {
		if !sv.rp.IsValidCardNumberIdExcludeId(*req.CardNumberId, id) {
			err := api.ServiceErrors[api.ErrConflict]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "The card number ID is already in use by another buyer",
			}
		}
	}

	updatedBuyer := mappers.RequestBuyerToBuyerUpadate(req, *existing)
	sv.rp.Update(id, updatedBuyer)

	resp := models.ResponseBuyer(updatedBuyer)
	return &resp, nil
}

func (sv *buyerService) Delete(id int) *api.ServiceError {
	if !sv.rp.ExistsById(id) {
		err := api.ServiceErrors[api.ErrNotFound]
		return &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The buyer you are trying to delete does not exist",
		}
	}
	sv.rp.Delete(id)
	return nil
}

func (sv *buyerService) FindAll() []models.ResponseBuyer {
	var rb []models.ResponseBuyer
	for _, b := range sv.rp.FindAll() {
		rb = append(rb, models.ResponseBuyer(b))
	}
	slices.SortFunc(rb, func(a, b models.ResponseBuyer) int {
		return a.Id - b.Id
	})
	return rb
}

func (sv *buyerService) FindById(id int) (*models.ResponseBuyer, *api.ServiceError) {
	b, err := sv.rp.FindById(id)
	if err != nil {
		return nil, err
	}
	resp := models.ResponseBuyer(*b)
	return &resp, nil
}
