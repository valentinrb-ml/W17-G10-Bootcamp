package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

const (
	queryCreateProductBatch    = `INSERT INTO product_batches (batch_number,current_quantity,current_temperature,due_date,initial_quantity,manufacturing_date,manufacturing_hour,minimum_temperature,product_id,section_id) VALUES (?,?,?,?,?,?,?,?,?,?)`
	queryGetReportProductsById = `SELECT s.id, s.section_number, SUM(p.current_quantity) FROM product_batches p INNER JOIN sections s on p.section_id = s.id  WHERE p.section_id = ? GROUP BY p.section_id`
	queryGetProductsReport     = `SELECT s.id, s.section_number, SUM(p.current_quantity) FROM product_batches p INNER JOIN sections s on p.section_id = s.id  GROUP BY p.section_id`
)

// CreateProductBatches inserts a new product batch into the database and returns the created batch.
// Returns error if a duplicate batch number or invalid foreign keys are provided.
func (r *productBatchesRepository) CreateProductBatches(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
	result, err := r.mysql.ExecContext(ctx, queryCreateProductBatch, proBa.BatchNumber, proBa.CurrentQuantity, proBa.CurrentTemperature, proBa.DueDate, proBa.InitialQuantity, proBa.ManufacturingDate, proBa.ManufacturingHour, proBa.MinimumTemperature, proBa.ProductId, proBa.SectionId)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, apperrors.NewAppError(apperrors.CodeConflict, "Batch number already exists.")
		}
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "Section id or product id does not exist.")
		}
		fmt.Println(err)
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while creating the Product Batch.")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	proBa.Id = int(id)

	return &proBa, nil
}

// GetReportProductById returns product report data by section id, including section info and sum of current quantities.
// Returns error if the section is not found.
func (r *productBatchesRepository) GetReportProductById(ctx context.Context, id int) (*models.ReportProduct, error) {
	var pr models.ReportProduct
	err := r.mysql.QueryRowContext(ctx, queryGetReportProductsById, id).Scan(&pr.SectionId, &pr.SectionNumber, &pr.ProductsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "The section you are looking for does not exist.")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the product report.")
	}
	return &pr, nil
}

// GetReportProduct returns a report for all sections, each with section info and sum of current quantities.
// Returns error if there was a problem during the query.
func (r *productBatchesRepository) GetReportProduct(ctx context.Context) ([]models.ReportProduct, error) {
	rows, err := r.mysql.QueryContext(ctx, queryGetProductsReport)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the products report.")
	}
	defer rows.Close()

	var productReport []models.ReportProduct

	for rows.Next() {
		var rp models.ReportProduct
		if err := rows.Scan(&rp.SectionId, &rp.SectionNumber, &rp.ProductsCount); err != nil {
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the product report.")
		}
		productReport = append(productReport, rp)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the product report.")
	}
	return productReport, nil

}
