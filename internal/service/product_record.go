package service

import (
	"context"
	"errors"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/productrecord"
)

type productRecordService struct {
	repo repository.ProductRecordRepository
}

func NewProductRecordService(repo repository.ProductRecordRepository) ProductRecordService {
	return &productRecordService{repo: repo}
}

func (s *productRecordService) Create(ctx context.Context, record productrecord.ProductRecord) (productrecord.ProductRecordResponse, error) {
	if err := validators.ValidateProductRecordBusinessRules(record); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return productrecord.ProductRecordResponse{}, err
		}
		return productrecord.ProductRecordResponse{}, apperrors.NewAppError(apperrors.CodeBadRequest, err.Error())
	}

	savedRecord, err := s.repo.Create(ctx, record)
	if err != nil {
		return productrecord.ProductRecordResponse{}, err
	}

	// No mapper needed, ProductRecordResponse is alias of ProductRecord
	return savedRecord, nil
}

func (s *productRecordService) GetRecordsReport(ctx context.Context, productID int) (productrecord.ProductRecordsReportResponse, error) {
	reports, err := s.repo.GetRecordsReport(ctx, productID)
	if err != nil {
		return productrecord.ProductRecordsReportResponse{}, err
	}

	return mappers.ProductRecordReportToResponse(reports), nil
}
