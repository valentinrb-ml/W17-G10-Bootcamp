package handler

import (
	"net/http"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

// GeographyHandler provides HTTP handlers for managing geography resources such as countries, provinces, and localities.
type GeographyHandler struct {
	sv     service.GeographyService
	logger logger.Logger
}

// NewGeographyHandler creates a new GeographyHandler using the given GeographyService.
func NewGeographyHandler(sv service.GeographyService) *GeographyHandler {
	return &GeographyHandler{
		sv: sv,
	}
}

// SetLogger allows injecting the logger after creation
func (h *GeographyHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

// Create handles HTTP POST requests to create a new geography resource (country, province, or locality).
// It validates the request payload and returns the created resource as JSON.
func (h *GeographyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(ctx, "geography-handler", "Create geography request received")
	}

	var rg models.RequestGeography
	if err := httputil.DecodeJSON(r, &rg); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "geography-handler", "Invalid JSON in create request", map[string]interface{}{
				"error": err.Error(),
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	err := validators.ValidateGeographyPost(rg)
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "geography-handler", "Validation failed for create request", map[string]interface{}{
				"locality_name":    rg.LocalityName,
				"province_name":    rg.ProvinceName,
				"country_name":     rg.CountryName,
				"validation_error": err.Error(),
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	s, err := h.sv.Create(ctx, rg)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "geography-handler", "Failed to create geography", err, map[string]interface{}{
				"locality_name": rg.LocalityName,
				"province_name": rg.ProvinceName,
				"country_name":  rg.CountryName,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "geography-handler", "Geography created successfully", map[string]interface{}{
			"locality_id":   s.LocalityId,
			"locality_name": s.LocalityName,
			"province_name": s.ProvinceName,
			"country_name":  s.CountryName,
		})
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
		if h.logger != nil {
			h.logger.Info(ctx, "geography-handler", "Count sellers grouped by locality request received")
		}

		resp, err := h.sv.CountSellersGroupedByLocality(ctx)
		if err != nil {
			if h.logger != nil {
				h.logger.Error(ctx, "geography-handler", "Failed to count sellers grouped by locality", err, nil)
			}
			response.ErrorWithRequest(w, r, err)
			return
		}

		if h.logger != nil {
			h.logger.Info(ctx, "geography-handler", "Count sellers grouped by locality completed successfully", map[string]interface{}{
				"localities_count": len(resp),
			})
		}

		response.JSON(w, http.StatusOK, resp)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "geography-handler", "Count sellers by locality request received", map[string]interface{}{
			"locality_id": id,
		})
	}

	s, err := h.sv.CountSellersByLocality(ctx, id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "geography-handler", "Failed to count sellers by locality", err, map[string]interface{}{
				"locality_id": id,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "geography-handler", "Count sellers by locality completed successfully", map[string]interface{}{
			"locality_id":   s.LocalityId,
			"locality_name": s.LocalityName,
			"sellers_count": s.SellersCount,
		})
	}

	response.JSON(w, http.StatusOK, s)
}
