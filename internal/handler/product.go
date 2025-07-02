package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"net/http"
	"strconv"
)

type ProductHandler struct{ svc service.ProductService }

func NewProductHandler(s service.ProductService) *ProductHandler { return &ProductHandler{svc: s} }

func (h *ProductHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.getAll)
	r.Get("/{id}", h.getByID)
	r.Post("/", h.create)
	r.Patch("/{id}", h.patch)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	list, err := h.svc.GetAll(r.Context())
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, list)
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req product.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := validators.ValidateCreateRequest(req); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	product := mappers.ToDomain(req)

	result, err := h.svc.Create(r.Context(), product)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := h.parseID(r)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	product, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) parseID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, errors.New("id parameter is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("id must be a valid integer")
	}

	if id <= 0 {
		return 0, errors.New("id must be a positive integer")
	}

	return id, nil
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := h.parseID(r)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.svc.Delete(r.Context(), id)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *ProductHandler) patch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := h.parseID(r)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var req product.ProductPatchRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err = validators.ValidatePatchRequest(req); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	result, err := h.svc.Patch(r.Context(), id, req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func HandleServiceError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var se api.ServiceError
	if errors.As(err, &se) {
		response.Error(w, se.ResponseCode, se.Message)
		return
	}

	response.Error(w, http.StatusInternalServerError, "unexpected error")
}
