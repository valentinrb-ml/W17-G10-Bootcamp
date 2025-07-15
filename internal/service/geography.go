package service

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func (s *geographyService) Create(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error) {
	tx, err := s.rp.BeginTx(ctx)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to start transaction").
			WithDetail("error", err.Error())
	}

	defer func() {
		if err != nil {
			s.rp.RollbackTx(tx)
		}
	}()

	country, err := s.handleCountry(ctx, tx, *gr.CountryName)
	if err != nil {
		return nil, err
	}

	province, err := s.handleProvince(ctx, tx, *gr.ProvinceName, country.Id)
	if err != nil {
		return nil, err
	}

	locality, err := s.handleLocality(ctx, tx, *gr.Id, *gr.LocalityName, province.Id)
	if err != nil {
		return nil, err
	}

	if err = s.rp.CommitTx(tx); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to commit transaction").
			WithDetail("error", err.Error())
	}

	return &models.ResponseGeography{
		LocalityId:   locality.Id,
		LocalityName: locality.Name,
		ProvinceName: province.Name,
		CountryName:  country.Name,
	}, nil
}

func (s *geographyService) CountSellersGroupedByLocality(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
	return s.rp.CountSellersGroupedByLocality(ctx)
}

func (s *geographyService) CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
	resp, err := s.rp.CountSellersByLocality(ctx, id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *geographyService) handleCountry(ctx context.Context, tx *sql.Tx, countryName string) (*models.Country, error) {
	country, err := s.rp.FindCountryByName(ctx, tx, countryName)
	if err != nil {
		if apperrors.IsAppError(err, apperrors.CodeNotFound) {
			newCountry := models.Country{Name: countryName}
			country, err = s.rp.CreateCountry(ctx, tx, newCountry)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return country, nil
}

func (s *geographyService) handleProvince(ctx context.Context, tx *sql.Tx, provinceName string, countryId int) (*models.Province, error) {
	province, err := s.rp.FindProvinceByName(ctx, tx, provinceName, countryId)
	if err != nil {
		if apperrors.IsAppError(err, apperrors.CodeNotFound) {
			newProvince := models.Province{Name: provinceName, CountryId: countryId}
			province, err = s.rp.CreateProvince(ctx, tx, newProvince)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return province, nil
}

func (s *geographyService) handleLocality(ctx context.Context, tx *sql.Tx, localityId string, localityName string, provinceId int) (*models.Locality, error) {
	_, err := s.rp.FindLocalityById(ctx, tx, localityId)

	if err == nil {
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "locality already exists")
	}

	if apperrors.IsAppError(err, apperrors.CodeNotFound) {
		newLocality := models.Locality{Id: localityId, Name: localityName, ProvinceId: provinceId}
		locality, err := s.rp.CreateLocality(ctx, tx, newLocality)
		if err != nil {
			return nil, err
		}
		return locality, nil
	}

	return nil, err
}
