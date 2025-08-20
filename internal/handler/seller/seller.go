package handler

import (
	"net/http"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

// SellerHandler provides HTTP handlers for managing seller resources.
type SellerHandler struct {
	sv     service.SellerService
	logger logger.Logger
}

// NewSellerHandler creates a new SellerHandler using the given SellerService.
func NewSellerHandler(sv service.SellerService) *SellerHandler {
	return &SellerHandler{
		sv: sv,
	}
}

// SetLogger allows injecting the logger after creation
func (h *SellerHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

// Create handles HTTP POST requests for creating a new seller.
// It validates the request payload and returns the created seller as JSON.
func (h *SellerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Create seller request received")
	}

	var sr models.RequestSeller
	if err := httputil.DecodeJSON(r, &sr); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Invalid JSON in create request", map[string]interface{}{
				"error": err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	err := validators.ValidateSellerPost(sr)
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Validation failed for create request", map[string]interface{}{
				"cid":              sr.Cid,
				"company_name":     sr.CompanyName,
				"validation_error": err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	s, err := h.sv.Create(ctx, sr)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "seller-handler", "Failed to create seller", err, map[string]interface{}{
				"cid":          sr.Cid,
				"company_name": sr.CompanyName,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Seller created successfully", map[string]interface{}{
			"seller_id":    s.Id,
			"cid":          s.Cid,
			"company_name": s.CompanyName,
		})
	}

	response.JSON(w, http.StatusCreated, s)
}

// Update handles HTTP PATCH or PUT requests to update an existing seller by ID.
// It validates the request payload and returns the updated seller as JSON.
func (h *SellerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Invalid seller ID in update request", map[string]interface{}{
				"id_param": r.URL.Query().Get("id"),
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Update seller request received", map[string]interface{}{
			"seller_id": id,
		})
	}

	var sr models.RequestSeller
	if err := httputil.DecodeJSON(r, &sr); err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Invalid JSON in update request", map[string]interface{}{
				"seller_id": id,
				"error":     err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	err = validators.ValidateSellerPatchNotEmpty(sr)
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Empty validation failed for update request", map[string]interface{}{
				"seller_id":        id,
				"validation_error": err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	err = validators.ValidateSellerPatch(sr)
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Validation failed for update request", map[string]interface{}{
				"seller_id":        id,
				"cid":              sr.Cid,
				"company_name":     sr.CompanyName,
				"validation_error": err.Error(),
			})
		}
		response.Error(w, err)
		return
	}

	s, err := h.sv.Update(ctx, id, sr)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "seller-handler", "Failed to update seller", err, map[string]interface{}{
				"seller_id":    id,
				"cid":          sr.Cid,
				"company_name": sr.CompanyName,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Seller updated successfully", map[string]interface{}{
			"seller_id":    s.Id,
			"cid":          s.Cid,
			"company_name": s.CompanyName,
		})
	}

	response.JSON(w, http.StatusOK, s)
}

// Delete handles HTTP DELETE requests to remove a seller by ID.
// It returns HTTP 204 No Content if the seller is deleted successfully.
func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Invalid seller ID in delete request", map[string]interface{}{
				"id_param": r.URL.Query().Get("id"),
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Delete seller request received", map[string]interface{}{
			"seller_id": id,
		})
	}

	err = h.sv.Delete(ctx, id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "seller-handler", "Failed to delete seller", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Seller deleted successfully", map[string]interface{}{
			"seller_id": id,
		})
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// FindAll handles HTTP GET requests to retrieve all sellers.
// It returns a JSON array of sellers.
func (h *SellerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Find all sellers request received")
	}

	s, err := h.sv.FindAll(ctx)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "seller-handler", "Failed to find all sellers", err, nil)
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Find all sellers completed successfully", map[string]interface{}{
			"sellers_count": len(s),
		})
	}

	response.JSON(w, http.StatusOK, s)
}

// FindById handles HTTP GET requests to retrieve a single seller by ID.
// It returns the seller as JSON, or an error if not found.
func (h *SellerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(ctx, "seller-handler", "Invalid seller ID in find by ID request", map[string]interface{}{
				"id_param": r.URL.Query().Get("id"),
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Find seller by ID request received", map[string]interface{}{
			"seller_id": id,
		})
	}

	s, err := h.sv.FindById(ctx, id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "seller-handler", "Failed to find seller by ID", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		response.Error(w, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(ctx, "seller-handler", "Seller found successfully", map[string]interface{}{
			"seller_id":    s.Id,
			"cid":          s.Cid,
			"company_name": s.CompanyName,
		})
	}

	response.JSON(w, http.StatusOK, s)
}
