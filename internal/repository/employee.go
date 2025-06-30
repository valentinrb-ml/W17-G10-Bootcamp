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
