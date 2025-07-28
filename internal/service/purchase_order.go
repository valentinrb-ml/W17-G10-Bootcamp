package service

import (
	"context"
	"time"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type purchaseOrderService struct {
	repo repository.PurchaseOrderRepository
}

func NewPurchaseOrderService(repo repository.PurchaseOrderRepository) PurchaseOrderService {
	return &purchaseOrderService{repo: repo}
}

func (s *purchaseOrderService) Create(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
	po, err := mappers.RequestPurchaseOrderToPurchaseOrder(req)
	if err != nil {
		return nil, apperrors.NewAppError(
			apperrors.CodeValidationError,
			"invalid date format, use YYYY-MM-DD",
		)
	}

	if po.OrderDate.After(time.Now()) {
		return nil, apperrors.NewAppError(
			apperrors.CodeValidationError,
			"order date cannot be in the future",
		)
	}

	createdPO, err := s.repo.Create(ctx, po)
	if err != nil {
		return nil, err
	}

	response := mappers.PurchaseOrderToResponse(*createdPO)
	return &response, nil
}

func (s *purchaseOrderService) GetAll(ctx context.Context) ([]models.ResponsePurchaseOrder, error) {
	pos, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return mappers.ToResponsePurchaseOrderList(pos), nil
}

func (s *purchaseOrderService) GetByID(ctx context.Context, id int) (*models.ResponsePurchaseOrder, error) {
	po, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mappers.PurchaseOrderToResponse(*po)
	return &response, nil
}

func (s *purchaseOrderService) GetReportByBuyer(ctx context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) {
	if buyerID != nil {
		report, err := s.repo.GetCountByBuyer(ctx, *buyerID)
		if err != nil {
			return nil, err
		}
		return report, nil
	}

	report, err := s.repo.GetAllWithPurchaseCount(ctx)
	if err != nil {
		return nil, err
	}
	return report, nil
}
