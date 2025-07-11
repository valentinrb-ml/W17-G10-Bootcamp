package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

const (
	querySellerCreate    = `INSERT INTO sellers (cid, company_name, address, telephone) VALUES (?, ?, ?, ?)`
	querySellerUpdate    = `UPDATE sellers SET cid = ?, company_name = ?, address = ?, telephone = ? WHERE id = ?`
	querySellerDelete    = `DELETE FROM sellers WHERE id = ?`
	querySellerFindAll   = `SELECT id, cid, company_name, address, telephone FROM sellers`
	querySellerFindById  = `SELECT id, cid, company_name, address, telephone FROM sellers WHERE id = ?`
	querySellerCIDExists = `SELECT EXISTS(SELECT 1 FROM sellers	WHERE LOWER(cid) = LOWER(?) AND id != ?)`
)

func (r *sellerRepository) Create(ctx context.Context, s models.Seller) (*models.Seller, error) {
	res, err := r.mysql.ExecContext(
		ctx,
		querySellerCreate,
		s.Cid, s.CompanyName, s.Address, s.Telephone,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	s.Id = int(id)

	return &s, nil
}

func (r *sellerRepository) Update(ctx context.Context, id int, s models.Seller) error {
	_, err := r.mysql.ExecContext(
		ctx,
		querySellerUpdate,
		s.Cid, s.CompanyName, s.Address, s.Telephone, s.Id,
	)

	return err
}

func (r *sellerRepository) Delete(ctx context.Context, id int) error {
	result, err := r.mysql.ExecContext(ctx, querySellerDelete, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
			errDef := api.ServiceErrors[api.ErrConflict]
			return &api.ServiceError{
				Code:         errDef.Code,
				ResponseCode: errDef.ResponseCode,
				Message:      "Cannot delete seller: there are products associated with this seller.",
			}
		}
		errDef := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      fmt.Sprintf("An internal server error occurred while deleting the seller: %s", err.Error()),
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errDef := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      fmt.Sprintf("An internal server error occurred while deleting the seller: %s", err.Error()),
		}
	}
	if rowsAffected == 0 {
		errDef := api.ServiceErrors[api.ErrNotFound]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "The seller you are trying to delete does not exist",
		}
	}
	return nil
}

func (r *sellerRepository) FindAll(ctx context.Context) ([]models.Seller, error) {
	rows, err := r.mysql.QueryContext(ctx, querySellerFindAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sellers []models.Seller
	for rows.Next() {
		var s models.Seller
		err := rows.Scan(&s.Id, &s.Cid, &s.CompanyName, &s.Address, &s.Telephone)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sellers, nil
}

func (r *sellerRepository) FindById(ctx context.Context, id int) (*models.Seller, error) {
	var s models.Seller
	row := r.mysql.QueryRowContext(ctx, querySellerFindById, id)
	err := row.Scan(&s.Id, &s.Cid, &s.CompanyName, &s.Address, &s.Telephone)
	if err != nil {
		if err == sql.ErrNoRows {
			err := api.ServiceErrors[api.ErrNotFound]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "The seller you are looking for does not exist.",
			}
		}
		err := api.ServiceErrors[api.ErrInternalServer]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "An internal server error occurred while retrieving the seller.",
		}
	}

	return &s, nil
}

func (r *sellerRepository) CIDExists(ctx context.Context, cid int, id int) bool {
	var exists bool

	r.mysql.QueryRowContext(
		ctx,
		querySellerCIDExists,
		cid, id,
	).Scan(&exists)

	return exists
}
