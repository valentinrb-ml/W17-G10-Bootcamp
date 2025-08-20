package handler

import (
	"net/http"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type BuyerHandler struct {
	sv     service.BuyerService
	logger logger.Logger
}

func NewBuyerHandler(sv service.BuyerService) *BuyerHandler {
	return &BuyerHandler{sv: sv}
}

// SetLogger allows injecting the logger after creation
func (h *BuyerHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Create buyer request received")
	}

	var br models.RequestBuyer
	if err := httputil.DecodeJSON(r, &br); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Invalid JSON in create request", map[string]interface{}{
				"error": err.Error(),
			})
		}
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"))
		return
	}

	if err := validators.ValidateRequestBuyer(br); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Validation failed for create request", map[string]interface{}{
				"card_number_id":   br.CardNumberId,
				"first_name":       br.FirstName,
				"last_name":        br.LastName,
				"validation_error": err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	b, err := h.sv.Create(ctx, br)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "buyer-handler", "Failed to create buyer", err, map[string]interface{}{
				"card_number_id": br.CardNumberId,
				"first_name":     br.FirstName,
				"last_name":      br.LastName,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Buyer created successfully", map[string]interface{}{
			"buyer_id":       b.Id,
			"card_number_id": b.CardNumberId,
			"first_name":     b.FirstName,
			"last_name":      b.LastName,
		})
	}

	response.JSON(w, http.StatusCreated, b)
}

func (h *BuyerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Invalid buyer ID in update request", map[string]interface{}{
				"id_param": r.URL.Query().Get("id"),
			})
		}
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid ID parameter"))
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Update buyer request received", map[string]interface{}{
			"buyer_id": id,
		})
	}

	var br models.RequestBuyer
	if err := httputil.DecodeJSON(r, &br); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Invalid JSON in update request", map[string]interface{}{
				"buyer_id": id,
				"error":    err.Error(),
			})
		}
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"))
		return
	}

	if err := validators.ValidateUpdateBuyer(br); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Validation failed for update request", map[string]interface{}{
				"buyer_id":         id,
				"validation_error": err.Error(),
			})
		}
		response.Error(w, err)
		return
	}
	if err := validators.ValidateBuyerPatchNotEmpty(br); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Empty validation failed for update request", map[string]interface{}{
				"buyer_id":         id,
				"validation_error": err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	updated, err := h.sv.Update(ctx, id, br)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "buyer-handler", "Failed to update buyer", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Buyer updated successfully", map[string]interface{}{
			"buyer_id":       updated.Id,
			"card_number_id": updated.CardNumberId,
			"first_name":     updated.FirstName,
			"last_name":      updated.LastName,
		})
	}

	response.JSON(w, http.StatusOK, updated)
}

func (h *BuyerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Invalid buyer ID in delete request", map[string]interface{}{
				"id_param": r.URL.Query().Get("id"),
			})
		}
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid ID parameter"))
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Delete buyer request received", map[string]interface{}{
			"buyer_id": id,
		})
	}

	if err := h.sv.Delete(ctx, id); err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "buyer-handler", "Failed to delete buyer", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Buyer deleted successfully", map[string]interface{}{
			"buyer_id": id,
		})
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *BuyerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Find all buyers request received")
	}

	result, err := h.sv.FindAll(ctx)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "buyer-handler", "Failed to find all buyers", err, nil)
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Find all buyers completed successfully", map[string]interface{}{
			"buyers_count": len(result),
		})
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *BuyerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "buyer-handler", "Invalid buyer ID in find by ID request", map[string]interface{}{
				"id_param": r.URL.Query().Get("id"),
			})
		}
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid ID parameter"))
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Find buyer by ID request received", map[string]interface{}{
			"buyer_id": id,
		})
	}

	b, err := h.sv.FindById(ctx, id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "buyer-handler", "Failed to find buyer by ID", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "buyer-handler", "Buyer found successfully", map[string]interface{}{
			"buyer_id":       b.Id,
			"card_number_id": b.CardNumberId,
			"first_name":     b.FirstName,
			"last_name":      b.LastName,
		})
	}

	response.JSON(w, http.StatusOK, b)
}
