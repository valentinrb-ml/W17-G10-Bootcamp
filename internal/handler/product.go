package handler

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"net/http"
)

type ProductHandler struct{ svc service.ProductService }

func NewProductHandler(s service.ProductService) *ProductHandler { return &ProductHandler{svc: s} }

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.GetAll(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, list)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req product.ProductRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validators.ValidateCreateRequest(req); err != nil {
		response.Error(w, err)
		return
	}

	newProduct := mappers.ToDomain(req)

	result, err := h.svc.Create(r.Context(), newProduct)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	currentProduct, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, currentProduct)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	err = h.svc.Delete(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *ProductHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var req product.ProductPatchRequest
	if err = httputil.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	if err = validators.ValidatePatchRequest(req); err != nil {
		response.Error(w, err)
		return
	}

	result, err := h.svc.Patch(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
