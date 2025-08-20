package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type purchaseOrderRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPurchaseOrderRepository(db *sql.DB) PurchaseOrderRepository {
	return &purchaseOrderRepository{db: db}
}

// SetLogger allows you to inject the logger after creation
func (r *purchaseOrderRepository) SetLogger(l logger.Logger) {
	r.logger = l
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
	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Creating purchase order", map[string]interface{}{
			"order_number":      po.OrderNumber,
			"buyer_id":          po.BuyerID,
			"product_record_id": po.ProductRecordID,
		})
	}

	if !r.recordExists(ctx, queryCheckBuyerExists, po.BuyerID) {
		if r.logger != nil {
			r.logger.Warning(ctx, "purchase-order-repository", "Buyer not found for purchase order creation", map[string]interface{}{
				"buyer_id": po.BuyerID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, fmt.Sprintf("buyer with id %d does not exist", po.BuyerID))
	}

	if !r.recordExists(ctx, queryCheckProductRecordExists, po.ProductRecordID) {
		if r.logger != nil {
			r.logger.Warning(ctx, "purchase-order-repository", "Product record not found for purchase order creation", map[string]interface{}{
				"product_record_id": po.ProductRecordID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, fmt.Sprintf("product record with id %d does not exist", po.ProductRecordID))
	}

	if r.ExistsOrderNumber(ctx, po.OrderNumber) {
		if r.logger != nil {
			r.logger.Warning(ctx, "purchase-order-repository", "Duplicate order number", map[string]interface{}{
				"order_number": po.OrderNumber,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeConflict, "order_number already exists")
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
			if r.logger != nil {
				r.logger.Error(ctx, "purchase-order-repository", "Database error during purchase order creation", err, map[string]interface{}{
					"order_number":     po.OrderNumber,
					"mysql_error_code": mysqlErr.Number,
				})
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, fmt.Sprintf("database error: %v", mysqlErr.Message))
		}
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error creating purchase order", err, map[string]interface{}{
				"order_number": po.OrderNumber,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error creating purchase order")
	}

	id, err := res.LastInsertId()
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Failed to get last insert ID", err, map[string]interface{}{
				"order_number": po.OrderNumber,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error getting last insert id")
	}

	po.ID = int(id)

	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Purchase order created successfully", map[string]interface{}{
			"purchase_order_id": po.ID,
			"order_number":      po.OrderNumber,
		})
	}

	return &po, nil
}

func (r *purchaseOrderRepository) recordExists(ctx context.Context, query string, id int) bool {
	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil && r.logger != nil {
		r.logger.Error(ctx, "purchase-order-repository", "Error checking record existence", err, map[string]interface{}{
			"query": query,
			"id":    id,
		})
	}
	return err == nil && exists
}

func (r *purchaseOrderRepository) GetAll(ctx context.Context) ([]models.PurchaseOrder, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Getting all purchase orders")
	}

	rows, err := r.db.QueryContext(ctx, queryPurchaseOrderGetAll)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error querying all purchase orders", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error querying all purchase orders")
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
			if r.logger != nil {
				r.logger.Error(ctx, "purchase-order-repository", "Error scanning purchase order", err, nil)
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "error scanning purchase order")
		}

		po.OrderDate, err = time.Parse("2006-01-02 15:04:05", orderDateStr)
		if err != nil {
			if r.logger != nil {
				r.logger.Error(ctx, "purchase-order-repository", "Error parsing order date", err, map[string]interface{}{
					"order_date_string": orderDateStr,
				})
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "error parsing order date")
		}

		pos = append(pos, po)
	}

	if err = rows.Err(); err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error after iterating rows", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error after iterating rows")
	}

	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Retrieved all purchase orders successfully", map[string]interface{}{
			"count": len(pos),
		})
	}

	return pos, nil
}

func (r *purchaseOrderRepository) GetByID(ctx context.Context, id int) (*models.PurchaseOrder, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Getting purchase order by ID", map[string]interface{}{
			"purchase_order_id": id,
		})
	}

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
			// No loggear este caso ya que es un comportamiento esperado
			return nil, apperrors.NewAppError(apperrors.CodeNotFound, "purchase order not found")
		}
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error querying purchase order by ID", err, map[string]interface{}{
				"purchase_order_id": id,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error querying purchase order by id")
	}

	po.OrderDate, err = time.Parse("2006-01-02 15:04:05", orderDateStr)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error parsing order date", err, map[string]interface{}{
				"order_date_string": orderDateStr,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error parsing order date")
	}

	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Purchase order found successfully", map[string]interface{}{
			"purchase_order_id": po.ID,
			"order_number":      po.OrderNumber,
		})
	}

	return &po, nil
}

func (r *purchaseOrderRepository) ExistsOrderNumber(ctx context.Context, orderNumber string) bool {
	var exists bool
	err := r.db.QueryRowContext(ctx, queryPurchaseOrderExists, orderNumber).Scan(&exists)
	if err != nil && r.logger != nil {
		r.logger.Error(ctx, "purchase-order-repository", "Error checking order number existence", err, map[string]interface{}{
			"order_number": orderNumber,
		})
	}
	return exists
}

func (r *purchaseOrderRepository) GetCountByBuyer(ctx context.Context, buyerID int) ([]models.BuyerWithPurchaseCount, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Getting purchase count by buyer", map[string]interface{}{
			"buyer_id": buyerID,
		})
	}

	rows, err := r.db.QueryContext(ctx, queryPurchaseCountByBuyer, buyerID)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error querying purchase count by buyer", err, map[string]interface{}{
				"buyer_id": buyerID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error querying purchase count by buyer")
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
			if r.logger != nil {
				r.logger.Error(ctx, "purchase-order-repository", "Error scanning purchase count result", err, nil)
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "error scanning purchase count result")
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error after iterating rows", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error after iterating rows")
	}

	if len(results) == 0 {
		if r.logger != nil {
			r.logger.Warning(ctx, "purchase-order-repository", "Buyer not found for purchase count", map[string]interface{}{
				"buyer_id": buyerID,
			})
		}
		return nil, apperrors.NewAppError(apperrors.CodeNotFound, "buyer not found")
	}

	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Purchase count by buyer retrieved successfully", map[string]interface{}{
			"buyer_id": buyerID,
			"count":    len(results),
		})
	}

	return results, nil
}

func (r *purchaseOrderRepository) GetAllWithPurchaseCount(ctx context.Context) ([]models.BuyerWithPurchaseCount, error) {
	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "Getting all purchase counts")
	}

	rows, err := r.db.QueryContext(ctx, queryAllPurchaseCount)
	if err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error querying all purchase counts", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error querying all purchase counts")
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
			if r.logger != nil {
				r.logger.Error(ctx, "purchase-order-repository", "Error scanning purchase count result", err, nil)
			}
			return nil, apperrors.NewAppError(apperrors.CodeInternal, "error scanning purchase count result")
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		if r.logger != nil {
			r.logger.Error(ctx, "purchase-order-repository", "Error after iterating rows", err, nil)
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "error after iterating rows")
	}

	if r.logger != nil {
		r.logger.Info(ctx, "purchase-order-repository", "All purchase counts retrieved successfully", map[string]interface{}{
			"buyers_count": len(results),
		})
	}

	return results, nil
}
