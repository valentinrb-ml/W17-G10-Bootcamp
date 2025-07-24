package service

import (
	"context"
	"errors"

	empRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
	wRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Servicio principal de empleados, implementa operaciones de negocio.
type EmployeeDefault struct {
	repo          empRepo.EmployeeRepository
	warehouseRepo wRepo.WarehouseRepository
}

// Constructor del servicio de empleados
func NewEmployeeDefault(r empRepo.EmployeeRepository, wrepo wRepo.WarehouseRepository) *EmployeeDefault {
	return &EmployeeDefault{
		repo:          r,
		warehouseRepo: wrepo,
	}
}

// Crea un nuevo empleado validando unicidad, existencia de warehouse y reglas de negocio.
func (s *EmployeeDefault) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	if err := validators.ValidateEmployee(e); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	warehouse, whErr := s.warehouseRepo.FindById(ctx, e.WarehouseID)
	if whErr != nil {
		var appErr *apperrors.AppError
		if errors.As(whErr, &appErr) && appErr.Code == apperrors.CodeNotFound {
			return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
		}
		return nil, apperrors.Wrap(whErr, "failed getting warehouse by id")
	}
	if warehouse == nil {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
	}
	emp, err := s.repo.FindByCardNumberID(ctx, e.CardNumberID)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed checking card_number_id uniqueness")
	}
	if emp != nil {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "card_number_id already exists")
	}
	return s.repo.Create(ctx, e)
}

// Devuelve todos los empleados
func (s *EmployeeDefault) FindAll(ctx context.Context) ([]*models.Employee, error) {
	emps, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed fetching all employees")
	}
	return emps, nil
}

// Busca un empleado por id, validando id y existencia
func (s *EmployeeDefault) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	if err := validators.ValidateEmployeeID(id); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	emp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed fetching employee by id")
	}
	if emp == nil {
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	return emp, nil
}

// Actualiza parcialmente un empleado, validando campos y relaciones.
func (s *EmployeeDefault) Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
	if err := validators.ValidateEmployeePatch(patch); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	if id <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, "invalid employee id")
	}
	found, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed fetching employee by id")
	}
	if found == nil {
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if patch.CardNumberID != nil {
		emp, err := s.repo.FindByCardNumberID(ctx, *patch.CardNumberID)
		if err != nil {
			return nil, apperrors.Wrap(err, "failed checking card_number_id")
		}
		if emp != nil && emp.ID != id {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "card_number_id already exists")
		}
		found.CardNumberID = *patch.CardNumberID
	}
	if patch.FirstName != nil {
		found.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		found.LastName = *patch.LastName
	}
	if patch.WarehouseID != nil && *patch.WarehouseID != 0 {
		warehouse, whErr := s.warehouseRepo.FindById(ctx, *patch.WarehouseID)
		if whErr != nil {
			var appErr *apperrors.AppError
			if errors.As(whErr, &appErr) && appErr.Code == apperrors.CodeNotFound {
				return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
			}
			return nil, apperrors.Wrap(whErr, "failed getting warehouse by id")
		}
		if warehouse == nil {
			return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
		}
		found.WarehouseID = *patch.WarehouseID
	}
	if err := s.repo.Update(ctx, id, found); err != nil {
		return nil, apperrors.Wrap(err, "failed updating employee")
	}
	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed fetching employee after update")
	}
	return updated, nil
}

// Elimina un empleado por id, validando su existencia.
func (s *EmployeeDefault) Delete(ctx context.Context, id int) error {
	if err := validators.ValidateEmployeeID(id); err != nil {
		return apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	found, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.Wrap(err, "failed fetching employee by id")
	}
	if found == nil {
		return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.Wrap(err, "failed deleting employee")
	}
	return nil
}
