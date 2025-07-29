package repository

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

// expectPrepareProductRepository registers the 4 Prepare statements executed
// in NewProductRepository() and returns the *ExpectedPrepare for the
// `SELECT … WHERE id = ?` so the caller can chain ExpectQuery().
func expectPrepareProductRepository(m sqlmock.Sqlmock) *sqlmock.ExpectedPrepare {
	sel := m.ExpectPrepare(regexp.QuoteMeta(baseSelect + " WHERE id = ?"))
	m.ExpectPrepare(regexp.QuoteMeta(insertProduct))
	m.ExpectPrepare(regexp.QuoteMeta(updateProduct))
	m.ExpectPrepare(regexp.QuoteMeta(deleteProduct))
	return sel
}

func allColumns() []string {
	return []string{
		"id", "product_code", "description", "width", "height", "length",
		"net_weight", "expiration_rate", "recommended_freezing_temperature",
		"freezing_rate", "product_type_id", "seller_id",
	}
}
