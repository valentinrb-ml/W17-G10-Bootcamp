package service

import (
	"context"
	"errors"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

const prServiceName = "product-record-service" // [LOG]

type productRecordService struct {
	repo   repository.ProductRecordRepository
	logger logger.Logger
}

func NewProductRecordService(repo repository.ProductRecordRepository) ProductRecordService {
	return &productRecordService{repo: repo}
}

// SetLogger allows injecting the logger after creation
func (s *productRecordService) SetLogger(l logger.Logger) { // [LOG]
	s.logger = l // [LOG]
}

// logging helpers (avoid repeating nil-check and set the service name)
func (s *productRecordService) debug(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if s.logger != nil {
		s.logger.Debug(ctx, prServiceName, msg, md...) // [LOG]
	}
}
func (s *productRecordService) info(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if s.logger != nil {
		s.logger.Info(ctx, prServiceName, msg, md...) // [LOG]
	}
}

func (s *productRecordService) Create(ctx context.Context, record models.ProductRecord) (models.ProductRecordResponse, error) {
	s.debug(ctx, "Validating product record business rules", map[string]interface{}{ // [LOG]
		"product_id": record.ProductID, // [LOG]
	})

	if err := validators.ValidateProductRecordBusinessRules(record); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return models.ProductRecordResponse{}, err
		}
		return models.ProductRecordResponse{}, apperrors.NewAppError(apperrors.CodeBadRequest, err.Error())
	}

	s.info(ctx, "Creating product record", map[string]interface{}{ // [LOG]
		"product_id": record.ProductID, // [LOG]
	})

	savedRecord, err := s.repo.Create(ctx, record)
	if err != nil {
		return models.ProductRecordResponse{}, err
	}

	s.info(ctx, "Product record created", map[string]interface{}{ // [LOG]
		"product_record_id": savedRecord.ID, // [LOG]
		"product_id":        savedRecord.ProductID,
	})

	// No mapper needed, ProductRecordResponse is alias of ProductRecord
	return savedRecord, nil
}

func (s *productRecordService) GetRecordsReport(ctx context.Context, productID int) ([]models.ProductRecordReport, error) {
	if productID != 0 {
		s.info(ctx, "Getting product records report (filtered)", map[string]interface{}{ // [LOG]
			"product_id": productID, // [LOG]
		})
	} else {
		s.info(ctx, "Getting product records report") // [LOG]
	}

	reports, err := s.repo.GetRecordsReport(ctx, productID)
	if err != nil {
		return []models.ProductRecordReport{}, err
	}

	if productID != 0 {
		s.info(ctx, "Product records report generated (filtered)", map[string]interface{}{ // [LOG]
			"product_id": productID,    // [LOG]
			"count":      len(reports), // [LOG]
		})
	} else {
		s.info(ctx, "Product records report generated", map[string]interface{}{ // [LOG]
			"count": len(reports), // [LOG]
		})
	}

	return reports, nil
}
