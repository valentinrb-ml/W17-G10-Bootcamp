package service

import (
	"context"
	"time"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/purchase_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type purchaseOrderService struct {
	repo   repository.PurchaseOrderRepository
	logger logger.Logger
}

func NewPurchaseOrderService(repo repository.PurchaseOrderRepository) PurchaseOrderService {
	return &purchaseOrderService{repo: repo}
}

// SetLogger allows you to inject the logger after creation
func (s *purchaseOrderService) SetLogger(l logger.Logger) {
	s.logger = l
	s.repo.SetLogger(l) // También inyectar el logger al repository
}

func (s *purchaseOrderService) Create(ctx context.Context, req models.RequestPurchaseOrder) (*models.ResponsePurchaseOrder, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "purchase-order-service", "Creating purchase order", map[string]interface{}{
			"order_number":      req.OrderNumber,
			"buyer_id":          req.BuyerID,
			"product_record_id": req.ProductRecordID,
		})
	}

	po, err := mappers.RequestPurchaseOrderToPurchaseOrder(req)
	if err != nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "purchase-order-service", "Invalid date format in purchase order", map[string]interface{}{
				"order_date": req.OrderDate,
				"error":      err.Error(),
			})
		}
		return nil, apperrors.NewAppError(
			apperrors.CodeValidationError,
			"invalid date format, use YYYY-MM-DD",
		)
	}

	if po.OrderDate.After(time.Now()) {
		if s.logger != nil {
			s.logger.Warning(ctx, "purchase-order-service", "Future date in purchase order", map[string]interface{}{
				"order_date": po.OrderDate,
			})
		}
		return nil, apperrors.NewAppError(
			apperrors.CodeValidationError,
			"order date cannot be in the future",
		)
	}

	createdPO, err := s.repo.Create(ctx, po)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "purchase-order-service", "Failed to create purchase order", err, map[string]interface{}{
				"order_number": req.OrderNumber,
			})
		}
		return nil, err
	}

	response := mappers.PurchaseOrderToResponse(*createdPO)

	if s.logger != nil {
		s.logger.Info(ctx, "purchase-order-service", "Purchase order created successfully", map[string]interface{}{
			"purchase_order_id": response.ID,
			"order_number":      response.OrderNumber,
		})
	}

	return &response, nil
}

func (s *purchaseOrderService) GetAll(ctx context.Context) ([]models.ResponsePurchaseOrder, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "purchase-order-service", "Getting all purchase orders")
	}

	pos, err := s.repo.GetAll(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "purchase-order-service", "Failed to get all purchase orders", err, nil)
		}
		return nil, err
	}

	response := mappers.ToResponsePurchaseOrderList(pos)

	if s.logger != nil {
		s.logger.Info(ctx, "purchase-order-service", "Retrieved all purchase orders successfully", map[string]interface{}{
			"count": len(response),
		})
	}

	return response, nil
}

func (s *purchaseOrderService) GetByID(ctx context.Context, id int) (*models.ResponsePurchaseOrder, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "purchase-order-service", "Getting purchase order by ID", map[string]interface{}{
			"purchase_order_id": id,
		})
	}

	po, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "purchase-order-service", "Failed to get purchase order by ID", err, map[string]interface{}{
				"purchase_order_id": id,
			})
		}
		return nil, err
	}

	response := mappers.PurchaseOrderToResponse(*po)

	if s.logger != nil {
		s.logger.Info(ctx, "purchase-order-service", "Purchase order found successfully", map[string]interface{}{
			"purchase_order_id": response.ID,
			"order_number":      response.OrderNumber,
		})
	}

	return &response, nil
}

func (s *purchaseOrderService) GetReportByBuyer(ctx context.Context, buyerID *int) ([]models.BuyerWithPurchaseCount, error) {
	if s.logger != nil {
		if buyerID != nil {
			s.logger.Info(ctx, "purchase-order-service", "Getting purchase report by buyer", map[string]interface{}{
				"buyer_id": *buyerID,
			})
		} else {
			s.logger.Info(ctx, "purchase-order-service", "Getting purchase report for all buyers")
		}
	}

	var report []models.BuyerWithPurchaseCount
	var err error

	if buyerID != nil {
		report, err = s.repo.GetCountByBuyer(ctx, *buyerID)
	} else {
		report, err = s.repo.GetAllWithPurchaseCount(ctx)
	}

	if err != nil {
		if s.logger != nil {
			if buyerID != nil {
				s.logger.Error(ctx, "purchase-order-service", "Failed to get purchase report by buyer", err, map[string]interface{}{
					"buyer_id": *buyerID,
				})
			} else {
				s.logger.Error(ctx, "purchase-order-service", "Failed to get purchase report for all buyers", err, nil)
			}
		}
		return nil, err
	}

	if s.logger != nil {
		if buyerID != nil {
			s.logger.Info(ctx, "purchase-order-service", "Purchase report by buyer retrieved successfully", map[string]interface{}{
				"buyer_id": *buyerID,
				"count":    len(report),
			})
		} else {
			s.logger.Info(ctx, "purchase-order-service", "Purchase report for all buyers retrieved successfully", map[string]interface{}{
				"buyers_count": len(report),
			})
		}
	}

	return report, nil
}
