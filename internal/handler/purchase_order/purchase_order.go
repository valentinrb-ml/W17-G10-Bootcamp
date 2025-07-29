package handler

import (
	"errors"
	"net/http"

	purchaseOrderService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type PurchaseOrderHandler struct {
	service purchaseOrderService.PurchaseOrderService
}

func NewPurchaseOrderHandler(s purchaseOrderService.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{service: s}
}

func (h *PurchaseOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var wrapper models.PurchaseOrderRequestWrapper
	if err := httputil.DecodeJSON(r, &wrapper); err != nil {
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"))
		return
	}

	req := wrapper.Data

	if err := validators.ValidatePurchaseOrderPost(req); err != nil {
		response.Error(w, err)
		return
	}

	createdPO, err := h.service.Create(ctx, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, createdPO)
}

func (h *PurchaseOrderHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	buyerID, err := httputil.ParseIntQueryParam(r, "id")
	if err != nil && !errors.Is(err, httputil.ErrParamNotProvided) {
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid buyer ID parameter"))
		return
	}

	var report []models.BuyerWithPurchaseCount
	if buyerID != nil {
		report, err = h.service.GetReportByBuyer(ctx, buyerID)
	} else {
		report, err = h.service.GetReportByBuyer(ctx, nil)
	}

	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, report)
}
