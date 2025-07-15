package handler

import (
	"errors"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func NewGeographyHandler(sv service.GeographyService) *GeographyHandler {
	return &GeographyHandler{
		sv: sv,
	}
}

type GeographyHandler struct {
	sv service.GeographyService
}

func (h *GeographyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rg models.RequestGeography
	err := request.JSON(r, &rg)
	if err != nil {
		switch {
		case errors.Is(err, request.ErrRequestContentTypeNotJSON):
			response.Error(w, http.StatusUnsupportedMediaType, err.Error())
		case errors.Is(err, request.ErrRequestJSONInvalid):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = validators.ValidateGeographyPost(rg)
	if handleApiError(w, err) {
		return
	}

	s, err := h.sv.Create(ctx, rg)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusCreated, s)
}
