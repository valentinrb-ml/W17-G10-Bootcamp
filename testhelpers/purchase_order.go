package testhelpers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	purchaseOrderRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/purchase_order"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// TestPurchaseOrderHelper provides helper functions for purchase order tests
type TestPurchaseOrderHelper struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
	Repo purchaseOrderRepo.PurchaseOrderRepository
	Ctx  context.Context
}

// NewTestPurchaseOrderHelper creates a new test helper instance
func NewTestPurchaseOrderHelper() (*TestPurchaseOrderHelper, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	//repo := purchaseOrderRepo.(db)
	repo := purchaseOrderRepo.NewPurchaseOrderRepository(db)
	ctx := context.Background()

	return &TestPurchaseOrderHelper{
		DB:   db,
		Mock: mock,
		Repo: repo,
		Ctx:  ctx,
	}, nil
}

// Close closes the database connection
func (h *TestPurchaseOrderHelper) Close() {
	h.DB.Close()
}

// CreateValidPurchaseOrder returns a valid purchase order for testing
func (h *TestPurchaseOrderHelper) CreateValidPurchaseOrder() models.PurchaseOrder {
	return models.PurchaseOrder{
		OrderNumber:     "ORD-001",
		OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		TrackingCode:    "TRK-001",
		BuyerID:         1,
		ProductRecordID: 1,
	}
}

// CreateValidPurchaseOrderWithID returns a valid purchase order with ID for testing
func (h *TestPurchaseOrderHelper) CreateValidPurchaseOrderWithID() models.PurchaseOrder {
	po := h.CreateValidPurchaseOrder()
	po.ID = 1
	return po
}

// CreateBuyerWithPurchaseCount returns a valid buyer with purchase count for testing
func (h *TestPurchaseOrderHelper) CreateBuyerWithPurchaseCount() models.BuyerWithPurchaseCount {
	return models.BuyerWithPurchaseCount{
		ID:                  1,
		CardNumberID:        "12345678",
		FirstName:           "John",
		LastName:            "Doe",
		PurchaseOrdersCount: 5,
	}
}

// MockBuyerExists mocks the buyer existence check
func (h *TestPurchaseOrderHelper) MockBuyerExists(buyerID int, exists bool) {
	query := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM buyers WHERE id = ?)")
	h.Mock.ExpectQuery(query).
		WithArgs(buyerID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(exists))
}

// MockProductRecordExists mocks the product record existence check
func (h *TestPurchaseOrderHelper) MockProductRecordExists(productRecordID int, exists bool) {
	query := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM product_records WHERE id = ?)")
	h.Mock.ExpectQuery(query).
		WithArgs(productRecordID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(exists))
}

// MockOrderNumberExists mocks the order number existence check
func (h *TestPurchaseOrderHelper) MockOrderNumberExists(orderNumber string, exists bool) {
	query := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM purchase_orders WHERE order_number = ?)")
	h.Mock.ExpectQuery(query).
		WithArgs(orderNumber).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(exists))
}

// MockCreatePurchaseOrderSuccess mocks a successful purchase order creation
func (h *TestPurchaseOrderHelper) MockCreatePurchaseOrderSuccess(po models.PurchaseOrder, insertID int64) {
	// Mock buyer exists
	h.MockBuyerExists(po.BuyerID, true)

	// Mock product record exists
	h.MockProductRecordExists(po.ProductRecordID, true)

	// Mock order number doesn't exist
	h.MockOrderNumberExists(po.OrderNumber, false)

	// Mock insert query
	insertQuery := regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES (?, ?, ?, ?, ?)")
	h.Mock.ExpectExec(insertQuery).
		WithArgs(po.OrderNumber, po.OrderDate, po.TrackingCode, po.BuyerID, po.ProductRecordID).
		WillReturnResult(sqlmock.NewResult(insertID, 1))
}

