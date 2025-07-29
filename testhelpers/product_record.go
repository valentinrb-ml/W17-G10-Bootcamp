package testhelpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_record"
	productrecordmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

/*
NewProductRecordRepoMock spins up a repository backed by sqlmock and
returns:

 1. the repo instance
 2. the sqlmock object (to program Exec/Query expectations)
 3. a cleanup() that checks ExpectationsWereMet and closes the DB
*/
func NewProductRecordRepoMock(tb testing.TB) (repo.ProductRecordRepository, sqlmock.Sqlmock, func()) {
	tb.Helper()

	db, mock, err := productrecordmock.NewSQLMock()
	require.NoError(tb, err)

	mock.ExpectPrepare("INSERT INTO product_records")
	mock.ExpectPrepare("FROM products p.*LEFT JOIN product_records")
	mock.ExpectPrepare("FROM products p.*WHERE p.id")

	repository, err := repo.NewProductRecordRepository(db)
	require.NoError(tb, err)

	cleanup := func() {
		require.NoError(tb, mock.ExpectationsWereMet())
		db.Close()
	}
	return repository, mock, cleanup
}

// BuildProductRecord returns a valid domain struct ready for tests.
func BuildProductRecord() models.ProductRecord {
	return models.ProductRecord{
		ProductRecordCore: models.ProductRecordCore{
			LastUpdateDate: time.Now(),
			PurchasePrice:  10.2,
			SalePrice:      15.5,
			ProductID:      3,
		},
	}
}

// RequireAppErr asserts that the error is *AppError and matches the code.
func RequireAppErr(tb testing.TB, err error, code string) {
	tb.Helper()
	var ae *apperrors.AppError
	require.True(tb, errors.As(err, &ae), "error is not AppError: %v", err)
	require.Equal(tb, code, ae.Code)
}

// BuildProductRecordRequest builds the JSON payload expected by POST /productRecords.
func BuildProductRecordRequest(ts time.Time, purchase, sale float64, productID int) models.ProductRecordRequest {
	return models.ProductRecordRequest{
		Data: models.ProductRecordCore{
			LastUpdateDate: ts,
			PurchasePrice:  purchase,
			SalePrice:      sale,
			ProductID:      productID,
		},
	}
}

/*
DoRequest executes an http.HandlerFunc using httptest.
body can be:
  - string  → sent verbatim (useful for malformed JSON tests)
  - any     → marshalled to JSON
*/
func DoRequest(tb testing.TB, method, url string, body any, h http.HandlerFunc) *httptest.ResponseRecorder {
	tb.Helper()

	var buf bytes.Buffer
	if body != nil {
		switch v := body.(type) {
		case string:
			buf.WriteString(v)
		default:
			require.NoError(tb, json.NewEncoder(&buf).Encode(v))
		}
	}

	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// DecodeAppErr converts an HTTP response body to *AppError.
// It works with two JSON layouts:
//
//  1. Flat error
//     { "code": "BAD_REQUEST", "message": "..." }
//
//  2. Wrapped error (envelope)
//     { "error": { "code": "BAD_REQUEST", "message": "..." } }
func DecodeAppErr(buf *bytes.Buffer) (*apperrors.AppError, error) {
	var raw map[string]any
	if err := json.NewDecoder(buf).Decode(&raw); err != nil {
		return nil, err
	}

	if c, ok := raw["code"].(string); ok {
		return &apperrors.AppError{Code: c}, nil
	}

	if errObj, ok := raw["error"].(map[string]any); ok {
		if c, ok := errObj["code"].(string); ok {
			return &apperrors.AppError{Code: c}, nil
		}
	}
	return &apperrors.AppError{}, nil
}
