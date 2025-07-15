package handler

import (
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

func NewCarryHandler(sv service.CarryService) *CarryHandler {
	return &CarryHandler{sv: sv}
}

type CarryHandler struct {
	// sv is the service that will be used by the handler
	sv service.CarryService
}

func (h *CarryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req = carry.CarryRequest{}
	if err := request.JSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validators.ValidateCarryCreateRequest(req); err != nil {
		response.Error(w, err)
		return
	}

	wh := mappers.RequestToCarry(req)

	newC, err := h.sv.Create(r.Context(), wh)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, mappers.CarryToDoc(newC))
}

func (h *CarryHandler) ReportCarries(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	var localityID *string = nil
	if idParam != "" {
		localityID = &idParam
	}

	result, err := h.sv.GetCarriesReport(r.Context(), localityID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
