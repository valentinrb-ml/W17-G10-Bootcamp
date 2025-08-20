package service

import (
	"context"
	"slices"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func (sv *sellerService) Create(ctx context.Context, reqs models.RequestSeller) (*models.ResponseSeller, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Creating seller", map[string]interface{}{
			"cid":          reqs.Cid,
			"company_name": reqs.CompanyName,
		})
	}

	ms := mappers.RequestSellerToSeller(reqs)

	s, err := sv.sellerRepo.Create(ctx, ms)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "seller-service", "Failed to create seller", err, map[string]interface{}{
				"cid":          reqs.Cid,
				"company_name": reqs.CompanyName,
			})
		}
		return nil, err
	}

	resps := mappers.ToResponseSeller(s)

	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Seller created successfully", map[string]interface{}{
			"seller_id":    resps.Id,
			"cid":          resps.Cid,
			"company_name": resps.CompanyName,
		})
	}

	return &resps, nil
}

func (sv *sellerService) Update(ctx context.Context, id int, reqs models.RequestSeller) (*models.ResponseSeller, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Updating seller", map[string]interface{}{
			"seller_id": id,
		})
	}

	s, err := sv.sellerRepo.FindById(ctx, id)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "seller-service", "Failed to find seller for update", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return nil, err
	}

	mappers.ApplySellerPatch(s, &reqs)

	err = sv.sellerRepo.Update(ctx, id, *s)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "seller-service", "Failed to update seller", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return nil, err
	}

	resps := mappers.ToResponseSeller(s)

	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Seller updated successfully", map[string]interface{}{
			"seller_id":    resps.Id,
			"cid":          resps.Cid,
			"company_name": resps.CompanyName,
		})
	}

	return &resps, nil
}

func (sv *sellerService) Delete(ctx context.Context, id int) error {
	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Deleting seller", map[string]interface{}{
			"seller_id": id,
		})
	}

	err := sv.sellerRepo.Delete(ctx, id)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "seller-service", "Failed to delete seller", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return err
	}

	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Seller deleted successfully", map[string]interface{}{
			"seller_id": id,
		})
	}

	return nil
}

func (sv *sellerService) FindAll(ctx context.Context) ([]models.ResponseSeller, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Finding all sellers")
	}

	s, err := sv.sellerRepo.FindAll(ctx)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "seller-service", "Failed to find all sellers", err, nil)
		}
		return []models.ResponseSeller{}, err
	}

	rs := mappers.ToResponseSellerList(s)

	slices.SortFunc(rs, func(a, b models.ResponseSeller) int {
		return a.Id - b.Id
	})

	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Found all sellers successfully", map[string]interface{}{
			"sellers_count": len(rs),
		})
	}

	return rs, nil
}

func (sv *sellerService) FindById(ctx context.Context, id int) (*models.ResponseSeller, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Finding seller by ID", map[string]interface{}{
			"seller_id": id,
		})
	}

	s, err := sv.sellerRepo.FindById(ctx, id)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "seller-service", "Failed to find seller by ID", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return nil, err
	}

	resps := mappers.ToResponseSeller(s)

	if sv.logger != nil {
		sv.logger.Info(ctx, "seller-service", "Seller found successfully", map[string]interface{}{
			"seller_id":    resps.Id,
			"cid":          resps.Cid,
			"company_name": resps.CompanyName,
		})
	}

	return &resps, nil
}
