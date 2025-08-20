package service

import (
	"context"
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func (s *geographyService) Create(ctx context.Context, gr models.RequestGeography) (*models.ResponseGeography, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Creating geography", map[string]interface{}{
			"locality_id":   gr.Id,
			"locality_name": gr.LocalityName,
			"province_name": gr.ProvinceName,
			"country_name":  gr.CountryName,
		})
	}

	tx, err := s.rp.BeginTx(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "geography-service", "Failed to start transaction", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to start transaction").
			WithDetail("error", err.Error())
	}

	defer func() {
		if err != nil {
			if s.logger != nil {
				s.logger.Warning(ctx, "geography-service", "Rolling back transaction due to error", map[string]interface{}{
					"locality_id": gr.Id,
				})
			}
			s.rp.RollbackTx(tx)
		}
	}()

	country, err := s.handleCountry(ctx, tx, *gr.CountryName)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "geography-service", "Failed to handle country", err, map[string]interface{}{
				"country_name": gr.CountryName,
			})
		}
		return nil, err
	}

	province, err := s.handleProvince(ctx, tx, *gr.ProvinceName, country.Id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "geography-service", "Failed to handle province", err, map[string]interface{}{
				"province_name": gr.ProvinceName,
				"country_id":    country.Id,
			})
		}
		return nil, err
	}

	locality, err := s.handleLocality(ctx, tx, *gr.Id, *gr.LocalityName, province.Id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "geography-service", "Failed to handle locality", err, map[string]interface{}{
				"locality_id":   gr.Id,
				"locality_name": gr.LocalityName,
				"province_id":   province.Id,
			})
		}
		return nil, err
	}

	if err = s.rp.CommitTx(tx); err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "geography-service", "Failed to commit transaction", err, map[string]interface{}{
				"locality_id": gr.Id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to commit transaction").
			WithDetail("error", err.Error())
	}

	response := &models.ResponseGeography{
		LocalityId:   locality.Id,
		LocalityName: locality.Name,
		ProvinceName: province.Name,
		CountryName:  country.Name,
	}

	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Geography created successfully", map[string]interface{}{
			"locality_id":   response.LocalityId,
			"locality_name": response.LocalityName,
			"province_name": response.ProvinceName,
			"country_name":  response.CountryName,
		})
	}

	return response, nil
}

func (s *geographyService) CountSellersGroupedByLocality(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Counting sellers grouped by locality")
	}

	result, err := s.rp.CountSellersGroupedByLocality(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "geography-service", "Failed to count sellers grouped by locality", err, nil)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Sellers grouped by locality counted successfully", map[string]interface{}{
			"localities_count": len(result),
		})
	}

	return result, nil
}

func (s *geographyService) CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Counting sellers by locality", map[string]interface{}{
			"locality_id": id,
		})
	}

	resp, err := s.rp.CountSellersByLocality(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(ctx, "geography-service", "Failed to count sellers by locality", err, map[string]interface{}{
				"locality_id": id,
			})
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Sellers by locality counted successfully", map[string]interface{}{
			"locality_id":   resp.LocalityId,
			"locality_name": resp.LocalityName,
			"sellers_count": resp.SellersCount,
		})
	}

	return resp, nil
}

func (s *geographyService) handleCountry(ctx context.Context, tx *sql.Tx, countryName string) (*models.Country, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Handling country", map[string]interface{}{
			"country_name": countryName,
		})
	}

	country, err := s.rp.FindCountryByName(ctx, countryName)
	if err != nil {
		if apperrors.IsAppError(err, apperrors.CodeNotFound) {
			if s.logger != nil {
				s.logger.Info(ctx, "geography-service", "Country not found, creating new one", map[string]interface{}{
					"country_name": countryName,
				})
			}

			newCountry := models.Country{Name: countryName}
			country, err = s.rp.CreateCountry(ctx, tx, newCountry)
			if err != nil {
				if s.logger != nil {
					s.logger.Error(ctx, "geography-service", "Failed to create country", err, map[string]interface{}{
						"country_name": countryName,
					})
				}
				return nil, err
			}

			if s.logger != nil {
				s.logger.Info(ctx, "geography-service", "Country created successfully", map[string]interface{}{
					"country_id":   country.Id,
					"country_name": country.Name,
				})
			}
		} else {
			if s.logger != nil {
				s.logger.Error(ctx, "geography-service", "Failed to find country", err, map[string]interface{}{
					"country_name": countryName,
				})
			}
			return nil, err
		}
	} else {
		if s.logger != nil {
			s.logger.Info(ctx, "geography-service", "Country found", map[string]interface{}{
				"country_id":   country.Id,
				"country_name": country.Name,
			})
		}
	}

	return country, nil
}

