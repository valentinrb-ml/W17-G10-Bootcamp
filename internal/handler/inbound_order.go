package handler

import (
	"encoding/json"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
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
