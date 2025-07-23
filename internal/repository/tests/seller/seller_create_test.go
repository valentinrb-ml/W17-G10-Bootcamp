package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func newTestSeller() models.Seller {
	return models.Seller{
		Cid:         101,
		CompanyName: "Acme",
		Address:     "Test Calle 123",
		Telephone:   "555-1234",
		LocalityId:  "5000",
	}
}

func TestSellerRepository_Create(t *testing.T) {
	type args struct {
		seller models.Seller
	}

	tests := []struct {
		name           string
		mock           func(mock sqlmock.Sqlmock, seller models.Seller)
		args           args
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			mock: func(mock sqlmock.Sqlmock, s models.Seller) {
				mock.ExpectExec("^INSERT INTO sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId).
					WillReturnResult(sqlmock.NewResult(31, 1))
			},
			args:    args{seller: newTestSeller()},
			wantErr: false,
		},
		{
			name: "error - cid duplicate",
			mock: func(mock sqlmock.Sqlmock, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1062,
					Message: "Duplicate entry '1A' for key 'cid'",
				}
				mock.ExpectExec("^INSERT INTO sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId).
					WillReturnError(mysqlErr)
			},
			args:           args{seller: newTestSeller()},
			wantErr:        true,
			expectedErrMsg: "cid is already used",
		},
		{
			name: "error - locality_id duplicate",
			mock: func(mock sqlmock.Sqlmock, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1062,
					Message: "Duplicate entry '42' for key 'locality_id'",
				}
				mock.ExpectExec("^INSERT INTO sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId).
					WillReturnError(mysqlErr)
			},
			args:           args{seller: newTestSeller()},
			wantErr:        true,
			expectedErrMsg: "locality is already used",
		},
		{
			name: "error - other data conflict",
			mock: func(mock sqlmock.Sqlmock, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1062,
					Message: "Duplicate entry 'FALSO'",
				}
				mock.ExpectExec("^INSERT INTO sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId).
					WillReturnError(mysqlErr)
			},
			args:           args{seller: newTestSeller()},
			wantErr:        true,
			expectedErrMsg: "data conflict",
		},
		{
			name: "error - locality does not exist",
			mock: func(mock sqlmock.Sqlmock, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1452,
					Message: "Cannot add or update a child row: a foreign key constraint fails (`sellers`, CONSTRAINT `sellers_ibfk_1` FOREIGN KEY (`locality_id`) REFERENCES `localities` (`id`))",
				}
				mock.ExpectExec("^INSERT INTO sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId).
					WillReturnError(mysqlErr)
			},
			args:           args{seller: newTestSeller()},
			wantErr:        true,
			expectedErrMsg: "locality does not exist",
		},
		{
			name: "error - unknown DB error",
			mock: func(mock sqlmock.Sqlmock, s models.Seller) {
				mock.ExpectExec("^INSERT INTO sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId).
					WillReturnError(errors.New("connection lost"))
			},
			args:           args{seller: newTestSeller()},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			tt.mock(mock, tt.args.seller)
			repo := repository.NewSellerRepository(db)

			got, err := repo.Create(context.Background(), tt.args.seller)

			if !tt.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
			} else {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
