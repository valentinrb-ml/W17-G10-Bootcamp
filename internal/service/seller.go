package service

import (
	"slices"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func (sv *sellerService) Create(reqs models.RequestSeller) (*models.ResponseSeller, error) {
	if sv.rp.CIDExists(*reqs.Cid, 0) {
		err := api.ServiceErrors[api.ErrConflict]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The CID is in use, use a valid one",
		}
	}

	ms := mappers.RequestSellerToSeller(reqs)

	s, err := sv.rp.Create(ms)
	if err != nil {
		return nil, err
	}

	resps := models.ResponseSeller(*s)

	return &resps, nil
}

func (sv *sellerService) Update(id int, reqs models.RequestSeller) (*models.ResponseSeller, error) {
	if sv.rp.CIDExists(*reqs.Cid, id) {
		err := api.ServiceErrors[api.ErrConflict]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The CID is in use, use a valid one",
		}
	}
	s, err := sv.rp.FindById(id)
	if err != nil {
		return nil, err
	}

	mappers.ApplySellerPatch(s, &reqs)

	sv.rp.Update(id, *s)
	resps := models.ResponseSeller(*s)

	return &resps, nil
}

func (sv *sellerService) Delete(id int) error {
	err := sv.rp.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (sv *sellerService) FindAll() ([]models.ResponseSeller, error) {
	var rs []models.ResponseSeller
	s, err := sv.rp.FindAll()
	if err != nil {
		return []models.ResponseSeller{}, err
	}

	for _, s := range s {
		rs = append(rs, models.ResponseSeller(s))
	}

	slices.SortFunc(rs, func(a, b models.ResponseSeller) int {
		return a.Id - b.Id
	})

	return rs, nil
}

func (sv *sellerService) FindById(id int) (*models.ResponseSeller, error) {
	s, err := sv.rp.FindById(id)
	if err != nil {
		return nil, err
	}

	resps := models.ResponseSeller(*s)

	return &resps, nil
}
