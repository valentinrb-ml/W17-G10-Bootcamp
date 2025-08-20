package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

type InboundOrderHandler struct {
	service service.InboundOrderService
	logger  logger.Logger
}

func NewInboundOrderHandler(s service.InboundOrderService) *InboundOrderHandler {
	return &InboundOrderHandler{service: s}
}
func (h *InboundOrderHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

// POST /api/v1/inboundOrders
func (h *InboundOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "inboundorder-handler", "Create inbound order request received")
	}
	var payload struct {
		Data models.InboundOrder `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "inboundorder-handler", "Invalid JSON in create request", map[string]interface{}{
				"error": err.Error(),
			})
		}
		response.ErrorWithRequest(w, r, apperrors.NewAppError(apperrors.CodeValidationError, "invalid JSON format"))
		return
	}
	created, err := h.service.Create(r.Context(), &payload.Data)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "inboundorder-handler", "Failed to create inbound order", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if h.logger != nil {
		h.logger.Info(r.Context(), "inboundorder-handler", "Inbound order created successfully", map[string]interface{}{
			"inbound_order_id": created.ID,
		})
	}
	response.JSON(w, http.StatusCreated, created)
}

// GET /api/v1/employees/reportInboundOrders?id=1
func (h *InboundOrderHandler) Report(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "inboundorder-handler", "Report inbound orders request received")
	}
	var idPtr *int
	if idStr := r.URL.Query().Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			if h.logger != nil {
				h.logger.Warning(r.Context(), "inboundorder-handler", "Invalid employee id in report request", map[string]interface{}{
					"id_param": idStr,
				})
			}
			response.ErrorWithRequest(w, r, apperrors.NewAppError(apperrors.CodeBadRequest, "id must be int"))
			return
		}
		idPtr = &id
	}
	report, err := h.service.Report(r.Context(), idPtr)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "inboundorder-handler", "Failed to generate inbound order report", err)
		}

		response.ErrorWithRequest(w, r, err)
		return
	}
	if h.logger != nil {
		h.logger.Info(r.Context(), "inboundorder-handler", "Inbound order report generated successfully")
	}
	response.JSON(w, http.StatusOK, report)
}
