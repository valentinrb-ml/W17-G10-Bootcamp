package handler

import (
	"net/http"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

// SellerHandler provides HTTP handlers for managing seller resources.
type SellerHandler struct {
	sv service.SellerService
}

// NewSellerHandler creates a new SellerHandler using the given SellerService.
func NewSellerHandler(sv service.SellerService) *SellerHandler {
	return &SellerHandler{
		sv: sv,
	}
}

// Create handles HTTP POST requests for creating a new seller.
// It validates the request payload and returns the created seller as JSON.
func (h *SellerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sr models.RequestSeller
	if err := httputil.DecodeJSON(r, &sr); err != nil {
		response.Error(w, err)
		return
	}

	err := validators.ValidateSellerPost(sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	s, err := h.sv.Create(ctx, sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, s)
}

// Update handles HTTP PATCH or PUT requests to update an existing seller by ID.
// It validates the request payload and returns the updated seller as JSON.
func (h *SellerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var sr models.RequestSeller
	if err := httputil.DecodeJSON(r, &sr); err != nil {
		response.Error(w, err)
		return
	}

	err = validators.ValidateSellerPatchNotEmpty(sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	err = validators.ValidateSellerPatch(sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	s, err := h.sv.Update(ctx, id, sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, s)
}

// Delete handles HTTP DELETE requests to remove a seller by ID.
// It returns HTTP 204 No Content if the seller is deleted successfully.
func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	err = h.sv.Delete(ctx, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// FindAll handles HTTP GET requests to retrieve all sellers.
// It returns a JSON array of sellers.
func (h *SellerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s, err := h.sv.FindAll(ctx)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, s)
}

// FindById handles HTTP GET requests to retrieve a single seller by ID.
// It returns the seller as JSON, or an error if not found.
func (h *SellerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	s, err := h.sv.FindById(ctx, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, s)
}
