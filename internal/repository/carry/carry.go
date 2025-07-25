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
	res, err := r.db.ExecContext(ctx, queryCarryCreate, c.Cid, c.CompanyName, c.Address, c.Telephone, c.LocalityId)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "cid already exists")
			case 1452:
				return nil, apperrors.NewAppError(apperrors.CodeConflict, "locality_id does not exist")
			}
		}
		return nil, apperrors.Wrap(err, "error creating carry")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperrors.Wrap(err, "error creating carry")
	}
	c.Id = int(id)
	return &c, nil
}

// GetCarriesCountByAllLocalities retrieves the count of carriers grouped by all localities
// Joins carriers with localities table to get locality names along with counts
// Returns a slice of CarriesReport containing locality information and carrier counts
func (r *CarryMySQL) GetCarriesCountByAllLocalities(ctx context.Context) ([]carry.CarriesReport, error) {
	rows, err := r.db.QueryContext(ctx, queryCarriesCountByAllLocalities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []carry.CarriesReport
	for rows.Next() {
		var cc carry.CarriesReport
		if err := rows.Scan(&cc.LocalityID, &cc.LocalityName, &cc.CarriesCount); err != nil {
			continue
		}
		results = append(results, cc)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, "error getting carries count by all localities")
	}
	return results, nil
}

// GetCarriesCountByLocalityID retrieves the count of carriers for a specific locality
// Joins carriers with localities table to get locality name along with count
// Returns a CarriesReport containing locality information and carrier count for the specified locality
func (r *CarryMySQL) GetCarriesCountByLocalityID(ctx context.Context, localityID string) (*carry.CarriesReport, error) {
	var cc carry.CarriesReport
	err := r.db.QueryRowContext(ctx, queryCarriesCountByLocalityID, localityID).Scan(&cc.LocalityID, &cc.LocalityName, &cc.CarriesCount)
	if err != nil {
		return nil, apperrors.Wrap(err, "error getting carries count by locality id")
	}
	return &cc, err
}
