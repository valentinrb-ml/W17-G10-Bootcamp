package service

import (
	"context"
	"errors"

	empRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
	wRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

// Servicio principal de empleados, implementa operaciones de negocio.
type EmployeeDefault struct {
	repo          empRepo.EmployeeRepository
	warehouseRepo wRepo.WarehouseRepository
	logger        logger.Logger
}

// Constructor del servicio de empleados
func NewEmployeeDefault(r empRepo.EmployeeRepository, wrepo wRepo.WarehouseRepository) *EmployeeDefault {
	return &EmployeeDefault{
		repo:          r,
		warehouseRepo: wrepo,
	}
}
func (s *EmployeeDefault) SetLogger(l logger.Logger) {
	s.logger = l
}

// Crea un nuevo empleado validando unicidad, existencia de warehouse y reglas de negocio.
func (s *EmployeeDefault) Create(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Creating new employee", map[string]interface{}{
			"card_number_id": e.CardNumberID,
		})
	}
	if err := validators.ValidateEmployee(e); err != nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Validation failed for create request", map[string]interface{}{
				"validation_error": err.Error(),
				"card_number_id":   e.CardNumberID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	warehouse, whErr := s.warehouseRepo.FindById(ctx, e.WarehouseID)
	if whErr != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed getting warehouse by id", whErr, map[string]interface{}{
				"warehouse_id": e.WarehouseID,
			})
		}
		var appErr *apperrors.AppError
		if errors.As(whErr, &appErr) && appErr.Code == apperrors.CodeNotFound {
			return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
		}
		return nil, apperrors.Wrap(whErr, "failed getting warehouse by id")
	}
	if warehouse == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Warehouse id does not exist for employee creation", map[string]interface{}{
				"warehouse_id": e.WarehouseID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
	}
	emp, err := s.repo.FindByCardNumberID(ctx, e.CardNumberID)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed checking card_number_id uniqueness", err, map[string]interface{}{
				"card_number_id": e.CardNumberID,
			})
		}
		return nil, apperrors.Wrap(err, "failed checking card_number_id uniqueness")
	}
	if emp != nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "card_number_id already exists", map[string]interface{}{
				"card_number_id": e.CardNumberID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "card_number_id already exists")
	}
	newEmp, err := s.repo.Create(ctx, e)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed to create employee", err, map[string]interface{}{
				"card_number_id": e.CardNumberID,
			})
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Employee created successfully", map[string]interface{}{
			"employee_id":    newEmp.ID,
			"card_number_id": newEmp.CardNumberID,
		})
	}
	return newEmp, nil
}

// Devuelve todos los empleados
func (s *EmployeeDefault) FindAll(ctx context.Context) ([]*models.Employee, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Fetching all employees")
	}
	emps, err := s.repo.FindAll(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed fetching all employees", err)
		}
		return nil, apperrors.Wrap(err, "failed fetching all employees")
	}
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Successfully fetched all employees", map[string]interface{}{
			"count": len(emps),
		})
	}
	return emps, nil
}

// Busca un empleado por id, validando id y existencia
func (s *EmployeeDefault) FindByID(ctx context.Context, id int) (*models.Employee, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Fetching employee by id", map[string]interface{}{
			"employee_id": id,
		})
	}
	if err := validators.ValidateEmployeeID(id); err != nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Invalid employee id format", map[string]interface{}{
				"validation_error": err.Error(),
				"employee_id":      id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	emp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed fetching employee by id", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.Wrap(err, "failed fetching employee by id")
	}
	if emp == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Employee not found", map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Successfully fetched employee by id", map[string]interface{}{
			"employee_id": emp.ID,
		})
	}
	return emp, nil
}

// Actualiza parcialmente un empleado, validando campos y relaciones.
func (s *EmployeeDefault) Update(ctx context.Context, id int, patch *models.EmployeePatch) (*models.Employee, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Updating employee", map[string]interface{}{
			"employee_id": id,
		})
	}
	if err := validators.ValidateEmployeePatch(patch); err != nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Validation failed for patch", map[string]interface{}{
				"employee_id":      id,
				"validation_error": err.Error(),
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	if id <= 0 {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Invalid employee id for update", map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeValidationError, "invalid employee id")
	}
	found, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed fetching employee by id (update)", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.Wrap(err, "failed fetching employee by id")
	}
	if found == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Employee not found for update", map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if patch.CardNumberID != nil {
		emp, err := s.repo.FindByCardNumberID(ctx, *patch.CardNumberID)
		if err != nil {
			if s.logger != nil {
				s.logger.Error(ctx, "employee-service", "Failed checking card_number_id (update)", err, map[string]interface{}{
					"employee_id":    id,
					"card_number_id": *patch.CardNumberID,
				})
			}
			return nil, apperrors.Wrap(err, "failed checking card_number_id")
		}
		if emp != nil && emp.ID != id {
			if s.logger != nil {
				s.logger.Warning(ctx, "employee-service", "card_number_id already exists (update)", map[string]interface{}{
					"employee_id":    id,
					"card_number_id": *patch.CardNumberID,
				})
			}
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
			if s.logger != nil {
				s.logger.Error(ctx, "employee-service", "Failed getting warehouse by id (update)", whErr, map[string]interface{}{
					"employee_id":  id,
					"warehouse_id": *patch.WarehouseID,
				})
			}
			var appErr *apperrors.AppError
			if errors.As(whErr, &appErr) && appErr.Code == apperrors.CodeNotFound {
				return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
			}
			return nil, apperrors.Wrap(whErr, "failed getting warehouse by id")
		}
		if warehouse == nil {
			if s.logger != nil {
				s.logger.Warning(ctx, "employee-service", "warehouse_id does not exist (update)", map[string]interface{}{
					"employee_id":  id,
					"warehouse_id": *patch.WarehouseID,
				})
			}
			return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse_id does not exist")
		}
		found.WarehouseID = *patch.WarehouseID
	}
	if err := s.repo.Update(ctx, id, found); err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed updating employee", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.Wrap(err, "failed updating employee")
	}
	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed fetching employee after update", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return nil, apperrors.Wrap(err, "failed fetching employee after update")
	}
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Employee updated successfully", map[string]interface{}{
			"employee_id": updated.ID,
		})
	}
	return updated, nil
}

// Elimina un empleado por id, validando su existencia.
func (s *EmployeeDefault) Delete(ctx context.Context, id int) error {
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Deleting employee", map[string]interface{}{
			"employee_id": id,
		})
	}
	if err := validators.ValidateEmployeeID(id); err != nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Invalid employee id for delete", map[string]interface{}{
				"employee_id":      id,
				"validation_error": err.Error(),
			})
		}
		return apperrors.NewAppError(apperrors.CodeValidationError, err.Error())
	}
	found, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed fetching employee by id (delete)", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return apperrors.Wrap(err, "failed fetching employee by id")
	}
	if found == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "employee-service", "Employee not found for delete", map[string]interface{}{
				"employee_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "employee-service", "Failed deleting employee", err, map[string]interface{}{
				"employee_id": id,
			})
		}
		return apperrors.Wrap(err, "failed deleting employee")
	}
	if s.logger != nil {
		s.logger.Info(ctx, "employee-service", "Employee deleted successfully", map[string]interface{}{
			"employee_id": id,
		})
	}
	return nil
}
