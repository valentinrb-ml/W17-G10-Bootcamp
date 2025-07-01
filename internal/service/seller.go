package service

import (
	"slices"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

type SellerService interface {
	Create(reqs models.RequestSeller) (*models.ResponseSeller, *api.ServiceError)
	Update(id int, reqs models.RequestSeller) (*models.ResponseSeller, *api.ServiceError)
	Delete(id int) *api.ServiceError
	FindAll() []models.ResponseSeller
	FindById(id int) (*models.ResponseSeller, *api.ServiceError)
}

type sellerService struct {
	rp repository.SellerRepository
}

func NewSellerService(rp repository.SellerRepository) SellerService {
	return &sellerService{rp: rp}
}

func (sv *sellerService) Create(reqs models.RequestSeller) (*models.ResponseSeller, *api.ServiceError) {
	if !sv.rp.IsValidCid(*reqs.Cid) {
		err := api.ServiceErrors[api.ErrConflict]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The CID is in use, use a valid one",
		}
	}

	ms := mappers.RequestSellerToSeller(reqs)

	s := sv.rp.Create(ms)
	resps := models.ResponseSeller(s)

	return &resps, nil
}

func (sv *sellerService) Update(id int, reqs models.RequestSeller) (*models.ResponseSeller, *api.ServiceError) {
	if !sv.rp.ExistsById(id) {
		err := api.ServiceErrors[api.ErrNotFound]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The seller you are trying to update does not exist",
		}
	}
	if !sv.rp.IsValidCidExcludeId(*reqs.Cid, id) {
		err := api.ServiceErrors[api.ErrConflict]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The CID is in use, use a valid one",
		}
	}

	ms := mappers.RequestSellerToSeller(reqs)

	ms.Id = id
	sv.rp.Update(id, ms)
	resps := models.ResponseSeller(ms)

	return &resps, nil
}

func (sv *sellerService) Delete(id int) *api.ServiceError {
	if !sv.rp.ExistsById(id) {
		err := api.ServiceErrors[api.ErrNotFound]
		return &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The seller you are trying to delete does not exist",
		}
	}

	sv.rp.Delete(id)

	return nil
}

func (sv *sellerService) FindAll() []models.ResponseSeller {
	var rs []models.ResponseSeller
	for _, s := range sv.rp.FindAll() {
		rs = append(rs, models.ResponseSeller(s))
	}

	slices.SortFunc(rs, func(a, b models.ResponseSeller) int {
		return a.Id - b.Id
	})

	return rs
}

func (sv *sellerService) FindById(id int) (*models.ResponseSeller, *api.ServiceError) {
	s, err := sv.rp.FindById(id)
	if err != nil {
		return nil, err
	}

	resps := models.ResponseSeller(*s)

	return &resps, nil
}
