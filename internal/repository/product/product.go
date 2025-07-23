package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product"
	"strings"
	"time"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	queryTimeout = 5 * time.Second
	baseSelect   = `SELECT id, product_code, description, width, height, length,
	                     net_weight, expiration_rate, recommended_freezing_temperature,
	                     freezing_rate, product_type_id, seller_id
	               FROM products`
	insertProduct = `
		INSERT INTO products (
		  product_code, description, width, height, length,
		  net_weight, expiration_rate, recommended_freezing_temperature,
		  freezing_rate, product_type_id, seller_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	updateProduct = `
		UPDATE products SET
		  product_code = ?, description = ?, width = ?, height = ?,
		  length = ?, net_weight = ?, expiration_rate = ?,
		  recommended_freezing_temperature = ?, freezing_rate = ?,
		  product_type_id = ?, seller_id = ?
		WHERE id = ?`
	deleteProduct = `DELETE FROM products WHERE id = ?`
)

type productMySQLRepository struct {
	db         *sqlx.DB
	stmtByID   *sqlx.Stmt
	stmtInsert *sqlx.Stmt
	stmtUpdate *sqlx.Stmt
	stmtDelete *sqlx.Stmt
}

func NewProductRepository(db *sql.DB) (ProductRepository, error) {
	xdb := sqlx.NewDb(db, "mysql")

	// Only critical sentences are prepared
	selByID, err := xdb.Preparex(baseSelect + " WHERE id = ?")
	if err != nil {
		return nil, err
	}
	insert, err := xdb.Preparex(insertProduct)
	if err != nil {
		return nil, err
	}
	update, err := xdb.Preparex(updateProduct)
	if err != nil {
		return nil, err
	}
	deleteStmt, err := xdb.Preparex(deleteProduct)
	if err != nil {
		return nil, err
	}

	return &productMySQLRepository{
		db:         xdb,
		stmtByID:   selByID,
		stmtInsert: insert,
		stmtUpdate: update,
		stmtDelete: deleteStmt,
	}, nil
}

// CRUD

func (r *productMySQLRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var dbRows []models.ProductDb
	query := baseSelect + " ORDER BY id"

	if err := r.db.SelectContext(ctx, &dbRows, query); err != nil {
		fmt.Println(err)
		return nil, apperrors.Wrap(err, "failed to get all products")
	}

	products := make([]models.Product, len(dbRows))
	for i, dp := range dbRows {
		products[i] = product.DbToDomain(dp)
	}
	return products, nil
}

func (r *productMySQLRepository) GetByID(ctx context.Context, id int) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var dp models.ProductDb
	if err := r.stmtByID.GetContext(ctx, &dp, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Product{}, apperrors.NewAppError(
				apperrors.CodeNotFound,
				fmt.Sprintf("product with id %d not found", id),
			)
		}
		return models.Product{}, apperrors.Wrap(err, "failed to get product by id")
	}
	return product.DbToDomain(dp), nil
}

func (r *productMySQLRepository) Save(ctx context.Context, p models.Product) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	if p.ID == 0 {
		return r.create(ctx, p)
	}
	return r.update(ctx, p)
}

func (r *productMySQLRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	res, err := r.stmtDelete.ExecContext(ctx, id)
	if err != nil {
		return r.handleDBError(err, "failed to delete product")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return apperrors.NewAppError(apperrors.CodeNotFound, "product not found")
	}
	return nil
}

func (r *productMySQLRepository) Patch(ctx context.Context, id int, req models.ProductPatchRequest) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var (
		fields []string
		args   []interface{}
	)

	if req.ProductCode != nil {
		fields = append(fields, "product_code = ?")
		args = append(args, *req.ProductCode)
	}
	if req.Description != nil {
		fields = append(fields, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Width != nil {
		fields = append(fields, "width = ?")
		args = append(args, *req.Width)
	}
	if req.Height != nil {
		fields = append(fields, "height = ?")
		args = append(args, *req.Height)
	}
	if req.Length != nil {
		fields = append(fields, "length = ?")
		args = append(args, *req.Length)
	}
	if req.NetWeight != nil {
		fields = append(fields, "net_weight = ?")
		args = append(args, *req.NetWeight)
	}
	if req.ExpirationRate != nil {
		fields = append(fields, "expiration_rate = ?")
		args = append(args, *req.ExpirationRate)
	}
	if req.RecommendedFreezingTemperature != nil {
		fields = append(fields, "recommended_freezing_temperature = ?")
		args = append(args, *req.RecommendedFreezingTemperature)
	}
	if req.FreezingRate != nil {
		fields = append(fields, "freezing_rate = ?")
		args = append(args, *req.FreezingRate)
	}
	if req.ProductTypeID != nil {
		fields = append(fields, "product_type_id = ?")
		args = append(args, *req.ProductTypeID)
	}
	if req.SellerID != nil {
		fields = append(fields, "seller_id = ?")
		args = append(args, *req.SellerID)
	}

	if len(fields) == 0 { // nothing to modify
		return r.GetByID(ctx, id)
	}

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = ?", strings.Join(fields, ", "))
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return models.Product{}, r.handleDBError(err, "failed to patch product")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return models.Product{}, apperrors.NewAppError(apperrors.CodeNotFound, "product not found")
	}

	return r.GetByID(ctx, id)
}

// privates (create / update)
func (r *productMySQLRepository) create(ctx context.Context, p models.Product) (models.Product, error) {
	d := product.FromDomainToDb(p)

	res, err := r.stmtInsert.ExecContext(ctx,
		d.Code, d.Description, d.Width, d.Height, d.Length,
		d.NetWeight, d.ExpRate, d.RecFreeze, d.FreezeRate,
		d.TypeID, d.SellerID)
	if err != nil {
		return models.Product{}, r.handleDBError(err, "failed to create product")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Product{}, apperrors.Wrap(err, "failed to fetch new id")
	}
	p.ID = int(id)
	return p, nil
}

func (r *productMySQLRepository) update(ctx context.Context, p models.Product) (models.Product, error) {
	d := product.FromDomainToDb(p)

	_, err := r.stmtUpdate.ExecContext(ctx,
		d.Code, d.Description, d.Width, d.Height, d.Length,
		d.NetWeight, d.ExpRate, d.RecFreeze, d.FreezeRate,
		d.TypeID, d.SellerID, d.ID)
	if err != nil {
		return models.Product{}, r.handleDBError(err, "failed to update product")
	}
	return p, nil
}

// MySQL error translation
func (r *productMySQLRepository) handleDBError(err error, msg string) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062: // ER_DUP_ENTRY - Duplicate key
			return apperrors.NewAppError(apperrors.CodeConflict, "product code already exists")

		case 1451: // ER_ROW_IS_REFERENCED_2 - Cannot delete (has product_records)
			return apperrors.NewAppError(apperrors.CodeConflict, "cannot delete product: it has associated product records")

		case 1452: // ER_NO_REFERENCED_ROW_2 - FK constraint fails
			if strings.Contains(mysqlErr.Message, "product_type_id") {
				return apperrors.NewAppError(apperrors.CodeBadRequest, "product_type_id does not exist")
			}
			if strings.Contains(mysqlErr.Message, "seller_id") {
				return apperrors.NewAppError(apperrors.CodeBadRequest, "seller_id does not exist")
			}
			return apperrors.NewAppError(apperrors.CodeBadRequest, "referenced record does not exist")

		case 1048: // ER_BAD_NULL_ERROR
			return apperrors.NewAppError(apperrors.CodeBadRequest, "required field cannot be null")

		case 1406: // ER_DATA_TOO_LONG
			return apperrors.NewAppError(apperrors.CodeBadRequest, "data too long for field")
		}
	}
	return apperrors.Wrap(err, msg)
}
