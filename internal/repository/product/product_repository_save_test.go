package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	mappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Create/update happy path plus various codes and generic DB errors.
func TestProductRepository_Save(t *testing.T) {
	t.Parallel()

	errLastID := sqlmock.NewErrorResult(sql.ErrNoRows)

	referencedErr := &mysql.MySQLError{Number: 1452, Message: "other_fk"}
	nullErr := &mysql.MySQLError{Number: 1048, Message: "null"}
	dupErr := &mysql.MySQLError{Number: 1062, Message: "dup"}
	dataLong := &mysql.MySQLError{Number: 1406, Message: "long"}

	tests := []struct {
		name    string
		product models.Product
		dbErr   error
		dbRes   driver.Result
		appCode string
	}{
		{"create ok", testhelpers.BuildProduct(0), nil, sqlmock.NewResult(5, 1), ""},
		{"update ok", testhelpers.BuildProduct(30), nil, sqlmock.NewResult(0, 1), ""},
		{"duplicate code", testhelpers.BuildProduct(0), dupErr, nil, apperrors.CodeConflict},
		{"fk product_type", testhelpers.BuildProduct(0),
			&mysql.MySQLError{Number: 1452, Message: "product_type_id"}, nil, apperrors.CodeBadRequest},
		{"fk seller", testhelpers.BuildProduct(40),
			&mysql.MySQLError{Number: 1452, Message: "seller_id"}, nil, apperrors.CodeBadRequest},
		{"fk other", testhelpers.BuildProduct(0), referencedErr, nil, apperrors.CodeBadRequest},
		{"null field", testhelpers.BuildProduct(0), nullErr, nil, apperrors.CodeBadRequest},
		{"data too long", testhelpers.BuildProduct(0), dataLong, nil, apperrors.CodeBadRequest},
		{"lastInsertId error", testhelpers.BuildProduct(0), nil, errLastID, apperrors.CodeInternal},
		{"generic db wrap", testhelpers.BuildProduct(55), sql.ErrConnDone, nil, apperrors.CodeInternal},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, cancel := productmock.NewMockDB(t)
			defer cancel()

			expectPrepareProductRepository(mock)

			if tc.product.ID == 0 { // create
				ex := mock.ExpectExec(regexp.QuoteMeta(insertProduct)).
					WithArgs(buildInsertArgs(tc.product)...)
				if tc.dbErr != nil {
					ex.WillReturnError(tc.dbErr)
				} else {
					ex.WillReturnResult(tc.dbRes)
				}
			} else { // update
				ex := mock.ExpectExec(regexp.QuoteMeta(updateProduct)).
					WithArgs(buildUpdateArgs(tc.product)...)
				if tc.dbErr != nil {
					ex.WillReturnError(tc.dbErr)
				} else {
					ex.WillReturnResult(tc.dbRes)
				}
			}

			repo, _ := NewProductRepository(db)
			_, err := repo.Save(context.Background(), tc.product)

			if tc.appCode != "" {
				require.Error(t, err)
				testhelpers.RequireAppErr(t, err, tc.appCode)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

/* local helpers */
func buildInsertArgs(p models.Product) []driver.Value {
	d := mappers.FromDomainToDb(p)
	return []driver.Value{
		d.Code, d.Description, d.Width, d.Height, d.Length,
		d.NetWeight, d.ExpRate, d.RecFreeze, d.FreezeRate,
		d.TypeID, d.SellerID,
	}
}
func buildUpdateArgs(p models.Product) []driver.Value {
	d := mappers.FromDomainToDb(p)
	return append(buildInsertArgs(p), d.ID)
}
