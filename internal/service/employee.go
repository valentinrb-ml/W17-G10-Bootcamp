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
	Update(id int, patch *models.EmployeePatch) (*models.Employee, error)
	Delete(id int) error
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

func (s *EmployeeDefault) Update(id int, patch *models.EmployeePatch) (*models.Employee, error) {
	if id <= 0 {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.Message = "invalid employee id"
		return nil, se
	}
	if err := validators.ValidateEmployeePatch(patch); err != nil {
		return nil, err
	}
	found, _ := s.repo.FindByID(id)
	if found == nil {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return nil, se
	}
	updated, err := s.repo.Update(id, patch)
	if err != nil {
		switch err.Error() {
		case "card_number_id already exists":
			se := api.ServiceErrors[api.ErrBadRequest]
			se.Message = err.Error()
			return nil, se
		case "not found":
			se := api.ServiceErrors[api.ErrNotFound]
			se.Message = "employee not found"
			return nil, se
		default:
			se := api.ServiceErrors[api.ErrInternalServer]
			se.Message = "update failed"
			se.InternalError = err
			return nil, se
		}
	}
	return updated, nil
}
func (s *EmployeeDefault) Delete(id int) error {
	if err := validators.ValidateEmployeeID(id); err != nil {
		return err
	}
	found, _ := s.repo.FindByID(id)
	if found == nil {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return se
	}
	if err := s.repo.Delete(id); err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "failed to delete employee"
		se.InternalError = err
		return se
	}
	return nil
}
