package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type purchaseOrderRepository struct {
	db *sql.DB
}

func NewPurchaseOrderRepository(db *sql.DB) PurchaseOrderRepository {
	return &purchaseOrderRepository{db: db}
}

const (
	queryPurchaseOrderCreate = `INSERT INTO purchase_orders 
		(order_number, order_date, tracking_code, buyer_id, product_record_id) 
		VALUES (?, ?, ?, ?, ?)`

	queryPurchaseOrderGetAll = `SELECT id, order_number, order_date, tracking_code, buyer_id, 
		product_record_id FROM purchase_orders`

	queryPurchaseOrderGetByID = `SELECT id, order_number, order_date, tracking_code, buyer_id, 
		product_record_id FROM purchase_orders WHERE id = ?`

	queryCheckProductRecordExists = `SELECT EXISTS(SELECT 1 FROM product_records WHERE id = ?)`
	queryCheckBuyerExists         = `SELECT EXISTS(SELECT 1 FROM buyers WHERE id = ?)`
	queryPurchaseOrderExists      = `SELECT EXISTS(SELECT 1 FROM purchase_orders WHERE order_number = ?)`
	queryPurchaseCountByBuyer     = `SELECT b.id, b.id_card_number, b.first_name, b.last_name, 
		COUNT(po.id) as purchase_orders_count 
		FROM buyers b 
		LEFT JOIN purchase_orders po ON b.id = po.buyer_id 
		WHERE b.id = ? 
		GROUP BY b.id`
	queryAllPurchaseCount = `SELECT b.id, b.id_card_number, b.first_name, b.last_name, 
		COUNT(po.id) as purchase_orders_count 
		FROM buyers b 
		LEFT JOIN purchase_orders po ON b.id = po.buyer_id 
		GROUP BY b.id`
)

func (r *purchaseOrderRepository) Create(ctx context.Context, po models.PurchaseOrder) (*models.PurchaseOrder, error) {
	if !r.recordExists(ctx, queryCheckBuyerExists, po.BuyerID) {
		return nil, createRelationError("buyer", po.BuyerID)
	}

	if !r.recordExists(ctx, queryCheckProductRecordExists, po.ProductRecordID) {
		return nil, createRelationError("product record", po.ProductRecordID)
	}

	// Validar existencia de relaciones
	if err := r.validateRelations(ctx, po); err != nil {
		return nil, err
	}

	// Verificar si el order_number ya existe
	if r.ExistsOrderNumber(ctx, po.OrderNumber) {
		errDef := api.ServiceErrors[api.ErrConflict]
		return nil, &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "order_number already exists",
		}
	}

	res, err := r.db.ExecContext(
		ctx,
		queryPurchaseOrderCreate,
		po.OrderNumber,
		po.OrderDate,
		po.TrackingCode,
		po.BuyerID,
		po.ProductRecordID,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			return nil, handleMySQLError(mysqlErr)
		}
		return nil, fmt.Errorf("error creating purchase order: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", err)
	}

	po.ID = int(id)
	return &po, nil
}

func (r *purchaseOrderRepository) validateRelations(ctx context.Context, po models.PurchaseOrder) error {
	if !r.recordExists(ctx, queryCheckBuyerExists, po.BuyerID) {
		return createRelationError("buyer", po.BuyerID)
	}

	if !r.recordExists(ctx, queryCheckProductRecordExists, po.ProductRecordID) {
		return createRelationError("product record", po.ProductRecordID)
	}

	return nil
}

func (r *purchaseOrderRepository) recordExists(ctx context.Context, query string, id int) bool {
	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return err == nil && exists
}

func createRelationError(entity string, id int) error {
	errDef := api.ServiceErrors[api.ErrConflict]
	return &api.ServiceError{
		Code:         errDef.Code,
		ResponseCode: errDef.ResponseCode,
		Message:      fmt.Sprintf("%s with id %d does not exist", entity, id),
	}
}

