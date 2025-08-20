package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	mappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
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

const repoServiceName = "product-repository" // [LOG]

type productMySQLRepository struct {
	db         *sqlx.DB
	stmtByID   *sqlx.Stmt
	stmtInsert *sqlx.Stmt
	stmtUpdate *sqlx.Stmt
	stmtDelete *sqlx.Stmt

	logger logger.Logger
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

// SetLogger allows you to inject the logger after creation
func (r *productMySQLRepository) SetLogger(l logger.Logger) {
	r.logger = l
}

// --- logging helpers (avoid repeating nil-check and set the service name) ---
func (r *productMySQLRepository) debug(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if r.logger != nil {
		r.logger.Debug(ctx, repoServiceName, msg, md...) // [LOG]
	}
}
func (r *productMySQLRepository) info(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if r.logger != nil {
		r.logger.Info(ctx, repoServiceName, msg, md...) // [LOG]
	}
}

// CRUD

func (r *productMySQLRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	r.debug(ctx, "Executing GetAll query") // [LOG]

	var dbRows []models.ProductDb
	query := baseSelect + " ORDER BY id"

	if err := r.db.SelectContext(ctx, &dbRows, query); err != nil {
		return nil, apperrors.Wrap(err, "failed to get all products")
	}

	r.debug(ctx, "GetAll query completed", map[string]interface{}{ // [LOG]
		"count": len(dbRows),
	})

	products := make([]models.Product, len(dbRows))
	for i, dp := range dbRows {
		products[i] = mappers.DbToDomain(dp)
	}
	return products, nil
}

func (r *productMySQLRepository) GetByID(ctx context.Context, id int) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	r.debug(ctx, "Executing GetByID", map[string]interface{}{ // [LOG]
		"product_id": id,
	})

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

	r.debug(ctx, "GetByID completed", map[string]interface{}{ // [LOG]
		"product_id": id,
	})

	return mappers.DbToDomain(dp), nil
}

func (r *productMySQLRepository) Save(ctx context.Context, p models.Product) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	if p.ID == 0 {
		r.info(ctx, "Creating product", map[string]interface{}{ // [LOG]
			"product_code": p.Code,
		})
		out, err := r.create(ctx, p)
		if err != nil {
			return models.Product{}, err // handler logs the error // [LOG]
		}
		r.info(ctx, "Product created", map[string]interface{}{ // [LOG]
			"product_id":   out.ID,
			"product_code": out.Code,
		})
		return out, nil
	}

	r.info(ctx, "Updating product", map[string]interface{}{ // [LOG]
		"product_id": p.ID,
	})
	out, err := r.update(ctx, p)
	if err != nil {
		return models.Product{}, err // handler logs the error // [LOG]
	}
	r.info(ctx, "Product updated", map[string]interface{}{ // [LOG]
		"product_id": p.ID,
	})
	return out, nil
}

func (r *productMySQLRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	r.info(ctx, "Deleting product", map[string]interface{}{ // [LOG]
		"product_id": id,
	})

	res, err := r.stmtDelete.ExecContext(ctx, id)
	if err != nil {
		return r.handleDBError(err, "failed to delete product")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return apperrors.NewAppError(apperrors.CodeNotFound, "product not found")
	}

	r.info(ctx, "Product deleted", map[string]interface{}{ // [LOG]
		"product_id": id,
	})

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
		r.debug(ctx, "Patch called with no fields; returning current product", map[string]interface{}{ // [LOG]
			"product_id": id,
		})
		return r.GetByID(ctx, id)
	}

	r.info(ctx, "Patching product", map[string]interface{}{ // [LOG]
		"product_id":   id,
		"fields_count": len(fields),
	})

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = ?", strings.Join(fields, ", "))
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return models.Product{}, r.handleDBError(err, "failed to patch product")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return models.Product{}, apperrors.NewAppError(apperrors.CodeNotFound, "product not found")
	}

	r.info(ctx, "Product patched", map[string]interface{}{ // [LOG]
		"product_id": id,
	})

	return r.GetByID(ctx, id)
}

// privates (create / update)
func (r *productMySQLRepository) create(ctx context.Context, p models.Product) (models.Product, error) {
	d := mappers.FromDomainToDb(p)

	r.debug(ctx, "Executing INSERT for product") // [LOG]

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

	r.debug(ctx, "INSERT completed", map[string]interface{}{ // [LOG]
		"product_id": p.ID,
	})

	return p, nil
}

func (r *productMySQLRepository) update(ctx context.Context, p models.Product) (models.Product, error) {
	d := mappers.FromDomainToDb(p)

	r.debug(ctx, "Executing UPDATE for product", map[string]interface{}{ // [LOG]
		"product_id": p.ID,
	})

	_, err := r.stmtUpdate.ExecContext(ctx,
		d.Code, d.Description, d.Width, d.Height, d.Length,
		d.NetWeight, d.ExpRate, d.RecFreeze, d.FreezeRate,
		d.TypeID, d.SellerID, d.ID)
	if err != nil {
		return models.Product{}, r.handleDBError(err, "failed to update product")
	}

	r.debug(ctx, "UPDATE completed", map[string]interface{}{ // [LOG]
		"product_id": p.ID,
	})

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
