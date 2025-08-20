package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

const (
	querySellerCreate   = `INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)`
	querySellerUpdate   = `UPDATE sellers SET cid = ?, company_name = ?, address = ?, telephone = ?, locality_id = ? WHERE id = ?`
	querySellerDelete   = `DELETE FROM sellers WHERE id = ?`
	querySellerFindAll  = `SELECT id, cid, company_name, address, telephone, locality_id FROM sellers`
	querySellerFindById = `SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`
)

func (r *sellerRepository) Create(ctx context.Context, s models.Seller) (*models.Seller, error) {
	res, err := r.mysql.ExecContext(
		ctx,
		querySellerCreate,
		s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				if r.logger != nil {
					r.logger.Warning(ctx, "seller-repository", "Duplicate key violation on seller creation", map[string]interface{}{
						"cid":         s.Cid,
						"locality_id": s.LocalityId,
						"mysql_error": mysqlErr.Number,
					})
				}
				switch {
				case strings.Contains(mysqlErr.Message, "cid"):
					return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create seller due to a data conflict: cid is already used. Please verify your input and try again.")
				case strings.Contains(mysqlErr.Message, "locality_id"):
					return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create seller due to a data conflict: locality is already used. Please verify your input and try again.")
				default:
					return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create seller due to a data conflict. Please verify your input and try again.")
				}
			case 1452:
				if r.logger != nil {
					r.logger.Warning(ctx, "seller-repository", "Foreign key constraint violation on seller creation", map[string]interface{}{
						"locality_id": s.LocalityId,
						"mysql_error": mysqlErr.Number,
					})
				}
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "Unable to create seller: The specified locality does not exist. Please check the locality information and try again.")
			}
		}
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Database error during seller creation", err, map[string]interface{}{
				"cid": s.Cid,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while creating a seller: %s", err.Error()))
	}

	id, err := res.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Failed to get last insert ID after seller creation", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while creating a seller: %s", err.Error()))
	}
	s.Id = int(id)
	return &s, nil
}

func (r *sellerRepository) Update(ctx context.Context, id int, s models.Seller) error {
	_, err := r.mysql.ExecContext(
		ctx,
		querySellerUpdate,
		s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId, s.Id,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				if r.logger != nil {
					r.logger.Warning(ctx, "seller-repository", "Duplicate key violation on seller update", map[string]interface{}{
						"seller_id":   id,
						"cid":         s.Cid,
						"locality_id": s.LocalityId,
						"mysql_error": mysqlErr.Number,
					})
				}
				switch {
				case strings.Contains(mysqlErr.Message, "cid"):
					return apperrors.NewAppError(apperrors.CodeConflict, "Could not update seller due to a data conflict: cid is already used. Please verify your input and try again.")
				case strings.Contains(mysqlErr.Message, "locality_id"):
					return apperrors.NewAppError(apperrors.CodeConflict, "Could not update seller due to a data conflict: locality_id is already used. Please verify your input and try again.")
				default:
					return apperrors.NewAppError(apperrors.CodeConflict, "Could not update seller due to a data conflict. Please verify your input and try again.")
				}
			case 1452:
				if r.logger != nil {
					r.logger.Warning(ctx, "seller-repository", "Foreign key constraint violation on seller update", map[string]interface{}{
						"seller_id":   id,
						"locality_id": s.LocalityId,
						"mysql_error": mysqlErr.Number,
					})
				}
				return apperrors.NewAppError(apperrors.CodeNotFound, "Unable to update seller: The specified locality does not exist. Please check the locality information and try again.")
			}
		}
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Database error during seller update", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while updating a seller: %s", err.Error()))
	}

	return nil
}

func (r *sellerRepository) Delete(ctx context.Context, id int) error {
	result, err := r.mysql.ExecContext(ctx, querySellerDelete, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
			if r.logger != nil {
				r.logger.Warning(ctx, "seller-repository", "Foreign key constraint violation on seller deletion", map[string]interface{}{
					"seller_id":   id,
					"mysql_error": mysqlErr.Number,
				})
			}
			return apperrors.NewAppError(apperrors.CodeConflict, "Cannot delete seller: there are products associated with this seller.")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Database error during seller deletion", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while deleting the seller: %s", err.Error()))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Failed to check rows affected after seller deletion", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while deleting the seller: %s", err.Error()))
	}
	if rowsAffected == 0 {
		if r.logger != nil {
			r.logger.Warning(ctx, "seller-repository", "Attempted to delete non-existent seller", map[string]interface{}{
				"seller_id": id,
			})
		}
		return apperrors.NewAppError(apperrors.CodeNotFound, "The seller you are trying to delete does not exist")
	}

	return nil
}

func (r *sellerRepository) FindAll(ctx context.Context) ([]models.Seller, error) {
	rows, err := r.mysql.QueryContext(ctx, querySellerFindAll)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Database error during find all sellers", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all sellers: %s", err.Error()))
	}
	defer rows.Close()

	sellers := make([]models.Seller, 0)
	for rows.Next() {
		var s models.Seller
		err := rows.Scan(&s.Id, &s.Cid, &s.CompanyName, &s.Address, &s.Telephone, &s.LocalityId)
		if err != nil {
			if r.logger != nil {
				r.logger.Error(ctx, "seller-repository", "Error scanning seller row", err, nil)
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all sellers: %s", err.Error()))
		}
		sellers = append(sellers, s)
	}

	if err := rows.Err(); err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Error iterating over seller rows", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all sellers: %s", err.Error()))
	}

	return sellers, nil
}

func (r *sellerRepository) FindById(ctx context.Context, id int) (*models.Seller, error) {
	var s models.Seller
	row := r.mysql.QueryRowContext(ctx, querySellerFindById, id)
	err := row.Scan(&s.Id, &s.Cid, &s.CompanyName, &s.Address, &s.Telephone, &s.LocalityId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No loggear este caso ya que es un comportamiento esperado
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The seller you are looking for does not exist.")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "seller-repository", "Database error during find seller by ID", err, map[string]interface{}{
				"seller_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "An internal server error occurred while retrieving the seller.")
	}

	return &s, nil
}
