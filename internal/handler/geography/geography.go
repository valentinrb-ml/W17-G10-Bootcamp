package handler

import (
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

// GeographyHandler provides HTTP handlers for managing geography resources such as countries, provinces, and localities.
type GeographyHandler struct {
	sv service.GeographyService
}

// NewGeographyHandler creates a new GeographyHandler using the given GeographyService.
func NewGeographyHandler(sv service.GeographyService) *GeographyHandler {
	return &GeographyHandler{
		sv: sv,
	}
}

// Create handles HTTP POST requests to create a new geography resource (country, province, or locality).
// It validates the request payload and returns the created resource as JSON.
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

// CountSellersByLocality handles HTTP GET requests to return the number of sellers by locality.
// - If 'id' is provided as a query parameter, it returns the count for that locality.
// - If 'id' is not provided, it returns the grouped counts for all localities.
func (h *GeographyHandler) CountSellersByLocality(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")

	if id == "" {
		resp, err := h.sv.CountSellersGroupedByLocality(ctx)
		if err != nil {
			response.Error(w, err)
			return
		}
		response.JSON(w, http.StatusOK, resp)
		return
	}

	s, err := h.sv.CountSellersByLocality(ctx, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, s)
}
