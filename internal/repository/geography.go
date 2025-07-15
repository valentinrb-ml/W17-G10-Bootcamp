package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

const (
	queryCountryCreate       = `INSERT INTO countries (name) VALUES (?)`
	queryCountryFindById     = `SELECT id, name FROM countries WHERE LOWER(name) = LOWER(?)`
	queryProvinceCreate      = `INSERT INTO provinces (name, country_id) VALUES (?, ?)`
	queryProvinceFindById    = `SELECT id, name, country_id FROM provinces WHERE LOWER(name) = LOWER(?) AND country_id = ?`
	queryLocalityCreate      = `INSERT INTO localities (id, name, province_id) VALUES (?, ?, ?)`
	queryLocalityFindById    = `SELECT id, name, province_id FROM localities WHERE id = ?`
	queryLocalityWithSellers = `SELECT l.id , l.name, COUNT(s.id) FROM localities l
								LEFT JOIN sellers s ON l.id = s.locality_id
								WHERE l.id = ?
								GROUP BY l.id, l.name`
)

func (r *geographyRepository) CreateCountry(ctx context.Context, exec Executor, c models.Country) (*models.Country, error) {
	res, err := exec.ExecContext(ctx, queryCountryCreate, c.Name)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to create country").WithDetail("error", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to get country ID after creation").WithDetail("error", err.Error())
	}

	c.Id = int(id)
	return &c, nil
}

func (r *geographyRepository) FindCountryByName(ctx context.Context, exec Executor, name string) (*models.Country, error) {
	var country models.Country
	err := exec.QueryRowContext(ctx, queryCountryFindById, name).Scan(&country.Id, &country.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "country not found")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to find country").WithDetail("error", err.Error())
	}
	return &country, nil
}

func (r *geographyRepository) CreateProvince(ctx context.Context, exec Executor, p models.Province) (*models.Province, error) {
	res, err := exec.ExecContext(ctx, queryProvinceCreate, p.Name, p.CountryId)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to create province").WithDetail("error", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to get province ID after creation").WithDetail("error", err.Error())
	}

	p.Id = int(id)
	return &p, nil
}

func (r *geographyRepository) FindProvinceByName(ctx context.Context, exec Executor, name string, countryId int) (*models.Province, error) {
	var province models.Province
	err := exec.QueryRowContext(ctx, queryProvinceFindById, name, countryId).Scan(&province.Id, &province.Name, &province.CountryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "province not found")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to find province").WithDetail("error", err.Error())
	}
	return &province, nil
}

func (r *geographyRepository) CreateLocality(ctx context.Context, exec Executor, l models.Locality) (*models.Locality, error) {
	_, err := exec.ExecContext(ctx, queryLocalityCreate, l.Id, l.Name, l.ProvinceId)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "postal code already exists").
				WithDetail("postal_code", l.Id).
				WithDetail("locality_name", l.Name)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to create locality").WithDetail("error", err.Error())
	}

	return &l, nil
}

func (r *geographyRepository) FindLocalityById(ctx context.Context, exec Executor, id string) (*models.Locality, error) {
	var locality models.Locality
	err := exec.QueryRowContext(ctx, queryLocalityFindById, id).Scan(&locality.Id, &locality.Name, &locality.ProvinceId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "locality not found")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to find locality").WithDetail("error", err.Error())
	}
	return &locality, nil
}

func (r *geographyRepository) CountSellersByLocality(ctx context.Context, id string) (*models.ResponseLocalitySellers, error) {
	var resp models.ResponseLocalitySellers
	row := r.mysql.QueryRowContext(ctx, queryLocalityWithSellers, id)
	err := row.Scan(&resp.LocalityId, &resp.LocalityName, &resp.SellersCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The locality you are looking for does not exist.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the locality sellers count.")
	}
	return &resp, nil
}
