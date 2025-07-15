package repository

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

const (
	queryCarryCreate = `INSERT INTO carriers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)`
)

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
