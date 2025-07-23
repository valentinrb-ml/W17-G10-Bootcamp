package product

import (
	productMappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product"
	productService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductHandler struct{ svc productService.ProductService }

func NewProductHandler(s productService.ProductService) *ProductHandler {
	return &ProductHandler{svc: s}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.GetAll(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, list)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.ProductRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validators.ValidateCreateRequest(req); err != nil {
		response.Error(w, err)
		return
	}

	newProduct := productMappers.ToDomain(req)

	result, err := h.svc.Create(r.Context(), newProduct)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	if err = validators.ValidateID(id, "product id"); err != nil {
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
	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	if err = validators.ValidateID(id, "product id"); err != nil {
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
	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	if err = validators.ValidateID(id, "product id"); err != nil {
		response.Error(w, err)
		return
	}

	var req models.ProductPatchRequest
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
