package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/seller"
)

func TestSellerRepository_Delete(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name           string
		mock           func(mock sqlmock.Sqlmock, id int)
		args           args
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			mock: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectExec("^DELETE FROM sellers").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args:    args{id: 1},
			wantErr: false,
		},
		{
			name: "error - DB foreign key conflict (1451, seller has products)",
			mock: func(mock sqlmock.Sqlmock, id int) {
				mysqlErr := &mysql.MySQLError{
					Number:  1451,
					Message: "Cannot delete or update a parent row: a foreign key constraint fails",
				}
				mock.ExpectExec("^DELETE FROM sellers").
					WithArgs(id).
					WillReturnError(mysqlErr)
			},
			args:           args{id: 1},
			wantErr:        true,
			expectedErrMsg: "associated with this seller",
		},
		{
			name: "error - other DB error",
			mock: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectExec("^DELETE FROM sellers").
					WithArgs(id).
					WillReturnError(errors.New("db connection lost"))
			},
			args:           args{id: 1},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
		{
			name: "error - rows affected error",
			mock: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectExec("^DELETE FROM sellers").
					WithArgs(id).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rowsAffected fail")))
			},
			args:           args{id: 42},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
		{
			name: "error - seller not found (rows affected = 0)",
			mock: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectExec("^DELETE FROM sellers").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			args:           args{id: 999},
			wantErr:        true,
			expectedErrMsg: "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mock(mock, tt.args.id)
			repo := repository.NewSellerRepository(db)

			err = repo.Delete(context.Background(), tt.args.id)
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
