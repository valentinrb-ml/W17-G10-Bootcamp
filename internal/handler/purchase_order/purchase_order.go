package handler

import (
	"errors"
	"net/http"

	purchaseOrderService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type PurchaseOrderHandler struct {
	service purchaseOrderService.PurchaseOrderService
	logger  logger.Logger
}

func NewPurchaseOrderHandler(s purchaseOrderService.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{service: s}
}

// SetLogger allows injecting the logger after creation
func (h *PurchaseOrderHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

func (h *PurchaseOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(ctx, "purchase-order-handler", "Create purchase order request received")
	}

	var wrapper models.PurchaseOrderRequestWrapper
	if err := httputil.DecodeJSON(r, &wrapper); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "purchase-order-handler", "Invalid JSON in create request", map[string]interface{}{
				"error": err.Error(),
			})
		}
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"))
		return
	}

	req := wrapper.Data

	if err := validators.ValidatePurchaseOrderPost(req); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "purchase-order-handler", "Validation failed for create request", map[string]interface{}{
				"order_number":      req.OrderNumber,
				"buyer_id":          req.BuyerID,
				"product_record_id": req.ProductRecordID,
				"validation_error":  err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	createdPO, err := h.service.Create(ctx, req)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "purchase-order-handler", "Failed to create purchase order", err, map[string]interface{}{
				"order_number":      req.OrderNumber,
				"buyer_id":          req.BuyerID,
				"product_record_id": req.ProductRecordID,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "purchase-order-handler", "Purchase order created successfully", map[string]interface{}{
			"purchase_order_id": createdPO.ID,
			"order_number":      createdPO.OrderNumber,
			"status_code":       http.StatusCreated,
		})
	}

	response.JSON(w, http.StatusCreated, createdPO)
}

func (h *PurchaseOrderHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(ctx, "purchase-order-handler", "Get purchase report request received")
	}

	buyerID, err := httputil.ParseIntQueryParam(r, "id")
	if err != nil && !errors.Is(err, httputil.ErrParamNotProvided) {
		if h.logger != nil {
			h.logger.Warning(ctx, "purchase-order-handler", "Invalid buyer ID parameter", map[string]interface{}{
				"id_param": r.URL.Query().Get("id"),
				"error":    err.Error(),
			})
		}
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid buyer ID parameter"))
		return
	}

	var report []models.BuyerWithPurchaseCount
	if buyerID != nil {
		if h.logger != nil {
			h.logger.Info(ctx, "purchase-order-handler", "Getting purchase report for specific buyer", map[string]interface{}{
				"buyer_id": *buyerID,
			})
		}
		report, err = h.service.GetReportByBuyer(ctx, buyerID)
	} else {
		if h.logger != nil {
			h.logger.Info(ctx, "purchase-order-handler", "Getting purchase report for all buyers")
		}
		report, err = h.service.GetReportByBuyer(ctx, nil)
	}

	if err != nil {
		if h.logger != nil {
			if buyerID != nil {
				h.logger.Error(ctx, "purchase-order-handler", "Failed to get purchase report for buyer", err, map[string]interface{}{
					"buyer_id": *buyerID,
				})
			} else {
				h.logger.Error(ctx, "purchase-order-handler", "Failed to get purchase report for all buyers", err, nil)
			}
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		if buyerID != nil {
			h.logger.Info(ctx, "purchase-order-handler", "Purchase report retrieved successfully", map[string]interface{}{
				"buyer_id":    *buyerID,
				"count":       len(report),
				"status_code": http.StatusOK,
			})
		} else {
			h.logger.Info(ctx, "purchase-order-handler", "Purchase report for all buyers retrieved successfully", map[string]interface{}{
				"buyers_count": len(report),
				"status_code":  http.StatusOK,
			})
		}
	}

	response.JSON(w, http.StatusOK, report)
}
