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

	er := validators.ValidateRequestBuyer(br)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	b, er := h.sv.Create(br)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusCreated, b)
}

func (h *BuyerHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID parameter")
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

	if er := validators.ValidateUpdateBuyer(br); er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	_, er := h.sv.FindById(id)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	updated, er := h.sv.Update(id, br)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusOK, updated)

}

func (h *BuyerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	er := h.sv.Delete(id)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)

}

func (h *BuyerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.sv.FindAll())

}

func (h *BuyerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	b, er := h.sv.FindById(id)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusOK, b)
}
