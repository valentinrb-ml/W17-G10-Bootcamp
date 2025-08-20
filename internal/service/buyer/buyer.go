package service

import (
	"context"
	"slices"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func (sv *buyerService) Create(ctx context.Context, req models.RequestBuyer) (*models.ResponseBuyer, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Creating buyer", map[string]interface{}{
			"card_number_id": req.CardNumberId,
			"first_name":     req.FirstName,
			"last_name":      req.LastName,
		})
	}

	if sv.rp.CardNumberExists(ctx, *req.CardNumberId, 0) {
		if sv.logger != nil {
			sv.logger.Warning(ctx, "buyer-service", "Card number already exists", map[string]interface{}{
				"card_number_id": req.CardNumberId,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "The card number ID is already in use, please use a different one")
	}

	mb := mappers.RequestBuyerToBuyer(req)

	b, err := sv.rp.Create(ctx, mb)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "buyer-service", "Failed to create buyer", err, map[string]interface{}{
				"card_number_id": req.CardNumberId,
				"first_name":     req.FirstName,
				"last_name":      req.LastName,
			})
		}
		return nil, err
	}

	resp := mappers.ToResponseBuyer(b)

	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Buyer created successfully", map[string]interface{}{
			"buyer_id":       resp.Id,
			"card_number_id": resp.CardNumberId,
			"first_name":     resp.FirstName,
			"last_name":      resp.LastName,
		})
	}

	return &resp, nil
}

func (sv *buyerService) Update(ctx context.Context, id int, req models.RequestBuyer) (*models.ResponseBuyer, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Updating buyer", map[string]interface{}{
			"buyer_id": id,
		})
	}

	if req.CardNumberId != nil {
		if sv.rp.CardNumberExists(ctx, *req.CardNumberId, id) {
			if sv.logger != nil {
				sv.logger.Warning(ctx, "buyer-service", "Card number already exists for another buyer", map[string]interface{}{
					"buyer_id":       id,
					"card_number_id": req.CardNumberId,
				})
			}
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "The card number ID is already in use by another buyer")
		}
	}

	b, err := sv.rp.FindById(ctx, id)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "buyer-service", "Failed to find buyer for update", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return nil, err
	}

	mappers.ApplyBuyerPatch(b, &req)

	if err := sv.rp.Update(ctx, id, *b); err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "buyer-service", "Failed to update buyer", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return nil, err
	}

	resp := mappers.ToResponseBuyer(b)

	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Buyer updated successfully", map[string]interface{}{
			"buyer_id":       resp.Id,
			"card_number_id": resp.CardNumberId,
			"first_name":     resp.FirstName,
			"last_name":      resp.LastName,
		})
	}

	return &resp, nil
}

func (sv *buyerService) Delete(ctx context.Context, id int) error {
	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Deleting buyer", map[string]interface{}{
			"buyer_id": id,
		})
	}

	err := sv.rp.Delete(ctx, id)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "buyer-service", "Failed to delete buyer", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return err
	}

	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Buyer deleted successfully", map[string]interface{}{
			"buyer_id": id,
		})
	}

	return nil
}

func (sv *buyerService) FindAll(ctx context.Context) ([]models.ResponseBuyer, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Finding all buyers")
	}

	bs, err := sv.rp.FindAll(ctx)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "buyer-service", "Failed to find all buyers", err, nil)
		}
		return nil, err
	}

	rs := mappers.ToResponseBuyerList(bs)
	slices.SortFunc(rs, func(a, b models.ResponseBuyer) int {
		return a.Id - b.Id
	})

	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Found all buyers successfully", map[string]interface{}{
			"buyers_count": len(rs),
		})
	}

	return rs, nil
}

func (sv *buyerService) FindById(ctx context.Context, id int) (*models.ResponseBuyer, error) {
	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Finding buyer by ID", map[string]interface{}{
			"buyer_id": id,
		})
	}

	b, err := sv.rp.FindById(ctx, id)
	if err != nil {
		if sv.logger != nil {
			sv.logger.Error(ctx, "buyer-service", "Failed to find buyer by ID", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return nil, err
	}

	resp := mappers.ToResponseBuyer(b)

	if sv.logger != nil {
		sv.logger.Info(ctx, "buyer-service", "Buyer found successfully", map[string]interface{}{
			"buyer_id":       resp.Id,
			"card_number_id": resp.CardNumberId,
			"first_name":     resp.FirstName,
			"last_name":      resp.LastName,
		})
	}

	return &resp, nil
}
