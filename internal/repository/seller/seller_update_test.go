package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/seller"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSellerRepository_Update(t *testing.T) {
	type args struct {
		id     int
		seller models.Seller
	}

	tests := []struct {
		name           string
		mock           func(mock sqlmock.Sqlmock, id int, s models.Seller)
		args           args
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			mock: func(mock sqlmock.Sqlmock, id int, s models.Seller) {
				mock.ExpectExec("^UPDATE sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId, s.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args:    args{id: 1, seller: testhelpers.SellersMapStub[1]},
			wantErr: false,
		},
		{
			name: "error - cid duplicate",
			mock: func(mock sqlmock.Sqlmock, id int, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1062,
					Message: "Duplicate entry '1A' for key 'cid'",
				}
				mock.ExpectExec("^UPDATE sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId, s.Id).
					WillReturnError(mysqlErr)
			},
			args:           args{id: 1, seller: testhelpers.SellersMapStub[1]},
			wantErr:        true,
			expectedErrMsg: "cid is already used",
		},
		{
			name: "error - locality_id duplicate",
			mock: func(mock sqlmock.Sqlmock, id int, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1062,
					Message: "Duplicate entry '5000' for key 'locality_id'",
				}
				mock.ExpectExec("^UPDATE sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId, s.Id).
					WillReturnError(mysqlErr)
			},
			args:           args{id: 1, seller: testhelpers.SellersMapStub[1]},
			wantErr:        true,
			expectedErrMsg: "locality_id is already used",
		},
		{
			name: "error - other data conflict",
			mock: func(mock sqlmock.Sqlmock, id int, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1062,
					Message: "Duplicate entry 'FALSO'",
				}
				mock.ExpectExec("^UPDATE sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId, s.Id).
					WillReturnError(mysqlErr)
			},
			args:           args{id: 1, seller: testhelpers.SellersMapStub[1]},
			wantErr:        true,
			expectedErrMsg: "data conflict",
		},
		{
			name: "error - locality does not exist",
			mock: func(mock sqlmock.Sqlmock, id int, s models.Seller) {
				mysqlErr := &mysql.MySQLError{
					Number:  1452,
					Message: "Cannot add or update a child row: a foreign key constraint fails (`sellers`, CONSTRAINT `sellers_ibfk_1` FOREIGN KEY (`locality_id`) REFERENCES `localities` (`id`))",
				}
				mock.ExpectExec("^UPDATE sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId, s.Id).
					WillReturnError(mysqlErr)
			},
			args:           args{id: 1, seller: testhelpers.SellersMapStub[1]},
			wantErr:        true,
			expectedErrMsg: "locality does not exist",
		},
		{
			name: "error - unknown DB error",
			mock: func(mock sqlmock.Sqlmock, id int, s models.Seller) {
				mock.ExpectExec("^UPDATE sellers").
					WithArgs(s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId, s.Id).
					WillReturnError(errors.New("connection lost"))
			},
			args:           args{id: 1, seller: testhelpers.SellersMapStub[1]},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			tt.mock(mock, tt.args.id, tt.args.seller)
			repo := repository.NewSellerRepository(db)

			err = repo.Update(context.Background(), tt.args.id, tt.args.seller)

			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
