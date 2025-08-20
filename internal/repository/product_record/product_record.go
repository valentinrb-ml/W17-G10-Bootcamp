package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

var (
	ErrPrepareInsert          = errors.New("repository: prepare insert stmt")
	ErrPrepareReportAll       = errors.New("repository: prepare report-all stmt")
	ErrPrepareReportByProduct = errors.New("repository: prepare report-by-product stmt")
)

const (
	productRecordQueryTimeout = 5 * time.Second
	insertProductRecord       = `
		INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) 
		VALUES (?, ?, ?, ?)`
	selectAllProductsReport = `
		SELECT 
			p.id as product_id,
			p.description,
			COALESCE(COUNT(pr.id), 0) as records_count
		FROM products p
		LEFT JOIN product_records pr ON p.id = pr.product_id
		GROUP BY p.id, p.description
		ORDER BY p.id`
	selectProductReportByID = `
		SELECT 
			p.id as product_id,
			p.description,
			COALESCE(COUNT(pr.id), 0) as records_count
		FROM products p
		LEFT JOIN product_records pr ON p.id = pr.product_id
		WHERE p.id = ?
		GROUP BY p.id, p.description`
)

const prRepoServiceName = "product-record-repository" // [LOG]

type productRecordMySQLRepository struct {
	db                  *sqlx.DB
	stmtInsert          *sqlx.Stmt
	stmtReportAll       *sqlx.Stmt
	stmtReportByProduct *sqlx.Stmt

	logger logger.Logger // [LOG]
}

func NewProductRecordRepository(db *sql.DB) (ProductRecordRepository, error) {
	xdb := sqlx.NewDb(db, "mysql")

	insert, err := xdb.Preparex(insertProductRecord)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPrepareInsert, err)
	}

	reportAll, err := xdb.Preparex(selectAllProductsReport)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPrepareReportAll, err)
	}

	reportByProduct, err := xdb.Preparex(selectProductReportByID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPrepareReportByProduct, err)
	}

	return &productRecordMySQLRepository{
		db:                  xdb,
		stmtInsert:          insert,
		stmtReportAll:       reportAll,
		stmtReportByProduct: reportByProduct,
	}, nil
}

// SetLogger allows injecting the logger after creation
func (r *productRecordMySQLRepository) SetLogger(l logger.Logger) { // [LOG]
	r.logger = l // [LOG]
}

// logging helpers (avoid repeating nil-check and set the service name)
func (r *productRecordMySQLRepository) debug(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if r.logger != nil {
		r.logger.Debug(ctx, prRepoServiceName, msg, md...) // [LOG]
	}
}
func (r *productRecordMySQLRepository) info(ctx context.Context, msg string, md ...map[string]interface{}) { // [LOG]
	if r.logger != nil {
		r.logger.Info(ctx, prRepoServiceName, msg, md...) // [LOG]
	}
}

func (r *productRecordMySQLRepository) Create(ctx context.Context, record models.ProductRecord) (models.ProductRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, productRecordQueryTimeout)
	defer cancel()

	r.info(ctx, "Inserting product record", map[string]interface{}{ // [LOG]
		"product_id": record.ProductID, // [LOG]
	})

	res, err := r.stmtInsert.ExecContext(ctx,
		record.LastUpdateDate,
		record.PurchasePrice,
		record.SalePrice,
		record.ProductID)
	if err != nil {
		return models.ProductRecord{}, r.handleDBError(err, "failed to create product record")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.ProductRecord{}, apperrors.Wrap(err, "failed to fetch new product record id")
	}

	record.ID = int(id)

	r.info(ctx, "Product record inserted", map[string]interface{}{ // [LOG]
		"product_record_id": record.ID, // [LOG]
		"product_id":        record.ProductID,
	})

	return record, nil
}

func (r *productRecordMySQLRepository) GetRecordsReport(ctx context.Context, productID int) ([]models.ProductRecordReport, error) {
	ctx, cancel := context.WithTimeout(ctx, productRecordQueryTimeout)
	defer cancel()

	var reports []models.ProductRecordReport

	if productID == 0 {
		r.debug(ctx, "Executing records report for all products") // [LOG]

		if err := r.stmtReportAll.SelectContext(ctx, &reports); err != nil {
			return nil, apperrors.Wrap(err, "failed to get records report for all products")
		}

		r.info(ctx, "Records report fetched (all products)", map[string]interface{}{ // [LOG]
			"count": len(reports), // [LOG]
		})
	} else {
		r.debug(ctx, "Executing records report by product", map[string]interface{}{ // [LOG]
			"product_id": productID, // [LOG]
		})

		if err := r.stmtReportByProduct.SelectContext(ctx, &reports, productID); err != nil {
			return nil, apperrors.Wrap(err, "failed to get records report for specific product")
		}

		if len(reports) == 0 {
			return nil, apperrors.NewAppError(
				apperrors.CodeNotFound,
				fmt.Sprintf("product with id %d not found", productID),
			)
		}

		r.info(ctx, "Records report fetched (filtered)", map[string]interface{}{ // [LOG]
			"product_id": productID,    // [LOG]
			"count":      len(reports), // [LOG]
		})
	}

	return reports, nil
}

func (r *productRecordMySQLRepository) handleDBError(err error, msg string) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1452: // Foreign key constraint violation
			return apperrors.NewAppError(
				apperrors.CodeConflict,
				"product_id does not exist",
			)
		case 1048: // Column cannot be null
			return apperrors.NewAppError(
				apperrors.CodeBadRequest,
				"required fields cannot be null",
			)
		case 1406: // Data too long
			return apperrors.NewAppError(
				apperrors.CodeBadRequest,
				"data too long for column",
			)
		}
	}
	return apperrors.Wrap(err, msg)
}
