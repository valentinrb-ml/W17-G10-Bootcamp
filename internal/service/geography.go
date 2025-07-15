package service

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func (s *geographyService) Create(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error) {
	tx, err := s.rp.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			s.rp.RollbackTx(tx)
		}
	}()

	country, err := s.rp.FindCountryByName(ctx, tx, *gr.CountryName)
	if err != nil {
		if err == sql.ErrNoRows {
			newCountry := models.Country{Name: *gr.CountryName}
			country, err = s.rp.CreateCountry(ctx, tx, newCountry)
			if err != nil {
				return nil, fmt.Errorf("error creating country: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error finding country: %w", err)
		}
	}

	province, err := s.rp.FindProvinceByName(ctx, tx, *gr.ProvinceName, country.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			newProvince := models.Province{Name: *gr.ProvinceName, CountryId: country.Id}
			province, err = s.rp.CreateProvince(ctx, tx, newProvince)
			if err != nil {
				return nil, fmt.Errorf("error creating province: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error finding province: %w", err)
		}
	}

	locality, err := s.rp.FindLocalityByName(ctx, tx, *gr.LocalityName, province.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			newLocality := models.Locality{Name: *gr.LocalityName, ProvinceId: province.Id}
			locality, err = s.rp.CreateLocality(ctx, tx, newLocality)
			if err != nil {
				return nil, fmt.Errorf("error creating locality: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error finding locality: %w", err)
		}
	}

	if err = s.rp.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &models.ResponseGeography{
		LocalityId:   locality.Id,
		LocalityName: locality.Name,
		ProvinceName: province.Name,
		CountryName:  country.Name,
	}, nil
}
