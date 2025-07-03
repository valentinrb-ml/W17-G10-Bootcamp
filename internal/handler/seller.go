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

	er := validators.ValidateRequestSeller(sr)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	s, er := h.sv.Create(sr)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusCreated, s)
}

func (h *SellerHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	er := validators.ValidateRequestSellerToPatch(sr)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	s, er := h.sv.Update(id, sr)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusOK, s)
}

func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID param.")
		return
	}

	er := h.sv.Delete(id)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *SellerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.sv.FindAll())
}

func (h *SellerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID param.")
		return
	}

	s, er := h.sv.FindById(id)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusOK, s)
}
