package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

type InboundOrderHandler struct {
	service service.InboundOrderService
}

func NewInboundOrderHandler(s service.InboundOrderService) *InboundOrderHandler {
	return &InboundOrderHandler{service: s}
}

// POST /api/v1/inboundOrders
func (h *InboundOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Data models.InboundOrder `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, apperrors.NewAppError(apperrors.CodeValidationError, "invalid JSON format"))
		return
	}
	created, err := h.service.Create(r.Context(), &payload.Data)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, created)
}

// GET /api/v1/employees/reportInboundOrders?id=1
func (h *InboundOrderHandler) Report(w http.ResponseWriter, r *http.Request) {
	var idPtr *int
	if idStr := r.URL.Query().Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "id must be int"))
			return
		}
		idPtr = &id
	}
	report, err := h.service.Report(r.Context(), idPtr)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, report)
}
