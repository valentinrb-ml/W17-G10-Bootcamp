package loader

import (
	"encoding/json"
	"os"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Interface
type EmployeeLoader interface {
	Load() (map[int]*models.Employee, error)
}

// Implementación: JSON file loader
type EmployeeJSONFile struct {
	path string
}

func NewEmployeeJSONFile(path string) *EmployeeJSONFile {
	return &EmployeeJSONFile{path: path}
}

func (l *EmployeeJSONFile) Load() (map[int]*models.Employee, error) {
	file, err := os.Open(l.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var employeesJSON []models.EmployeeDoc
	err = json.NewDecoder(file).Decode(&employeesJSON)
	if err != nil {
		return nil, err
	}

	employees := make(map[int]*models.Employee)
	for _, doc := range employeesJSON {
		emp := mappers.MapEmployeeDocToEmployee(doc)
		employees[emp.ID] = emp
	}
	return employees, nil
}
