package service

import (
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Interface
type EmployeeService interface {
	Create(e *models.Employee) (*models.Employee, error)
	SaveToFile(filename string) error
	FindAll() ([]*models.Employee, error)
	FindByID(id int) (*models.Employee, error)
}

// Implementación
type EmployeeDefault struct {
	repo repository.EmployeeRepository
}

func NewEmployeeDefault(r repository.EmployeeRepository) *EmployeeDefault {
	return &EmployeeDefault{repo: r}
}

func (s *EmployeeDefault) Create(e *models.Employee) (*models.Employee, error) {
	if err := validators.ValidateEmployee(e); err != nil {
		return nil, err
	}
	emp, _ := s.repo.FindByCardNumberID(e.CardNumberID)
	if emp != nil {
		se := api.ServiceErrors[api.ErrBadRequest]
		se.InternalError = errors.New("card_number_id already exists")
		return nil, se
	}
	return s.repo.Create(e)
}

func (s *EmployeeDefault) SaveToFile(filename string) error {
	if repo, ok := s.repo.(*repository.EmployeeMap); ok {
		return repo.SaveToFile(filename)
	}
	return errors.New("repository does not support saving to file")
}

func (s *EmployeeDefault) FindAll() ([]*models.Employee, error) {
	return s.repo.FindAll()
}
func (s *EmployeeDefault) FindByID(id int) (*models.Employee, error) {
	if err := validators.ValidateEmployeeID(id); err != nil {
		return nil, err
	}
	emp, _ := s.repo.FindByID(id)
	if emp == nil {
		se := api.ServiceErrors[api.ErrNotFound]
		se.InternalError = nil
		se.Message = "employee not found"
		return nil, se
	}
	return emp, nil
}
