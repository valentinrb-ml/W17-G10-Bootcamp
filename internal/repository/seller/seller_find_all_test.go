package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/seller"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSellerRepository_FindAll(t *testing.T) {
	type args struct{}
	tests := []struct {
		name           string
		mock           func(mock sqlmock.Sqlmock)
		args           args
		wantErr        bool
		expectedErrMsg string
		expectedResult []models.Seller
	}{
		{
			name: "success - sellers found",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
				for _, s := range testhelpers.FindAllSellersDummy() {
					rows.AddRow(s.Id, s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId)
				}
				mock.ExpectQuery("^SELECT (.+) FROM sellers").WillReturnRows(rows)
			},
			wantErr:        false,
			expectedResult: testhelpers.FindAllSellersDummy(),
		},
		{
			name: "success - empty slice",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
				mock.ExpectQuery("^SELECT (.+) FROM sellers").WillReturnRows(rows)
			},
			wantErr:        false,
			expectedResult: []models.Seller{},
		},
		{
			name: "error - db query fails",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM sellers").WillReturnError(errors.New("connection lost"))
			},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
		{
			name: "error - scan fails",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
					AddRow(nil, "ABC", "ML", "calle falsa", "1234", 999)
				mock.ExpectQuery("^SELECT (.+) FROM sellers").WillReturnRows(rows)
			},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
		{
			name: "error - rows.Err fails",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
				for _, s := range testhelpers.FindAllSellersDummy() {
					rows.AddRow(s.Id, s.Cid, s.CompanyName, s.Address, s.Telephone, s.LocalityId)
				}
				rows.RowError(0, errors.New("row failure"))
				mock.ExpectQuery("^SELECT (.+) FROM sellers").WillReturnRows(rows)
			},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			tt.mock(mock)
			repo := repository.NewSellerRepository(db)

			got, err := repo.FindAll(context.Background())

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
