package handler

import (
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

func NewCarryHandler(sv service.CarryService) *CarryHandler {
	return &CarryHandler{sv: sv}
}

type CarryHandler struct {
	// sv is the service that will be used by the handler
	sv     service.CarryService
	logger logger.Logger
}

// SetLogger allows injecting the logger after creation
func (h *CarryHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

func (h *CarryHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Context(), "carry", "Starting carry creation request")

	var req = carry.CarryRequest{}
	if err := request.JSON(r, &req); err != nil {
		h.logger.Error(r.Context(), "carry", "Failed to parse JSON request", err)
		response.ErrorWithRequest(w, r, err)
		return
	}

	h.logger.Debug(r.Context(), "carry", "Carry request parsed successfully", map[string]interface{}{
		"company_name": req.CompanyName,
		"cid":          req.Cid,
	})

	if err := validators.ValidateCarryCreateRequest(req); err != nil {
		h.logger.Warning(r.Context(), "carry", "Carry validation failed", map[string]interface{}{
			"error":        err.Error(),
			"company_name": req.CompanyName,
		})
		response.ErrorWithRequest(w, r, err)
		return
	}

	wh := mappers.RequestToCarry(req)

	newC, err := h.sv.Create(r.Context(), wh)
	if err != nil {
		h.logger.Error(r.Context(), "carry", "Failed to create carry", err, map[string]interface{}{
			"company_name": req.CompanyName,
			"cid":          req.Cid,
		})
		response.ErrorWithRequest(w, r, err)
		return
	}

	h.logger.Info(r.Context(), "carry", "Carry created successfully", map[string]interface{}{
		"carry_id":     newC.Id,
		"company_name": newC.CompanyName,
		"cid":          newC.Cid,
	})
	response.JSON(w, http.StatusCreated, mappers.CarryToDoc(newC))
}

func (h *CarryHandler) ReportCarries(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Context(), "carry", "Starting carries report request")

	idParam := r.URL.Query().Get("id")
	var localityID *string = nil
	if idParam != "" {
		localityID = &idParam
		h.logger.Debug(r.Context(), "carry", "Filtering carries report by locality", map[string]interface{}{
			"locality_id": idParam,
		})
	}

	result, err := h.sv.GetCarriesReport(r.Context(), localityID)
	if err != nil {
		h.logger.Error(r.Context(), "carry", "Failed to get carries report", err, map[string]interface{}{
			"locality_id": idParam,
		})
		response.ErrorWithRequest(w, r, err)
		return
	}

	h.logger.Info(r.Context(), "carry", "Carries report generated successfully", map[string]interface{}{
		"locality_id": idParam,
	})
	response.JSON(w, http.StatusOK, result)
}
