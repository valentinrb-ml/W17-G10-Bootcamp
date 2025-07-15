package service

import (
	"context"
	"slices"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func (sv *sellerService) Create(ctx context.Context, reqs models.RequestSeller) (*models.ResponseSeller, error) {
	ms := mappers.RequestSellerToSeller(reqs)

	s, err := sv.rp.Create(ctx, ms)
	if err != nil {
		return nil, err
	}

	resps := mappers.ToResponseSeller(s)

	return &resps, nil
}

func (sv *sellerService) Update(ctx context.Context, id int, reqs models.RequestSeller) (*models.ResponseSeller, error) {
	s, err := sv.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	mappers.ApplySellerPatch(s, &reqs)

	err = sv.rp.Update(ctx, id, *s)
	if err != nil {
		return nil, err
	}

	resps := mappers.ToResponseSeller(s)

	return &resps, nil
}

func (sv *sellerService) Delete(ctx context.Context, id int) error {
	err := sv.rp.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (sv *sellerService) FindAll(ctx context.Context) ([]models.ResponseSeller, error) {
	s, err := sv.rp.FindAll(ctx)
	if err != nil {
		return []models.ResponseSeller{}, err
	}

	rs := mappers.ToResponseSellerList(s)

	slices.SortFunc(rs, func(a, b models.ResponseSeller) int {
		return a.Id - b.Id
	})

	return rs, nil
}

func (sv *sellerService) FindById(ctx context.Context, id int) (*models.ResponseSeller, error) {
	s, err := sv.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	resps := mappers.ToResponseSeller(s)

	return &resps, nil
}