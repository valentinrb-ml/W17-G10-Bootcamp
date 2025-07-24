package employee

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Handler para las rutas relacionadas a Employee
type EmployeeHandler struct {
	service service.EmployeeService
}

// Constructor del handler
func NewEmployeeHandler(s service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

// POST /employees - crea un nuevo empleado
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.EmployeeRequest
	if err := request.JSON(r, &req); err != nil {
		response.Error(w, err)
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
		response.Error(w, err)
		return
	}
	// Convierte el modelo a doc para presentarlo al cliente
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(created)
	response.JSON(w, http.StatusCreated, employeeDoc)
}

// GET /employees - devuelve todos los empleados
func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	employees, err := h.service.FindAll(r.Context())
	if err != nil {
		response.Error(w, err)
		return
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
	if err != nil {
		response.Error(w, err)
		return
	}
	emp, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(emp)
	response.JSON(w, http.StatusOK, employeeDoc)
}

// PATCH /employees/{id} - actualiza parcialmente un empleado
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(w, err)
		return
	}
	var patch models.EmployeePatch
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&patch); err != nil {
		response.Error(w, err)
		return
	}
	updated, err := h.service.Update(r.Context(), id, &patch)
	if err != nil {
		response.Error(w, err)
		return
	}
	employeeDoc := mappers.MapEmployeeToEmployeeDoc(updated)
	response.JSON(w, http.StatusOK, employeeDoc)
}

// DELETE /employees/{id} - elimina un empleado por id
func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
