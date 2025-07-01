package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Interface
type EmployeeRepository interface {
	Create(e *models.Employee) (*models.Employee, error)
	SaveToFile(filename string) error
	FindByCardNumberID(cardNumberID string) (*models.Employee, error)
	FindAll() ([]*models.Employee, error)
	FindByID(id int) (*models.Employee, error)
	Update(id int, patch *models.EmployeePatch) (*models.Employee, error)
	Delete(id int) error
}

// Implementación en memoria
type EmployeeMap struct {
	mu     sync.Mutex
	nextID int
	data   map[int]*models.Employee
}

func NewEmployeeMap() *EmployeeMap {
	return &EmployeeMap{
		nextID: 1,
		data:   make(map[int]*models.Employee),
	}
}

func (r *EmployeeMap) Create(e *models.Employee) (*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, emp := range r.data {
		if emp.CardNumberID == e.CardNumberID {
			return nil, errors.New("card_number_id already exists")
		}
	}
	e.ID = r.nextID
	r.nextID++
	r.data[e.ID] = e
	return e, nil
}

func (r *EmployeeMap) SaveToFile(filename string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var employees []models.EmployeeDoc
	for _, emp := range r.data {
		employees = append(employees, mappers.MapEmployeeToEmployeeDoc(emp))
	}
	file, err := os.Create(filename)
	if err != nil {
		return errors.New("could not create the file")
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(employees)
}

func (r *EmployeeMap) FindByCardNumberID(cardNumberID string) (*models.Employee, error) {
	for _, emp := range r.data {
		if emp.CardNumberID == cardNumberID {
			return emp, nil
		}
	}
	return nil, nil
}

func (r *EmployeeMap) FindAll() ([]*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	employees := make([]*models.Employee, 0, len(r.data))
	for _, emp := range r.data {
		employees = append(employees, emp)
	}
	return employees, nil
}
func (r *EmployeeMap) FindByID(id int) (*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	emp, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return emp, nil
}
func (r *EmployeeMap) Update(id int, patch *models.EmployeePatch) (*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	emp, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	if patch.CardNumberID != nil && *patch.CardNumberID != emp.CardNumberID {
		for _, v := range r.data {
			if v.CardNumberID == *patch.CardNumberID {
				return nil, errors.New("card_number_id already exists")
			}
		}
		emp.CardNumberID = *patch.CardNumberID
	}
	if patch.FirstName != nil {
		emp.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		emp.LastName = *patch.LastName
	}
	if patch.WarehouseID != nil && *patch.WarehouseID != 0 {
		emp.WarehouseID = *patch.WarehouseID
	}
	r.data[emp.ID] = emp
	return emp, nil
}
func (r *EmployeeMap) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[id]; !exists {
		return errors.New("not found")
	}
	delete(r.data, id)
	return nil
}
