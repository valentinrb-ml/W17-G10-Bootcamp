package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Handler para las rutas relacionadas a Employee
type EmployeeHandler struct {
	service service.EmployeeService
	logger  logger.Logger
}

// Constructor del handler
func NewEmployeeHandler(s service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}
func (h *EmployeeHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

// POST /employees - crea un nuevo empleado
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Create employee request received")
	}
	var req models.EmployeeRequest
	if err := request.JSON(r, &req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "employee-handler", "Invalid JSON in create request", map[string]interface{}{
				"error": err.Error(),
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	// Construye el objeto empleado a partir de la request
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
	// Llama al service para crear
	created, err := h.service.Create(r.Context(), emp)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "employee-handler", "Failed to create employee", err, map[string]interface{}{
				"card_number_id": emp.CardNumberID,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Employee created successfully", map[string]interface{}{
			"employee_id":    created.ID,
			"card_number_id": created.CardNumberID,
		})
	}
	// Convierte el modelo a doc para presentarlo al cliente
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(created)
	response.JSON(w, http.StatusCreated, employeeDoc)
}

// GET /employees - devuelve todos los empleados
func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Get all employees request received")
	}
	employees, err := h.service.FindAll(r.Context())
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "employee-handler", "Failed to fetch all employees", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "All employees fetched", map[string]interface{}{
			"count": len(employees),
		})
	}
	if len(employees) == 0 {
		// Responde lista vacía si no hay empleados
		response.JSON(w, http.StatusOK, []interface{}{})
		return
	}
	employeeDocs := make([]models.EmployeeDoc, 0, len(employees))
	for _, emp := range employees {
		employeeDocs = append(employeeDocs, mappers.MapEmployeeToEmployeeDoc(emp))
	}
	response.JSON(w, http.StatusOK, employeeDocs)
}

// GET /employees/{id} - devuelve un empleado por id
func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Get employee by ID request received", map[string]interface{}{
			"employee_id": idParam,
		})
	}
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "employee-handler", "Invalid employee ID in request", map[string]interface{}{
				"id_param": idParam,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	emp, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "employee-handler", "Failed to get employee by ID", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Employee found by ID", map[string]interface{}{
			"employee_id": id,
		})
	}
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(emp)
	response.JSON(w, http.StatusOK, employeeDoc)
}

// PATCH /employees/{id} - actualiza parcialmente un empleado
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Update employee request received", map[string]interface{}{
			"employee_id": idParam,
		})
	}
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "employee-handler", "Invalid employee ID in update request", map[string]interface{}{
				"id_param": idParam,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	var patch models.EmployeePatch
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&patch); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "employee-handler", "Invalid JSON in update request", map[string]interface{}{
				"error":       err.Error(),
				"employee_id": id,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	updated, err := h.service.Update(r.Context(), id, &patch)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "employee-handler", "Failed to update employee", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Employee updated successfully", map[string]interface{}{
			"employee_id": updated.ID,
		})
	}
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(updated)
	response.JSON(w, http.StatusOK, employeeDoc)
}

// DELETE /employees/{id} - elimina un empleado por id
func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Delete employee request received", map[string]interface{}{
			"employee_id": idParam,
		})
	}
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "employee-handler", "Invalid employee ID in delete request", map[string]interface{}{
				"id_param": idParam,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "employee-handler", "Failed to delete employee", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if h.logger != nil {
		h.logger.Info(r.Context(), "employee-handler", "Employee deleted successfully", map[string]interface{}{
			"employee_id": id,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}
