package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

const (
	queryBuyerCreate           = `INSERT INTO buyers (id_card_number, first_name, last_name) VALUES (?, ?, ?)`
	queryBuyerUpdate           = `UPDATE buyers SET id_card_number = ?, first_name = ?, last_name = ? WHERE id = ?`
	queryBuyerDelete           = `DELETE FROM buyers WHERE id = ?`
	queryBuyerFindAll          = `SELECT id, id_card_number, first_name, last_name FROM buyers`
	queryBuyerFindById         = `SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = ?`
	queryBuyerCardNumberExists = `SELECT EXISTS(SELECT 1 FROM buyers WHERE LOWER(id_card_number) = LOWER(?) AND id != ?)`
)

func (r *buyerRepository) Create(ctx context.Context, b models.Buyer) (*models.Buyer, error) {
	res, err := r.mysql.ExecContext(
		ctx,
		queryBuyerCreate,
		b.CardNumberId, b.FirstName, b.LastName,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	b.Id = int(id)
	return &b, nil
}

func (r *buyerRepository) Update(ctx context.Context, id int, b models.Buyer) error {
	_, err := r.mysql.ExecContext(
		ctx,
		queryBuyerUpdate,
		b.CardNumberId, b.FirstName, b.LastName, id,
	)
	return err
}

func (r *buyerRepository) Delete(ctx context.Context, id int) error {
	result, err := r.mysql.ExecContext(ctx, queryBuyerDelete, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
			errDef := api.ServiceErrors[api.ErrConflict]
			return &api.ServiceError{
				Code:         errDef.Code,
				ResponseCode: errDef.ResponseCode,
				Message:      "Cannot delete buyer: there are purchase orders associated with this buyer.",
			}
		}
		errDef := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      fmt.Sprintf("An internal server error occurred while deleting the buyer: %s", err.Error()),
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errDef := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      fmt.Sprintf("An internal server error occurred while deleting the buyer: %s", err.Error()),
		}
	}

	if rowsAffected == 0 {
		errDef := api.ServiceErrors[api.ErrNotFound]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "The buyer you are trying to delete does not exist",
		}
	}
	return nil
}

func (r *buyerRepository) FindAll(ctx context.Context) ([]models.Buyer, error) {
	rows, err := r.mysql.QueryContext(ctx, queryBuyerFindAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buyers []models.Buyer
	for rows.Next() {
		var b models.Buyer
		err := rows.Scan(&b.Id, &b.CardNumberId, &b.FirstName, &b.LastName)
		if err != nil {
			return nil, err
		}
		buyers = append(buyers, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buyers, nil
}

func (r *buyerRepository) FindById(ctx context.Context, id int) (*models.Buyer, error) {
	var b models.Buyer
	row := r.mysql.QueryRowContext(ctx, queryBuyerFindById, id)
	err := row.Scan(&b.Id, &b.CardNumberId, &b.FirstName, &b.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			err := api.ServiceErrors[api.ErrNotFound]
			return nil, &api.ServiceError{
				Code:         err.Code,
				ResponseCode: err.ResponseCode,
				Message:      "The buyer you are looking for does not exist.",
			}
		}
		err := api.ServiceErrors[api.ErrInternalServer]
		return nil, &api.ServiceError{
			Code:         err.Code,
			ResponseCode: err.ResponseCode,
			Message:      "An internal server error occurred while retrieving the buyer.",
		}
	}

	return &b, nil
}

func (r *buyerRepository) CardNumberExists(ctx context.Context, cardNumber string, excludeId int) bool {
	var exists bool
	r.mysql.QueryRowContext(
		ctx,
		queryBuyerCardNumberExists,
		cardNumber, excludeId,
	).Scan(&exists)
	return exists
}
