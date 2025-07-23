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
				switch {
				case strings.Contains(mysqlErr.Message, "cid"):
					return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create seller due to a data conflict: cid is already used. Please verify your input and try again.")
				case strings.Contains(mysqlErr.Message, "locality_id"):
					return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create seller due to a data conflict: locality is already used. Please verify your input and try again.")
				default:
					return nil, apperrors.NewAppError(apperrors.CodeConflict, "Could not create seller due to a data conflict. Please verify your input and try again.")
				}
			case 1452:
				return nil, apperrors.NewAppError(apperrors.CodeNotFound, "Unable to create seller: The specified locality does not exist. Please check the locality information and try again.")
			}
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while creating a seller: %s", err.Error()))
	}

	id, err := res.LastInsertId()
	if err != nil {
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
				switch {
				case strings.Contains(mysqlErr.Message, "cid"):
					return apperrors.NewAppError(apperrors.CodeConflict, "Could not update seller due to a data conflict: cid is already used. Please verify your input and try again.")
				case strings.Contains(mysqlErr.Message, "locality_id"):
					return apperrors.NewAppError(apperrors.CodeConflict, "Could not update seller due to a data conflict: locality_id is already used. Please verify your input and try again.")
				default:
					return apperrors.NewAppError(apperrors.CodeConflict, "Could not update seller due to a data conflict. Please verify your input and try again.")
				}
			case 1452:
				return apperrors.NewAppError(apperrors.CodeNotFound, "Unable to update seller: The specified locality does not exist. Please check the locality information and try again.")
			}
		}
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while updating a seller: %s", err.Error()))
	}

	return nil
}

func (r *sellerRepository) Delete(ctx context.Context, id int) error {
	result, err := r.mysql.ExecContext(ctx, querySellerDelete, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
			return apperrors.NewAppError(apperrors.CodeConflict, "Cannot delete seller: there are products associated with this seller.")
		}
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while deleting the seller: %s", err.Error()))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while deleting the seller: %s", err.Error()))
	}
	if rowsAffected == 0 {
		return apperrors.NewAppError(apperrors.CodeNotFound, "The seller you are trying to delete does not exist")
	}

	return nil
}

func (r *sellerRepository) FindAll(ctx context.Context) ([]models.Seller, error) {
	rows, err := r.mysql.QueryContext(ctx, querySellerFindAll)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all sellers: %s", err.Error()))
	}
	defer rows.Close()

	sellers := make([]models.Seller, 0)
	for rows.Next() {
		var s models.Seller
		err := rows.Scan(&s.Id, &s.Cid, &s.CompanyName, &s.Address, &s.Telephone, &s.LocalityId)
		if err != nil {
			return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("An internal server error occurred while finding all sellers: %s", err.Error()))
		}
		sellers = append(sellers, s)
	}

	if err := rows.Err(); err != nil {
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
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The seller you are looking for does not exist.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "An internal server error occurred while retrieving the seller.")
	}

	return &s, nil
}
