package service

import (
	"context"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	wRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

// Servicio principal para Inbound Orders, incluye lógica de validación
type InboundOrderDefault struct {
	repo          repository.InboundOrderRepository
	employeeRepo  repository.EmployeeRepository
	warehouseRepo wRepo.WarehouseRepository
}

// Constructor del servicio de inbound orders
func NewInboundOrderService(r repository.InboundOrderRepository, er repository.EmployeeRepository, wr wRepo.WarehouseRepository) *InboundOrderDefault {
	return &InboundOrderDefault{
		repo:          r,
		employeeRepo:  er,
		warehouseRepo: wr,
	}
}

// Crea un nuevo inbound order con todas las validaciones de negocio y de integridad
func (s *InboundOrderDefault) Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
	// Validación de campos obligatorios (devuelve 422 si falta alguno)
	if err := validators.ValidateInboundOrder(o); err != nil {
		return nil, err
	}
	// Verifica unicidad del order_number
	exists, err := s.repo.ExistsByOrderNumber(ctx, o.OrderNumber)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed checking for order_number uniqueness")
	}
	if exists {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "order_number already exists")
	}
	// Valida que el empleado referenciado exista (FK)
	emp, err := s.employeeRepo.FindByID(ctx, o.EmployeeID)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed getting employee by id")
	}
	if emp == nil {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "employee_id does not exist")
	}
	// Valida que el warehouse referenciado exista (FK)
	warehouse, whErr := s.warehouseRepo.FindById(ctx, o.WarehouseID)
	if whErr != nil {
		var appErr *apperrors.AppError
		if errors.As(whErr, &appErr) && appErr.Code == apperrors.CodeNotFound {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_id does not exist")
		}
		return nil, apperrors.Wrap(whErr, "failed getting warehouse by id")
	}
	if warehouse == nil {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_id does not exist")
	}
	// Si pasa todas las validaciones, crea el inbound order
	return s.repo.Create(ctx, o)
}

// Genera un reporte de inbound orders por empleado o global si employeeID es nil
func (s *InboundOrderDefault) Report(ctx context.Context, employeeID *int) (interface{}, error) {
	if employeeID == nil {
		// Retorna el reporte general para todos los empleados
		return s.repo.ReportAll(ctx)
	}
	// Retorna reporte solo para el empleado solicitado
	return s.repo.ReportByID(ctx, *employeeID)
}