func handleMySQLError(err *mysql.MySQLError) error {
	errDef := api.ServiceErrors[api.ErrInternalServer]
	return &api.ServiceError{
		Code:         errDef.Code,
		ResponseCode: errDef.ResponseCode,
		Message:      fmt.Sprintf("database error: %v", err.Message),
	}
}

func (r *purchaseOrderRepository) GetAll(ctx context.Context) ([]models.PurchaseOrder, error) {
	rows, err := r.db.QueryContext(ctx, queryPurchaseOrderGetAll)
	if err != nil {
		return nil, fmt.Errorf("error querying all purchase orders: %w", err)
	}
	defer rows.Close()

	var pos []models.PurchaseOrder
	for rows.Next() {
		var po models.PurchaseOrder
		var orderDateStr string
		err := rows.Scan(
			&po.ID,
			&po.OrderNumber,
			&orderDateStr,
			&po.TrackingCode,
			&po.BuyerID,
			&po.ProductRecordID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning purchase order: %w", err)
		}

		po.OrderDate, err = time.Parse("2006-01-02 15:04:05", orderDateStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing order date: %w", err)
		}

		pos = append(pos, po)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return pos, nil
}

func (r *purchaseOrderRepository) GetByID(ctx context.Context, id int) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	var orderDateStr string

	row := r.db.QueryRowContext(ctx, queryPurchaseOrderGetByID, id)
	err := row.Scan(
		&po.ID,
		&po.OrderNumber,
		&orderDateStr,
		&po.TrackingCode,
		&po.BuyerID,
		&po.ProductRecordID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errDef := api.ServiceErrors[api.ErrNotFound]
			return nil, &api.ServiceError{
				Code:         errDef.Code,
				ResponseCode: errDef.ResponseCode,
				Message:      "purchase order not found",
			}
		}
		return nil, fmt.Errorf("error querying purchase order by id: %w", err)
	}

	po.OrderDate, err = time.Parse("2006-01-02 15:04:05", orderDateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing order date: %w", err)
	}

	return &po, nil
}

func (r *purchaseOrderRepository) ExistsOrderNumber(ctx context.Context, orderNumber string) bool {
	var exists bool
	r.db.QueryRowContext(ctx, queryPurchaseOrderExists, orderNumber).Scan(&exists)
	return exists
}

func (r *purchaseOrderRepository) GetCountByBuyer(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) {
	rows, err := r.db.QueryContext(ctx, queryPurchaseCountByBuyer, buyerID)
	if err != nil {
		return nil, fmt.Errorf("error querying purchase count by buyer: %w", err)
	}
	defer rows.Close()

	var results []models.BuyerWithPurchaseCount
	for rows.Next() {
		var result models.BuyerWithPurchaseCount
		err := rows.Scan(
			&result.ID,
			&result.CardNumberID,
			&result.FirstName,
			&result.LastName,
			&result.PurchaseOrdersCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning purchase count result: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	if len(results) == 0 {
		errDef := api.ServiceErrors[api.ErrNotFound]
		return nil, &api.ServiceError{
			Code:         errDef.Code,
			ResponseCode: errDef.ResponseCode,
			Message:      "buyer not found",
		}
	}

	return results, nil
}

func (r *purchaseOrderRepository) GetAllWithPurchaseCount(ctx context.Context) ([]models.BuyerWithPurchaseCount, error) {
	rows, err := r.db.QueryContext(ctx, queryAllPurchaseCount)
	if err != nil {
		return nil, fmt.Errorf("error querying all purchase counts: %w", err)
	}
	defer rows.Close()

	var results []models.BuyerWithPurchaseCount
	for rows.Next() {
		var result models.BuyerWithPurchaseCount
		err := rows.Scan(
			&result.ID,
			&result.CardNumberID,
			&result.FirstName,
			&result.LastName,
			&result.PurchaseOrdersCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning purchase count result: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return results, nil
}
