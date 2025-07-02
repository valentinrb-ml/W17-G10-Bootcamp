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

func NewWarehouseDefault(sv service.WarehouseService) *WarehouseDefault {
	return &WarehouseDefault{sv: sv}
}

type WarehouseDefault struct {
	// sv is the service that will be used by the handler
	sv service.WarehouseService
}

func (h *WarehouseDefault) Create(w http.ResponseWriter, r *http.Request) {
	var req = warehouse.WarehouseRequest{}
	if err := request.JSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.ValidateWarehouseCreateRequest(req); err != nil {
		response.Error(w, err.ResponseCode, err.Message)
		return
	}

	wh := mappers.RequestToWarehouse(req)

	newW, err := h.sv.Create(wh)
	if err != nil {
		response.Error(w, err.ResponseCode, err.Message)
		return
	}

	response.JSON(w, http.StatusCreated, mappers.WarehouseToDoc(newW))
}

func (h *WarehouseDefault) FindAll(w http.ResponseWriter, r *http.Request) {
	whs, err := h.sv.FindAll()
	if err != nil {
		response.Error(w, err.ResponseCode, err.Message)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDocSlice(whs))
}

func (h *WarehouseDefault) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := api.ServiceErrors[api.ErrBadRequest]
		response.Error(w, err.ResponseCode, "Invalid id")
		return
	}

	wh, er := h.sv.FindById(id)
	if er != nil {
		response.Error(w, er.ResponseCode, er.Message)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDoc(wh))

}

func (h *WarehouseDefault) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := api.ServiceErrors[api.ErrBadRequest]
		response.Error(w, err.ResponseCode, "Invalid id")
		return
	}

	var req warehouse.WarehousePatchDTO
	if err := request.JSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, serviceErr := h.sv.Update(id, req)
	if serviceErr != nil {
		response.Error(w, serviceErr.ResponseCode, serviceErr.Message)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDoc(updated))
}

func (h *WarehouseDefault) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := api.ServiceErrors[api.ErrBadRequest]
		response.Error(w, err.ResponseCode, "Invalid id")
		return
	}

	serviceErr := h.sv.Delete(id)
	if serviceErr != nil {
		response.Error(w, serviceErr.ResponseCode, serviceErr.Message)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
