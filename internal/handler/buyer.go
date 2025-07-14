package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func NewBuyerHandler(sv service.BuyerService) *BuyerHandler {
	return &BuyerHandler{
		sv: sv,
	}
}

type BuyerHandler struct {
	sv service.BuyerService
}

func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var br models.RequestBuyer
	err := request.JSON(r, &br)
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

	err = validators.ValidateRequestBuyer(br)
	if handleApiError(w, err) {
		return
	}

	b, err := h.sv.Create(ctx, br)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusCreated, b)
}

func (h *BuyerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID param.")
		return
	}

	var br models.RequestBuyer
	err = request.JSON(r, &br)
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

	// Validación PATCH puede partirse igual que seller si quieres más granularidad
	err = validators.ValidateUpdateBuyer(br)
	if handleApiError(w, err) {
		return
	}

	updated, err := h.sv.Update(ctx, id, br)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusOK, updated)
}

func (h *BuyerHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *BuyerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bs, err := h.sv.FindAll(ctx)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusOK, bs)
}

func (h *BuyerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID param.")
		return
	}

	b, err := h.sv.FindById(ctx, id)
	if handleApiError(w, err) {
		return
	}

	response.JSON(w, http.StatusOK, b)
}