func (s *geographyService) handleProvince(ctx context.Context, tx *sql.Tx, provinceName string, countryId int) (*models.Province, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Handling province", map[string]interface{}{
			"province_name": provinceName,
			"country_id":    countryId,
		})
	}

	province, err := s.rp.FindProvinceByName(ctx, provinceName, countryId)
	if err != nil {
		if apperrors.IsAppError(err, apperrors.CodeNotFound) {
			if s.logger != nil {
				s.logger.Info(ctx, "geography-service", "Province not found, creating new one", map[string]interface{}{
					"province_name": provinceName,
					"country_id":    countryId,
				})
			}

			newProvince := models.Province{Name: provinceName, CountryId: countryId}
			province, err = s.rp.CreateProvince(ctx, tx, newProvince)
			if err != nil {
				if s.logger != nil {
					s.logger.Error(ctx, "geography-service", "Failed to create province", err, map[string]interface{}{
						"province_name": provinceName,
						"country_id":    countryId,
					})
				}
				return nil, err
			}

			if s.logger != nil {
				s.logger.Info(ctx, "geography-service", "Province created successfully", map[string]interface{}{
					"province_id":   province.Id,
					"province_name": province.Name,
					"country_id":    province.CountryId,
				})
			}
		} else {
			if s.logger != nil {
				s.logger.Error(ctx, "geography-service", "Failed to find province", err, map[string]interface{}{
					"province_name": provinceName,
					"country_id":    countryId,
				})
			}
			return nil, err
		}
	} else {
		if s.logger != nil {
			s.logger.Info(ctx, "geography-service", "Province found", map[string]interface{}{
				"province_id":   province.Id,
				"province_name": province.Name,
			})
		}
	}

	return province, nil
}

func (s *geographyService) handleLocality(ctx context.Context, tx *sql.Tx, localityId string, localityName string, provinceId int) (*models.Locality, error) {
	if s.logger != nil {
		s.logger.Info(ctx, "geography-service", "Handling locality", map[string]interface{}{
			"locality_id":   localityId,
			"locality_name": localityName,
			"province_id":   provinceId,
		})
	}

	_, err := s.rp.FindLocalityById(ctx, localityId)

	if err == nil {
		if s.logger != nil {
			s.logger.Warning(ctx, "geography-service", "Locality already exists", map[string]interface{}{
				"locality_id": localityId,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "locality already exists")
	}

	if apperrors.IsAppError(err, apperrors.CodeNotFound) {
		if s.logger != nil {
			s.logger.Info(ctx, "geography-service", "Locality not found, creating new one", map[string]interface{}{
				"locality_id":   localityId,
				"locality_name": localityName,
				"province_id":   provinceId,
			})
		}

		newLocality := models.Locality{Id: localityId, Name: localityName, ProvinceId: provinceId}
		locality, err := s.rp.CreateLocality(ctx, tx, newLocality)
		if err != nil {
			if s.logger != nil {
				s.logger.Error(ctx, "geography-service", "Failed to create locality", err, map[string]interface{}{
					"locality_id":   localityId,
					"locality_name": localityName,
					"province_id":   provinceId,
				})
			}
			return nil, err
		}

		if s.logger != nil {
			s.logger.Info(ctx, "geography-service", "Locality created successfully", map[string]interface{}{
				"locality_id":   locality.Id,
				"locality_name": locality.Name,
				"province_id":   locality.ProvinceId,
			})
		}

		return locality, nil
	}

	if s.logger != nil {
		s.logger.Error(ctx, "geography-service", "Failed to find locality", err, map[string]interface{}{
			"locality_id": localityId,
		})
	}

	return nil, err
}
