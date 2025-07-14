package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func NewWarehouseDefault(sv service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{sv: sv}
}

type WarehouseHandler struct {
	// sv is the service that will be used by the handler
	sv service.WarehouseService
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req = warehouse.WarehouseRequest{}
	if err := request.JSON(r, &req); err != nil {
		// response.Error(w, http.StatusBadRequest, err.Error())
		response.Error(w, err)
		return
	}

	if err := validators.ValidateWarehouseCreateRequest(req); err != nil {
		// response.Error(w, err.ResponseCode, err.Message)
		response.Error(w, err)
		return
	}

	wh := mappers.RequestToWarehouse(req)

	newW, err := h.sv.Create(r.Context(), wh)
	if err != nil {
		// response.Error(w, err.ResponseCode, err.Message)
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, mappers.WarehouseToDoc(newW))
}

func (h *WarehouseHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	whs, err := h.sv.FindAll(r.Context())
	if err != nil {
		// response.Error(w, err.ResponseCode, err.Message)
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDocSlice(whs))
}

func (h *WarehouseHandler) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := api.ServiceErrors[api.ErrBadRequest]
		err.Message = "Invalid id"
		// response.Error(w, err.ResponseCode, err.Message)
		response.Error(w, err)
		return
	}

	wh, er := h.sv.FindById(r.Context(), id)
	if er != nil {
		// response.Error(w, er.ResponseCode, er.Message)
		response.Error(w, er)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDoc(wh))

}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := api.ServiceErrors[api.ErrBadRequest]
		err.Message = "Invalid id"
		// response.Error(w, err.ResponseCode, err.Message)
		response.Error(w, err)
		return
	}

	var req warehouse.WarehousePatchDTO
	if err := request.JSON(r, &req); err != nil {
		// response.Error(w, http.StatusBadRequest, err.Error())
		response.Error(w, err)
		return
	}

	updated, serviceErr := h.sv.Update(r.Context(), id, req)
	if serviceErr != nil {
		// response.Error(w, serviceErr.ResponseCode, serviceErr.Message)
		response.Error(w, serviceErr)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDoc(updated))
}

func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := api.ServiceErrors[api.ErrBadRequest]
		err.Message = "Invalid id"
		// response.Error(w, err.ResponseCode, err.Message)
		response.Error(w, err)
		return
	}

	serviceErr := h.sv.Delete(r.Context(), id)
	if serviceErr != nil {
		// response.Error(w, serviceErr.ResponseCode, serviceErr.Message)
		response.Error(w, serviceErr)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
