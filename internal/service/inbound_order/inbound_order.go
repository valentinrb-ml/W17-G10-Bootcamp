package service

import (
	"context"
	"errors"

	empRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
	inbRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/inbound_order"
	wRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

// Servicio principal para Inbound Orders, incluye lógica de validación
type InboundOrderDefault struct {
	repo          inbRepo.InboundOrderRepository
	employeeRepo  empRepo.EmployeeRepository
	warehouseRepo wRepo.WarehouseRepository
	logger        logger.Logger
}

// Constructor del servicio de inbound orders
func NewInboundOrderService(
	r inbRepo.InboundOrderRepository,
	er empRepo.EmployeeRepository,
	wr wRepo.WarehouseRepository) *InboundOrderDefault {
	return &InboundOrderDefault{
		repo:          r,
		employeeRepo:  er,
		warehouseRepo: wr,
	}
}
func (s *InboundOrderDefault) SetLogger(l logger.Logger) {
	s.logger = l
}

// Crea un nuevo inbound order con todas las validaciones de negocio y de integridad
func (s *InboundOrderDefault) Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "inboundorder-service", "Creating inbound order", map[string]interface{}{
			"order_number": o.OrderNumber,
			"employee_id":  o.EmployeeID,
			"warehouse_id": o.WarehouseID,
		})
	}
	// Validación de campos obligatorios (devuelve 422 si falta alguno)
	if err := validators.ValidateInboundOrder(o); err != nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "inboundorder-service", "Validation failed for inbound order", map[string]interface{}{
				"validation_error": err.Error(),
				"order_number":     o.OrderNumber,
			})
		}
		return nil, err
	}
	// Verifica unicidad del order_number
	exists, err := s.repo.ExistsByOrderNumber(ctx, o.OrderNumber)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "inboundorder-service", "Failed checking order_number uniqueness", err, map[string]interface{}{
				"order_number": o.OrderNumber,
			})
		}
		return nil, apperrors.Wrap(err, "failed checking for order_number uniqueness")
	}
	if exists {
		if s.logger != nil {
			s.logger.Warning(ctx, "inboundorder-service", "order_number already exists", map[string]interface{}{
				"order_number": o.OrderNumber,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "order_number already exists")
	}
	// Valida que el empleado referenciado exista (FK)
	emp, err := s.employeeRepo.FindByID(ctx, o.EmployeeID)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "inboundorder-service", "Failed getting employee by id", err, map[string]interface{}{
				"employee_id": o.EmployeeID,
			})
		}
		return nil, apperrors.Wrap(err, "failed getting employee by id")
	}
	if emp == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "inboundorder-service", "employee_id does not exist", map[string]interface{}{
				"employee_id": o.EmployeeID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "employee_id does not exist")
	}
	// Valida que el warehouse referenciado exista (FK)
	warehouse, whErr := s.warehouseRepo.FindById(ctx, o.WarehouseID)
	if whErr != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "inboundorder-service", "Failed getting warehouse by id", whErr, map[string]interface{}{
				"warehouse_id": o.WarehouseID,
			})
		}
		var appErr *apperrors.AppError
		if errors.As(whErr, &appErr) && appErr.Code == apperrors.CodeNotFound {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_id does not exist")
		}
		return nil, apperrors.Wrap(whErr, "failed getting warehouse by id")
	}
	if warehouse == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "inboundorder-service", "warehouse_id does not exist", map[string]interface{}{
				"warehouse_id": o.WarehouseID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_id does not exist")
	}
	// Si pasa todas las validaciones, crea el inbound order
	inboundOrder, err := s.repo.Create(ctx, o)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "inboundorder-service", "Failed to create inbound order", err, map[string]interface{}{
				"order_number": o.OrderNumber,
			})
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.Info(ctx, "inboundorder-service", "Inbound order created successfully", map[string]interface{}{
			"inbound_order_id": inboundOrder.ID,
			"order_number":     inboundOrder.OrderNumber,
		})
	}
	return inboundOrder, nil
}

// Genera un reporte de inbound orders por empleado o global si employeeID es nil
func (s *InboundOrderDefault) Report(ctx context.Context, employeeID *int) (interface{}, error) {
	if s.logger != nil {
		if employeeID == nil {
			s.logger.Info(ctx, "inboundorder-service", "Generating inbound orders report for all employees")
		} else {
			s.logger.Info(ctx, "inboundorder-service", "Generating inbound orders report for employee", map[string]interface{}{
				"employee_id": *employeeID,
			})
		}
	}

	var (
		result interface{}
		err    error
	)

	if employeeID == nil {
		result, err = s.repo.ReportAll(ctx)
	} else {
		result, err = s.repo.ReportByID(ctx, *employeeID)
	}

	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "inboundorder-service", "Failed to generate inbound orders report", err, map[string]interface{}{
				"employee_id": func() interface{} {
					if employeeID != nil {
						return *employeeID
					} else {
						return nil
					}
				}(),
			})
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.Info(ctx, "inboundorder-service", "Inbound orders report generated successfully", map[string]interface{}{
			"employee_id": func() interface{} {
				if employeeID != nil {
					return *employeeID
				} else {
					return nil
				}
			}(),
		})
	}
	return result, nil
}
