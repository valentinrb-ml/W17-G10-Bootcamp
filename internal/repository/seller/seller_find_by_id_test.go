package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/seller"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSellerRepository_FindById(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name           string
		mock           func(mock sqlmock.Sqlmock, id int)
		args           args
		wantErr        bool
		expectedResult *models.Seller
		expectedErrMsg string
	}{
		{
			name: "success",
			mock: func(mock sqlmock.Sqlmock, id int) {
				s := testhelpers.SellersMapStub[1]
				rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
					AddRow(s.Id, s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId)
				mock.ExpectQuery("^SELECT (.+) FROM sellers WHERE id").
					WithArgs(id).
					WillReturnRows(rows)
			},
			args:    args{id: 1},
			wantErr: false,
			expectedResult: func() *models.Seller {
				s := testhelpers.SellersMapStub[1]
				return &s
			}(),
		},
		{
			name: "error - not found",
			mock: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectQuery("^SELECT (.+) FROM sellers WHERE id").
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			args:           args{id: 42},
			wantErr:        true,
			expectedErrMsg: "does not exist",
		},
		{
			name: "error - scan fails",
			mock: func(mock sqlmock.Sqlmock, id int) {
				rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
					AddRow(nil, "xxx", "zzz", "dir", "tel", "locid")
				mock.ExpectQuery("^SELECT (.+) FROM sellers WHERE id").
					WithArgs(id).
					WillReturnRows(rows)
			},
			args:           args{id: 7},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mock(mock, tt.args.id)
			repo := repository.NewSellerRepository(db)

			got, err := repo.FindById(context.Background(), tt.args.id)

			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, got)
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
