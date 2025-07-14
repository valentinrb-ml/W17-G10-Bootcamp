package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(s service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

type EmployeeRequest struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  *int   `json:"warehouse_id"`
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req EmployeeRequest
	if err := request.JSON(r, &req); err != nil {
		// response.Error(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		response.Error(w, err)
		return
	}
	emp := &models.Employee{
		CardNumberID: req.CardNumberID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}
	if req.WarehouseID != nil {
		emp.WarehouseID = *req.WarehouseID
	} else {
		emp.WarehouseID = 0
	}
	created, err := h.service.Create(r.Context(), emp)
	if err != nil {
		var se *api.ServiceError
		if errors.As(err, &se) {
			// response.Error(w, se.ResponseCode, se.Message)
			response.Error(w, err)
		} else {
			// response.Error(w, http.StatusInternalServerError, "Internal error")
			response.Error(w, err)
		}
		return
	}
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(created)
	response.JSON(w, http.StatusCreated, employeeDoc)
}

func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	employees, err := h.service.FindAll(r.Context())
	if err != nil {
		// response.Error(w, http.StatusInternalServerError, "cannot get employees")
		response.Error(w, err)
		return
	}
	if len(employees) == 0 {
		response.JSON(w, http.StatusOK, []interface{}{})
		return
	}
	employeeDocs := make([]models.EmployeeDoc, 0, len(employees))
	for _, emp := range employees {
		employeeDocs = append(employeeDocs, mappers.MapEmployeeToEmployeeDoc(emp))
	}
	response.JSON(w, http.StatusOK, employeeDocs)
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// response.Error(w, http.StatusBadRequest, "id must be a number")
		response.Error(w, err)
		return
	}
	emp, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		var se *api.ServiceError
		if errors.As(err, &se) {
			// response.Error(w, se.ResponseCode, se.Message)
			response.Error(w, err)
		} else {
			// response.Error(w, http.StatusInternalServerError, "Internal error")
			response.Error(w, err)
		}
		return
	}
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(emp)
	response.JSON(w, http.StatusOK, employeeDoc)
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// response.Error(w, http.StatusBadRequest, "Invalid id")
		response.Error(w, err)
		return
	}
	var patch models.EmployeePatch
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&patch); err != nil {
		// response.Error(w, http.StatusBadRequest, "Invalid JSON or unknown field: "+err.Error())
		response.Error(w, err)
		return
	}
	updated, err := h.service.Update(r.Context(), id, &patch)
	if err != nil {
		var se *api.ServiceError
		if errors.As(err, &se) {
			// response.Error(w, se.ResponseCode, se.Message)
			response.Error(w, err)
		} else {
			// response.Error(w, http.StatusInternalServerError, "Internal error")
			response.Error(w, err)
		}
		return
	}
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(updated)
	response.JSON(w, http.StatusOK, employeeDoc)
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// response.Error(w, http.StatusBadRequest, "Invalid id")
		response.Error(w, err)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		var se *api.ServiceError
		if errors.As(err, &se) {
			// response.Error(w, se.ResponseCode, se.Message)
			response.Error(w, err)
		} else {
			// response.Error(w, http.StatusInternalServerError, "Internal error")
			response.Error(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
