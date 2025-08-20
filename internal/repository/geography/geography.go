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
	queryLocalityWithSellers = `SELECT l.id, l.name, COUNT(s.id) FROM localities l
								LEFT JOIN sellers s ON l.id = s.locality_id
								WHERE l.id = ?
								GROUP BY l.id, l.name`
	queryAllLocalitiesWithSellers = `SELECT l.id, l.name, COUNT(s.id) FROM localities l
									 LEFT JOIN sellers s ON l.id = s.locality_id
									 GROUP BY l.id, l.name`
)

func (r *geographyRepository) CreateCountry(ctx context.Context, exec Executor, c models.Country) (*models.Country, error) {
	res, err := exec.ExecContext(ctx, queryCountryCreate, c.Name)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Failed to create country", err, map[string]interface{}{
				"country_name": c.Name,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to create country").WithDetail("error", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Failed to get country ID after creation", err, map[string]interface{}{
				"country_name": c.Name,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to get country ID after creation").WithDetail("error", err.Error())
	}

	c.Id = int(id)
	return &c, nil
}

func (r *geographyRepository) FindCountryByName(ctx context.Context, name string) (*models.Country, error) {
	var country models.Country
	err := r.mysql.QueryRowContext(ctx, queryCountryFindById, name).Scan(&country.Id, &country.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			// No loggear este caso ya que es comportamiento esperado
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "country not found")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Database error finding country by name", err, map[string]interface{}{
				"country_name": name,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to find country").WithDetail("error", err.Error())
	}
	return &country, nil
}

func (r *geographyRepository) CreateProvince(ctx context.Context, exec Executor, p models.Province) (*models.Province, error) {
	res, err := exec.ExecContext(ctx, queryProvinceCreate, p.Name, p.CountryId)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Failed to create province", err, map[string]interface{}{
				"province_name": p.Name,
				"country_id":    p.CountryId,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to create province").WithDetail("error", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Failed to get province ID after creation", err, map[string]interface{}{
				"province_name": p.Name,
				"country_id":    p.CountryId,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to get province ID after creation").WithDetail("error", err.Error())
	}

	p.Id = int(id)
	return &p, nil
}

func (r *geographyRepository) FindProvinceByName(ctx context.Context, name string, countryId int) (*models.Province, error) {
	var province models.Province
	err := r.mysql.QueryRowContext(ctx, queryProvinceFindById, name, countryId).Scan(&province.Id, &province.Name, &province.CountryId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No loggear este caso ya que es comportamiento esperado
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "province not found")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Database error finding province by name", err, map[string]interface{}{
				"province_name": name,
				"country_id":    countryId,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to find province").WithDetail("error", err.Error())
	}
	return &province, nil
}

func (r *geographyRepository) CreateLocality(ctx context.Context, exec Executor, l models.Locality) (*models.Locality, error) {
	_, err := exec.ExecContext(ctx, queryLocalityCreate, l.Id, l.Name, l.ProvinceId)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			if r.logger != nil {
				r.logger.Warning(ctx, "geography-repository", "Duplicate locality creation attempted", map[string]interface{}{
					"locality_id":   l.Id,
					"locality_name": l.Name,
					"province_id":   l.ProvinceId,
					"mysql_error":   mysqlErr.Number,
				})
			}
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "The locality you are creating already exists.").
				WithDetail("postal_code", l.Id).
				WithDetail("locality_name", l.Name)
		}
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Failed to create locality", err, map[string]interface{}{
				"locality_id":   l.Id,
				"locality_name": l.Name,
				"province_id":   l.ProvinceId,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "failed to create locality").WithDetail("error", err.Error())
	}

	return &l, nil
}

func (r *geographyRepository) FindLocalityById(ctx context.Context, id string) (*models.Locality, error) {
	var locality models.Locality
	err := r.mysql.QueryRowContext(ctx, queryLocalityFindById, id).Scan(&locality.Id, &locality.Name, &locality.ProvinceId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No loggear este caso ya que es comportamiento esperado
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The locality you are looking for does not exist.")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Database error finding locality by ID", err, map[string]interface{}{
				"locality_id": id,
			})
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
			// No loggear este caso ya que es comportamiento esperado
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The locality you are looking for does not exist.")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Database error counting sellers by locality", err, map[string]interface{}{
				"locality_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the locality sellers count.")
	}
	return &resp, nil
}

func (r *geographyRepository) CountSellersGroupedByLocality(ctx context.Context) ([]models.ResponseLocalitySellers, error) {
	rows, err := r.mysql.QueryContext(ctx, queryAllLocalitiesWithSellers)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Database error querying sellers grouped by locality", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the sellers count by locality.")
	}
	defer rows.Close()

	var results []models.ResponseLocalitySellers
	for rows.Next() {
		var resp models.ResponseLocalitySellers
		if err := rows.Scan(&resp.LocalityId, &resp.LocalityName, &resp.SellersCount); err != nil {
			if r.logger != nil {
				r.logger.Error(ctx, "geography-repository", "Error scanning locality sellers count row", err, nil)
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "Failed to scan locality sellers count.")
		}
		results = append(results, resp)
	}
	if err := rows.Err(); err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "geography-repository", "Error iterating over locality sellers count rows", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An error occurred while iterating over the localities.")
	}

	return results, nil
}

func (r *geographyRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.mysql.BeginTx(ctx, nil)
}

func (r *geographyRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *geographyRepository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *geographyRepository) GetDB() *sql.DB {
	return r.mysql
}
