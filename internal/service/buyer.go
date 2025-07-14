package service

import (
	"context"
	"slices"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func (sv *buyerService) Create(ctx context.Context, req models.RequestBuyer) (*models.ResponseBuyer, error) {
	if sv.rp.CardNumberExists(ctx, *req.CardNumberId, 0) {
		err := api.ServiceErrors[api.ErrConflict]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "The card number ID is already in use, please use a different one",
		}
	}

	mb := mappers.RequestBuyerToBuyer(req)

	b, err := sv.rp.Create(ctx, mb)
	if err != nil {
		return nil, err
	}

	resp := mappers.ToResponseBuyer(b)
	return &resp, nil
}

func (sv *buyerService) Update(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error) {
	if req.CardNumberId != nil {
		if sv.rp.CardNumberExists(ctx, *req.CardNumberId, id) {
			err := api.ServiceErrors[api.ErrConflict]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "The card number ID is already in use by another buyer",
			}
		}
	}
	b, err := sv.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	mappers.ApplyBuyerPatch(b, &req)

	if err := sv.rp.Update(ctx, id, *b); err != nil {
		return nil, err
	}
	resp := mappers.ToResponseBuyer(b)
	return &resp, nil
}

func (sv *buyerService) Delete(ctx context.Context, id int) error {
	err := sv.rp.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (sv *buyerService) FindAll(ctx context.Context) ([]models.ResponseBuyer, error) {
	bs, err := sv.rp.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	rs := mappers.ToResponseBuyerList(bs)
	slices.SortFunc(rs, func(a, b models.ResponseBuyer) int {
		return a.Id - b.Id
	})

	return rs, nil
}

func (sv *buyerService) FindById(ctx context.Context, id int) (*models.ResponseBuyer, error) {
	b, err := sv.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := mappers.ToResponseBuyer(b)
	return &resp, nil
}
