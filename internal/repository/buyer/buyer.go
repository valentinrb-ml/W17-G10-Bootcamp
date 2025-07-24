package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
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
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create buyer: card number already exists.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while creating a buyer: %s", err.Error()))
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while creating a buyer: %s", err.Error()))
	}

	b.Id = int(id)
	return &b, nil
}

// func (r *buyerRepository) Delete(ctx context.Context, id int) error {
// 	result, err := r.mysql.ExecContext(ctx, queryBuyerDelete, id)
// 	if err != nil {
// 		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
// 			return apperrors.NewAppError(apperrors.CodeConflict, "Cannot delete buyer: there are purchase orders associated with this buyer.")
// 		}
// 		return apperrors.NewAppError(apperrors.CodeInternal,
// 			fmt.Sprintf("An internal server error occurred while deleting the buyer: %s", err.Error()))
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return apperrors.NewAppError(apperrors.CodeInternal,
// 			fmt.Sprintf("An internal server error occurred while deleting the buyer: %s", err.Error()))
// 	}

//		if rowsAffected == 0 {
//			return apperrors.NewAppError(apperrors.CodeNotFound, "The buyer you are trying to delete does not exist")
//		}
//		return nil
//	}

func (r *buyerRepository) Update(ctx context.Context, id int, b models.Buyer) error {
	_, err := r.mysql.ExecContext(
		ctx,
		queryBuyerUpdate,
		b.CardNumberId, b.FirstName, b.LastName, id,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return apperrors.NewAppError(apperrors.CodeConflict, "Could not update buyer: card number already exists.")
		}
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while updating a buyer: %s", err.Error()))
	}
	return nil
}
func (r *buyerRepository) Delete(ctx context.Context, id int) error {
	result, err := r.mysql.ExecContext(ctx, queryBuyerDelete, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1451 {
				return apperrors.NewAppError(apperrors.CodeConflict,
					"cannot delete buyer: there are purchase orders associated")
			}
		}
		return apperrors.NewAppError(apperrors.CodeInternal,
			fmt.Sprintf("an internal server error occurred while deleting the buyer: %s", err.Error()))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternal,
			"error verifying affected rows")
	}

	if rowsAffected == 0 {
		return apperrors.NewAppError(apperrors.CodeNotFound,
			"buyer not found")
	}

	return nil
}

func (r *buyerRepository) FindAll(ctx context.Context) ([]models.Buyer, error) {
	rows, err := r.mysql.QueryContext(ctx, queryBuyerFindAll)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all buyers: %s", err.Error()))
	}
	defer rows.Close()

	//var buyers []models.Buyer
	buyers := []models.Buyer{}
	for rows.Next() {
		var b models.Buyer
		err := rows.Scan(&b.Id, &b.CardNumberId, &b.FirstName, &b.LastName)
		if err != nil {
			return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all buyers: %s", err.Error()))
		}
		buyers = append(buyers, b)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all buyers: %s", err.Error()))
	}

	return buyers, nil
}

func (r *buyerRepository) FindById(ctx context.Context, id int) (*models.Buyer, error) {
	var b models.Buyer
	row := r.mysql.QueryRowContext(ctx, queryBuyerFindById, id)
	err := row.Scan(&b.Id, &b.CardNumberId, &b.FirstName, &b.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The buyer you are looking for does not exist.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the buyer.")
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
