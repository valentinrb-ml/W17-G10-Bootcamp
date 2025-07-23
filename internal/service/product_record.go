package service

import (
	"context"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

type productRecordService struct {
	repo repository.ProductRecordRepository
}

func NewProductRecordService(repo repository.ProductRecordRepository) ProductRecordService {
	return &productRecordService{repo: repo}
}

func (s *productRecordService) Create(ctx context.Context, record models.ProductRecord) (models.ProductRecordResponse, error) {
	if err := validators.ValidateProductRecordBusinessRules(record); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return models.ProductRecordResponse{}, err
		}
		return models.ProductRecordResponse{}, apperrors.NewAppError(apperrors.CodeBadRequest, err.Error())
	}

	savedRecord, err := s.repo.Create(ctx, record)
	if err != nil {
		return models.ProductRecordResponse{}, err
	}

	// No mapper needed, ProductRecordResponse is alias of ProductRecord
	return savedRecord, nil
}

func (s *productRecordService) GetRecordsReport(ctx context.Context, productID int) (models.ProductRecordsReportResponse, error) {
	reports, err := s.repo.GetRecordsReport(ctx, productID)
	if err != nil {
		return models.ProductRecordsReportResponse{}, err
	}

	return mappers.ProductRecordReportToResponse(reports), nil
}
