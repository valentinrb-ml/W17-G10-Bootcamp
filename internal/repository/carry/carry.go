package repository

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

// SQL queries for carrier operations
const (
	queryCarryCreate                 = `INSERT INTO carriers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)`
	queryCarriesCountByAllLocalities = `SELECT c.locality_id, l.name, COUNT(*) as carries_count FROM carriers c INNER JOIN localities l ON c.locality_id = l.id GROUP BY c.locality_id`
	queryCarriesCountByLocalityID    = `SELECT c.locality_id, l.name, COUNT(*) as carries_count FROM carriers c INNER JOIN localities l ON c.locality_id = l.id WHERE c.locality_id = ? GROUP BY c.locality_id`
)

// Create inserts a new carrier into the database
// Handles specific MySQL errors for duplicate CID and invalid locality_id
// Returns the created carrier with its generated ID or an error if the operation fails
func (r *CarryMySQL) Create(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
	r.logger.Info(ctx, "carry-repository", "Starting carry creation in database", map[string]interface{}{
		"company_name": c.CompanyName,
		"cid":          c.Cid,
		"locality_id":  c.LocalityId,
	})

	res, err := r.db.ExecContext(ctx, queryCarryCreate, c.Cid, c.CompanyName, c.Address, c.Telephone, c.LocalityId)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				r.logger.Warning(ctx, "carry-repository", "Carry creation failed - CID already exists", map[string]interface{}{
					"cid":         c.Cid,
					"mysql_error": mysqlErr.Number,
				})
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "cid already exists")
			case 1452:
				r.logger.Warning(ctx, "carry-repository", "Carry creation failed - locality_id does not exist", map[string]interface{}{
					"locality_id": c.LocalityId,
					"mysql_error": mysqlErr.Number,
				})
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "locality_id does not exist")
			}
		}
		r.logger.Error(ctx, "carry-repository", "Database error during carry creation", err, map[string]interface{}{
			"company_name": c.CompanyName,
			"cid":          c.Cid,
		})
		return nil, apperrors.Wrap(err, "error creating carry")
	}

	id, err := res.LastInsertId()
	if err != nil {
		r.logger.Error(ctx, "carry-repository", "Failed to get last insert ID", err, map[string]interface{}{
			"company_name": c.CompanyName,
			"cid":          c.Cid,
		})
		return nil, apperrors.Wrap(err, "error creating carry")
	}

	c.Id = int(id)
	r.logger.Info(ctx, "carry-repository", "Carry created successfully in database", map[string]interface{}{
		"carry_id":     c.Id,
		"company_name": c.CompanyName,
		"cid":          c.Cid,
	})
	return &c, nil
}

// GetCarriesCountByAllLocalities retrieves the count of carriers grouped by all localities
// Joins carriers with localities table to get locality names along with counts
// Returns a slice of CarriesReport containing locality information and carrier counts
func (r *CarryMySQL) GetCarriesCountByAllLocalities(ctx context.Context) ([]carry.CarriesReport, error) {
	r.logger.Info(ctx, "carry-repository", "Starting query for carries count by all localities")

	rows, err := r.db.QueryContext(ctx, queryCarriesCountByAllLocalities)
	if err != nil {
		r.logger.Error(ctx, "carry-repository", "Database query failed for carries count by all localities", err)
		return nil, err
	}
	defer rows.Close()

	var results []carry.CarriesReport
	for rows.Next() {
		var cc carry.CarriesReport
		if err := rows.Scan(&cc.LocalityID, &cc.LocalityName, &cc.CarriesCount); err != nil {
			r.logger.Warning(ctx, "carry-repository", "Failed to scan row in carries count query", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		results = append(results, cc)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, "carry-repository", "Row iteration error in carries count query", err)
		return nil, apperrors.Wrap(err, "error getting carries count by all localities")
	}

	r.logger.Info(ctx, "carry-repository", "Successfully retrieved carries count by all localities", map[string]interface{}{
		"results_count": len(results),
	})
	return results, nil
}

// GetCarriesCountByLocalityID retrieves the count of carriers for a specific locality
// Joins carriers with localities table to get locality name along with count
// Returns a CarriesReport containing locality information and carrier count for the specified locality
func (r *CarryMySQL) GetCarriesCountByLocalityID(ctx context.Context, localityID string) (*carry.CarriesReport, error) {
	r.logger.Debug(ctx, "carry-repository", "Starting query for carries count by specific locality", map[string]interface{}{
		"locality_id": localityID,
	})

	var cc carry.CarriesReport
	err := r.db.QueryRowContext(ctx, queryCarriesCountByLocalityID, localityID).Scan(&cc.LocalityID, &cc.LocalityName, &cc.CarriesCount)
	if err != nil {
		r.logger.Error(ctx, "carry-repository", "Database query failed for carries count by locality ID", err, map[string]interface{}{
			"locality_id": localityID,
		})
		return nil, apperrors.Wrap(err, "error getting carries count by locality id")
	}

	r.logger.Info(ctx, "carry-repository", "Successfully retrieved carries count by locality ID", map[string]interface{}{
		"locality_id":   localityID,
		"locality_name": cc.LocalityName,
		"carries_count": cc.CarriesCount,
	})
	return &cc, nil
}