// MockGetAllPurchaseOrdersSuccess mocks a successful get all purchase orders
func (h *TestPurchaseOrderHelper) MockGetAllPurchaseOrdersSuccess(pos []models.PurchaseOrder) {
	query := regexp.QuoteMeta("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders")

	rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"})
	for _, po := range pos {
		rows.AddRow(po.ID, po.OrderNumber, po.OrderDate.Format("2006-01-02 15:04:05"), po.TrackingCode, po.BuyerID, po.ProductRecordID)
	}

	h.Mock.ExpectQuery(query).WillReturnRows(rows)
}

// MockGetPurchaseOrderByIDSuccess mocks a successful get purchase order by ID
func (h *TestPurchaseOrderHelper) MockGetPurchaseOrderByIDSuccess(po models.PurchaseOrder) {
	query := regexp.QuoteMeta("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders WHERE id = ?")

	rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
		AddRow(po.ID, po.OrderNumber, po.OrderDate.Format("2006-01-02 15:04:05"), po.TrackingCode, po.BuyerID, po.ProductRecordID)

	h.Mock.ExpectQuery(query).
		WithArgs(po.ID).
		WillReturnRows(rows)
}

// MockGetPurchaseOrderByIDNotFound mocks a not found scenario for get by ID
func (h *TestPurchaseOrderHelper) MockGetPurchaseOrderByIDNotFound(id int) {
	query := regexp.QuoteMeta("SELECT id, order_number, order_date, tracking_code, buyer_id, product_record_id FROM purchase_orders WHERE id = ?")

	h.Mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)
}

// MockGetCountByBuyerSuccess mocks a successful get count by buyer
func (h *TestPurchaseOrderHelper) MockGetCountByBuyerSuccess(buyerID int, buyers []models.BuyerWithPurchaseCount) {
	query := regexp.QuoteMeta("SELECT b.id, b.id_card_number, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count FROM buyers b LEFT JOIN purchase_orders po ON b.id = po.buyer_id WHERE b.id = ? GROUP BY b.id")

	rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"})
	for _, buyer := range buyers {
		rows.AddRow(buyer.ID, buyer.CardNumberID, buyer.FirstName, buyer.LastName, buyer.PurchaseOrdersCount)
	}

	h.Mock.ExpectQuery(query).
		WithArgs(buyerID).
		WillReturnRows(rows)
}

// MockGetCountByBuyerNotFound mocks a not found scenario for get count by buyer
func (h *TestPurchaseOrderHelper) MockGetCountByBuyerNotFound(buyerID int) {
	query := regexp.QuoteMeta("SELECT b.id, b.id_card_number, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count FROM buyers b LEFT JOIN purchase_orders po ON b.id = po.buyer_id WHERE b.id = ? GROUP BY b.id")

	rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"})

	h.Mock.ExpectQuery(query).
		WithArgs(buyerID).
		WillReturnRows(rows)
}

// MockGetAllWithPurchaseCountSuccess mocks a successful get all with purchase count
func (h *TestPurchaseOrderHelper) MockGetAllWithPurchaseCountSuccess(buyers []models.BuyerWithPurchaseCount) {
	query := regexp.QuoteMeta("SELECT b.id, b.id_card_number, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count FROM buyers b LEFT JOIN purchase_orders po ON b.id = po.buyer_id GROUP BY b.id")

	rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"})
	for _, buyer := range buyers {
		rows.AddRow(buyer.ID, buyer.CardNumberID, buyer.FirstName, buyer.LastName, buyer.PurchaseOrdersCount)
	}

	h.Mock.ExpectQuery(query).WillReturnRows(rows)
}

// MockDatabaseError mocks a database error
func (h *TestPurchaseOrderHelper) MockDatabaseError(query string, args ...driver.Value) {
	h.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(args...).
		WillReturnError(sql.ErrConnDone)
}

// MockExecError mocks an execution error
func (h *TestPurchaseOrderHelper) MockExecError(query string, args ...driver.Value) {
	h.Mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(args...).
		WillReturnError(sql.ErrConnDone)
}

// AssertExpectations verifies that all expectations were met
func (h *TestPurchaseOrderHelper) AssertExpectations() error {
	return h.Mock.ExpectationsWereMet()
}
