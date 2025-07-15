package handler

import (
	"errors"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
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

	var req models.RequestPurchaseOrder
	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	createdPO, err := h.service.Create(ctx, req)
	if err != nil {
		response.Error(w, convertServiceError(err))
		return
	}

	response.JSON(w, http.StatusCreated, createdPO)
}

func (h *PurchaseOrderHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener parámetro de query 'id' si existe
	buyerID, err := httputil.ParseIntQueryParam(r, "id")
	if err != nil && !errors.Is(err, httputil.ErrParamNotProvided) {
		response.Error(w, convertServiceError(err))
		return
	}

	var report interface{}
	if buyerID != nil {
		// Reporte para un buyer específico
		report, err = h.service.GetReportByBuyer(ctx, buyerID)
	} else {
		// Reporte general
		report, err = h.service.GetReportByBuyer(ctx, nil)
	}

	if err != nil {
		response.Error(w, convertServiceError(err))
		return
	}

	response.JSON(w, http.StatusOK, report)
}

// convertServiceError asegura que todos los errores sean del tipo esperado por el handler
func convertServiceError(err error) error {
	if err == nil {
		return nil
	}

	// Si ya es un ServiceError, lo retornamos directamente
	if svcErr, ok := err.(*api.ServiceError); ok {
		return svcErr
	}

	// Para otros tipos de errores, los convertimos a ServiceError
	return &api.ServiceError{
		Code:         http.StatusInternalServerError,
		ResponseCode: http.StatusInternalServerError,
		Message:      err.Error(),
	}
}
