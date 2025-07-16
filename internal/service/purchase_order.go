package service

import (
	"context"
	"net/http"
	"time"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type purchaseOrderService struct {
	repo repository.PurchaseOrderRepository
}

func NewPurchaseOrderService(repo repository.PurchaseOrderRepository) PurchaseOrderService {
	return &purchaseOrderService{repo: repo}
}

func (s *purchaseOrderService) Create(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
	// Convertir Request a PurchaseOrder
	po, err := mappers.RequestPurchaseOrderToPurchaseOrder(req)
	if err != nil {
		return nil, &api.ServiceError{
			Code:         http.StatusBadRequest,
			ResponseCode: http.StatusBadRequest,
			Message:      "invalid date format, use YYYY-MM-DD",
		}
	}

	// Validar fecha no sea futura
	if po.OrderDate.After(time.Now()) {
		return nil, &api.ServiceError{
			Code:         http.StatusBadRequest,
			ResponseCode: http.StatusBadRequest,
			Message:      "order date cannot be in the future",
		}
	}

	// Crear en repositorio
	createdPO, err := s.repo.Create(ctx, po)
	if err != nil {
		return nil, err
	}

	// Convertir a Response
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
		// Reporte para un buyer específico
		return s.repo.GetCountByBuyer(ctx, *buyerID)
	}
	// Reporte general de todos los buyers
	return s.repo.GetAllWithPurchaseCount(ctx)
}
