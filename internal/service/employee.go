package service

import (
	"context"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

type EmployeeService interface {
	Create(ctx context.Context, e *models.Employee) (*models.Employee, error)
	FindAll(ctx context.Context) ([]*models.Employee, error)
	FindByID(ctx context.Context, id int) (*models.Employee, error)
	Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error)
	Delete(ctx context.Context, id int) error
}

type EmployeeDefault struct {
	repo          repository.EmployeeRepository
	warehouseRepo repository.WarehouseRepository
}

func NewEmployeeDefault(r repository.EmployeeRepository, wrepo repository.WarehouseRepository) *EmployeeDefault {
	return &EmployeeDefault{
		repo:          r,
		warehouseRepo: wrepo,
	}
}

func (s *EmployeeDefault) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	if err := validators.ValidateEmployee(e); err != nil {
		return nil, err
	}

	warehouse, whErr := s.warehouseRepo.FindById(e.WarehouseID)
	if whErr != nil {
		var se *api.ServiceError
		if errors.As(whErr, &se) && se.Code == api.ErrNotFound {
			se := api.ServiceErrors[api.ErrBadRequest]
			se.Message = "warehouse_id does not exist"
			return nil, &se
		}
		return nil, whErr
	}
	if warehouse == nil {
		se := api.ServiceErrors[api.ErrBadRequest]
		se.Message = "warehouse_id does not exist"
		return nil, &se
	}

	emp, err := s.repo.FindByCardNumberID(ctx, e.CardNumberID)
	if err != nil {
		return nil, err
	}
	if emp != nil {
		se := api.ServiceErrors[api.ErrBadRequest]
		se.Message = "card_number_id already exists"
		return nil, &se
	}
	return s.repo.Create(ctx, e)
}

func (s *EmployeeDefault) FindAll(ctx context.Context) ([]*models.Employee, error) {
	return s.repo.FindAll(ctx)
}

func (s *EmployeeDefault) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	if err := validators.ValidateEmployeeID(id); err != nil {
		return nil, err
	}
	emp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return nil, &se
	}
	return emp, nil
}

func (s *EmployeeDefault) Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
	if id <= 0 {
		se := api.ServiceErrors[api.ErrUnprocessableEntity]
		se.Message = "invalid employee id"
		return nil, &se
	}
	if err := validators.ValidateEmployeePatch(patch); err != nil {
		return nil, err
	}

	if patch.WarehouseID != nil {
		warehouse, whErr := s.warehouseRepo.FindById(*patch.WarehouseID)
		if whErr != nil {
			var se *api.ServiceError
			if errors.As(whErr, &se) && se.Code == api.ErrNotFound {
				se := api.ServiceErrors[api.ErrBadRequest]
				se.Message = "warehouse_id does not exist"
				return nil, &se
			}
			return nil, whErr
		}
		if warehouse == nil {
			se := api.ServiceErrors[api.ErrBadRequest]
			se.Message = "warehouse_id does not exist"
			return nil, &se
		}
	}

	found, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if found == nil {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return nil, &se
	}

	updated, err := s.repo.Update(ctx, id, patch)
	if err != nil {
		switch err.Error() {
		case "card_number_id already exists":
			se := api.ServiceErrors[api.ErrBadRequest]
			se.Message = err.Error()
			return nil, &se
		case "not found":
			se := api.ServiceErrors[api.ErrNotFound]
			se.Message = "employee not found"
			return nil, &se
		default:
			se := api.ServiceErrors[api.ErrInternalServer]
			se.Message = "update failed"
			se.InternalError = err
			return nil, &se
		}
	}
	return updated, nil
}

func (s *EmployeeDefault) Delete(ctx context.Context, id int) error {
	if err := validators.ValidateEmployeeID(id); err != nil {
		return err
	}
	found, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if found == nil {
		se := api.ServiceErrors[api.ErrNotFound]
		se.Message = "employee not found"
		return &se
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		se := api.ServiceErrors[api.ErrInternalServer]
		se.Message = "failed to delete employee"
		se.InternalError = err
		return &se
	}
	return nil
}
