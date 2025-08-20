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
			if r.logger != nil {
				r.logger.Warning(ctx, "buyer-repository", "Duplicate key violation on buyer creation", map[string]interface{}{
					"card_number_id": b.CardNumberId,
					"mysql_error":    mysqlErr.Number,
				})
			}
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create buyer: card number already exists.")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Database error during buyer creation", err, map[string]interface{}{
				"card_number_id": b.CardNumberId,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while creating a buyer: %s", err.Error()))
	}

	id, err := res.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Failed to get last insert ID after buyer creation", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while creating a buyer: %s", err.Error()))
	}

	b.Id = int(id)

	if r.logger != nil {
		r.logger.Info(ctx, "buyer-repository", "Buyer created successfully", map[string]interface{}{
			"buyer_id":       b.Id,
			"card_number_id": b.CardNumberId,
		})
	}

	return &b, nil
}

func (r *buyerRepository) Update(ctx context.Context, id int, b models.Buyer) error {
	_, err := r.mysql.ExecContext(
		ctx,
		queryBuyerUpdate,
		b.CardNumberId, b.FirstName, b.LastName, id,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			if r.logger != nil {
				r.logger.Warning(ctx, "buyer-repository", "Duplicate key violation on buyer update", map[string]interface{}{
					"buyer_id":       id,
					"card_number_id": b.CardNumberId,
					"mysql_error":    mysqlErr.Number,
				})
			}
			return apperrors.NewAppError(apperrors.CodeConflict, "Could not update buyer: card number already exists.")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Database error during buyer update", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while updating a buyer: %s", err.Error()))
	}

	if r.logger != nil {
		r.logger.Info(ctx, "buyer-repository", "Buyer updated successfully", map[string]interface{}{
			"buyer_id": id,
		})
	}

	return nil
}

func (r *buyerRepository) Delete(ctx context.Context, id int) error {
	result, err := r.mysql.ExecContext(ctx, queryBuyerDelete, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1451 {
				if r.logger != nil {
					r.logger.Warning(ctx, "buyer-repository", "Foreign key constraint violation on buyer deletion", map[string]interface{}{
						"buyer_id":    id,
						"mysql_error": mysqlErr.Number,
					})
				}
				return apperrors.NewAppError(apperrors.CodeConflict,
					"cannot delete buyer: there are purchase orders associated")
			}
		}
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Database error during buyer deletion", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal,
			fmt.Sprintf("an internal server error occurred while deleting the buyer: %s", err.Error()))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Failed to check rows affected after buyer deletion", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal,
			"error verifying affected rows")
	}

	if rowsAffected == 0 {
		if r.logger != nil {
			r.logger.Warning(ctx, "buyer-repository", "Attempted to delete non-existent buyer", map[string]interface{}{
				"buyer_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeNotFound,
			"buyer not found")
	}

	if r.logger != nil {
		r.logger.Info(ctx, "buyer-repository", "Buyer deleted successfully", map[string]interface{}{
			"buyer_id": id,
		})
	}

	return nil
}

func (r *buyerRepository) FindAll(ctx context.Context) ([]models.Buyer, error) {
	rows, err := r.mysql.QueryContext(ctx, queryBuyerFindAll)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Database error during find all buyers", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all buyers: %s", err.Error()))
	}
	defer rows.Close()

	buyers := []models.Buyer{}
	for rows.Next() {
		var b models.Buyer
		err := rows.Scan(&b.Id, &b.CardNumberId, &b.FirstName, &b.LastName)
		if err != nil {
			if r.logger != nil {
				r.logger.Error(ctx, "buyer-repository", "Error scanning buyer row", err, nil)
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all buyers: %s", err.Error()))
		}
		buyers = append(buyers, b)
	}

	if err := rows.Err(); err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Error iterating over buyer rows", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all buyers: %s", err.Error()))
	}

	if r.logger != nil {
		r.logger.Info(ctx, "buyer-repository", "Find all buyers completed successfully", map[string]interface{}{
			"buyers_count": len(buyers),
		})
	}

	return buyers, nil
}

func (r *buyerRepository) FindById(ctx context.Context, id int) (*models.Buyer, error) {
	var b models.Buyer
	row := r.mysql.QueryRowContext(ctx, queryBuyerFindById, id)
	err := row.Scan(&b.Id, &b.CardNumberId, &b.FirstName, &b.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			// No loggear este caso ya que es un comportamiento esperado
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The buyer you are looking for does not exist.")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "buyer-repository", "Database error during find buyer by ID", err, map[string]interface{}{
				"buyer_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the buyer.")
	}

	if r.logger != nil {
		r.logger.Info(ctx, "buyer-repository", "Buyer found successfully", map[string]interface{}{
			"buyer_id": id,
		})
	}

	return &b, nil
}

func (r *buyerRepository) CardNumberExists(ctx context.Context, cardNumber string, excludeId int) bool {
	var exists bool
	err := r.mysql.QueryRowContext(
		ctx,
		queryBuyerCardNumberExists,
		cardNumber, excludeId,
	).Scan(&exists)

	if err != nil && r.logger != nil {
		r.logger.Error(ctx, "buyer-repository", "Error checking card number existence", err, map[string]interface{}{
			"card_number": cardNumber,
			"exclude_id":  excludeId,
		})
	}

	return exists
}
