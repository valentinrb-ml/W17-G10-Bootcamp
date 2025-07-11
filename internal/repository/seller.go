package repository

import (
	"database/sql"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func (r *sellerRepository) Create(s models.Seller) (*models.Seller, error) {
	res, err := r.mysql.Exec(
		"INSERT INTO sellers (`cid`, `company_name`, `address`, `telephone`) VALUES  (?, ?, ?, ?)",
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

func (r *sellerRepository) Update(id int, s models.Seller) error {
	_, err := r.mysql.Exec(
		"UPDATE sellers SET `cid`=?, `company_name`=?, `address`=?, `telephone`=? WHERE id=?",
		s.Cid, s.CompanyName, s.Address, s.Telephone, s.Id,
	)

	return err
}

func (r *sellerRepository) Delete(id int) error {
	result, err := r.mysql.Exec("DELETE FROM sellers WHERE id = ?", id)
	if err != nil {
		errDef := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "An internal server error occurred while deleting the seller.",
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errDef := api.ServiceErrors[api.ErrInternalServer]
		return &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "An internal server error occurred while deleting the seller.",
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

func (r *sellerRepository) FindAll() ([]models.Seller, error) {
	rows, err := r.mysql.Query(`
        SELECT id, cid, company_name, address, telephone
        FROM sellers`)
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

func (r *sellerRepository) FindById(id int) (*models.Seller, error) {
	var s models.Seller
	row := r.mysql.QueryRow(`
        SELECT id, cid, company_name, address, telephone 
        FROM sellers WHERE id = ?`, id)
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

func (r *sellerRepository) CIDExists(cid int, id int) bool {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1 FROM sellers
            WHERE LOWER(cid) = LOWER(?) AND id != ?
        )
    `
	r.mysql.QueryRow(query, cid, id).Scan(&exists)
	return exists
}
