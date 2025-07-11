package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func NewSellerHandler(sv service.SellerService) *SellerHandler {
	return &SellerHandler{
		sv: sv,
	}
}

type SellerHandler struct {
	sv service.SellerService
}

func (h *SellerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sr models.RequestSeller
	err := request.JSON(r, &sr)
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

	err = validators.ValidateSellerPost(sr)
	if handleApiError(w, err) {
		return
	}

	s, err := h.sv.Create(ctx, sr)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusCreated, s)
}

func (h *SellerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID param.")
		return
	}

	var sr models.RequestSeller
	err = request.JSON(r, &sr)
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

	err = validators.ValidateSellerPatchNotEmpty(sr)
	if handleApiError(w, err) {
		return
	}

	err = validators.ValidateSellerPatch(sr)
	if handleApiError(w, err) {
		return
	}

	s, err := h.sv.Update(ctx, id, sr)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusOK, s)
}

func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID param.")
		return
	}

	err = h.sv.Delete(ctx, id)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *SellerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s, err := h.sv.FindAll(ctx)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusOK, s)
}

func (h *SellerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID param.")
		return
	}

	s, err := h.sv.FindById(ctx, id)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusOK, s)
}

func handleApiError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if errorResp, ok := err.(*api.ServiceError); ok {
		response.Error(w, errorResp.ResponseCode, errorResp.Message)
	} else {
		response.Error(w, http.StatusInternalServerError, err.Error())
	}

	return true
}
