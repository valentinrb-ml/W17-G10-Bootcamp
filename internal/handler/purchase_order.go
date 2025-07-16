package handler

import (
	"errors"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type PurchaseOrderHandler struct {
	service service.PurchaseOrderService
}

func NewPurchaseOrderHandler(s service.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{service: s}
}
func (h *PurchaseOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var wrapper models.PurchaseOrderRequestWrapper
	if err := httputil.DecodeJSON(r, &wrapper); err != nil {
		response.Error(w, err)
		return
	}

	req := wrapper.Data // Extraemos los datos del wrapper

	if err := validators.ValidatePurchaseOrderPost(req); err != nil {
		response.Error(w, httputil.ConvertServiceErrorToAppError(err))
		return
	}

	createdPO, err := h.service.Create(ctx, req)
	if err != nil {
		response.Error(w, httputil.ConvertServiceErrorToAppError(err))
		return
	}

	// Envuelve la respuesta también si es necesario
	response.JSON(w, http.StatusCreated, createdPO)
}

func (h *PurchaseOrderHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	buyerID, err := httputil.ParseIntQueryParam(r, "id")
	if err != nil && !errors.Is(err, httputil.ErrParamNotProvided) {
		//response.Error(w, convertServiceError(err))
		response.Error(w, httputil.ConvertServiceErrorToAppError(err))
		return
	}

	var report []models.BuyerWithPurchaseCount
	if buyerID != nil {
		report, err = h.service.GetReportByBuyer(ctx, buyerID)
	} else {
		report, err = h.service.GetReportByBuyer(ctx, nil)
	}

	if err != nil {
		//response.Error(w, convertServiceError(err))
		response.Error(w, httputil.ConvertServiceErrorToAppError(err))
		return
	}

	response.JSON(w, http.StatusOK, report)
}
