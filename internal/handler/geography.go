package handler

import (
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
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
	if err := httputil.DecodeJSON(r, &rg); err != nil {
		response.Error(w, err)
		return
	}

	err := validators.ValidateGeographyPost(rg)
	if err != nil {
		response.Error(w, err)
		return
	}

	s, err := h.sv.Create(ctx, rg)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, s)
}
