package repository_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_record"
	productrecordmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_record"
)

/*
Checks:
 1. Happy path → all Preparex succeed.
 2. Failure on first Prepare → constructor must wrap ErrPrepareInsert.
*/
func TestNewProductRecordRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock) // expectations for sqlmock
		wantErr   error                 // sentinel error expected, nil = success
	}{
		{
			name: "success",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectPrepare("INSERT INTO product_records")
				m.ExpectPrepare("FROM products p.*LEFT JOIN product_records")
				m.ExpectPrepare("FROM products p.*WHERE p.id")
			},
		},
		{
			name: "insert_prepare_fails",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectPrepare("INSERT INTO product_records").
					WillReturnError(errors.New("syntax"))
			},
			wantErr: repo.ErrPrepareInsert,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := productrecordmock.NewSQLMock()
			require.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)

			_, err = repo.NewProductRecordRepository(db)
			
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
