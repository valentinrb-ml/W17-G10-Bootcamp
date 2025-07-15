package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
)

type InboundOrderDefault struct {
	repo         repository.InboundOrderRepository
	employeeRepo repository.EmployeeRepository
}

func NewInboundOrderService(r repository.InboundOrderRepository, er repository.EmployeeRepository) *InboundOrderDefault {
	return &InboundOrderDefault{
		repo:         r,
		employeeRepo: er,
	}
}

func (s *InboundOrderDefault) Create(ctx context.Context, o *models.InboundOrder) (*models.InboundOrder, error) {
	exists, err := s.repo.ExistsByOrderNumber(ctx, o.OrderNumber)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed checking for order_number uniqueness")
	}
	if exists {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "order_number already exists")
	}
	emp, err := s.employeeRepo.FindByID(ctx, o.EmployeeID)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed getting employee by id")
	}
	if emp == nil {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "employee_id does not exist")
	}
	return s.repo.Create(ctx, o)
}

func (s *InboundOrderDefault) Report(ctx context.Context, employeeID *int) (interface{}, error) {
	if employeeID == nil {
		return s.repo.ReportAll(ctx)
	}
	return s.repo.ReportByID(ctx, *employeeID)
}
