package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

const (
	queryCountryCreate   = "INSERT INTO countries"
	queryCountryFindById = "SELECT id, name FROM countries WHERE LOWER\\(name\\) = LOWER\\(\\?\\)"
)

func TestGeographyRepository_CreateCountry(t *testing.T) {
	type args struct {
		country models.Country
	}
	tests := []struct {
		name           string
		mock           func(mock sqlmock.Sqlmock, country models.Country)
		args           args
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			mock: func(mock sqlmock.Sqlmock, c models.Country) {
				mock.ExpectExec("^" + queryCountryCreate).
					WithArgs(c.Name).
					WillReturnResult(sqlmock.NewResult(42, 1))
			},
			args:    args{country: models.Country{Name: "Argentina"}},
			wantErr: false,
		},
		{
			name: "error - db failure",
			mock: func(mock sqlmock.Sqlmock, c models.Country) {
				mock.ExpectExec("^" + queryCountryCreate).
					WithArgs(c.Name).
					WillReturnError(errors.New("db is down"))
			},
			args:           args{country: models.Country{Name: "Uruguay"}},
			wantErr:        true,
			expectedErrMsg: "failed to create country",
		},
		{
			name: "error - LastInsertId fails",
			mock: func(mock sqlmock.Sqlmock, c models.Country) {
				r := sqlmock.NewErrorResult(errors.New("last insert ID error"))
				mock.ExpectExec("^" + queryCountryCreate).
					WithArgs(c.Name).
					WillReturnResult(r)
			},
			args:           args{country: models.Country{Name: "Chile"}},
			wantErr:        true,
			expectedErrMsg: "failed to get country ID after creation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mock(mock, tt.args.country)
			repo := repository.NewGeographyRepository(db)

			got, err := repo.CreateCountry(context.Background(), db, tt.args.country)

			if !tt.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotZero(t, got.Id)
				require.Equal(t, tt.args.country.Name, got.Name)
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

func TestGeographyRepository_FindCountryByName(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		arg     string
		wantErr bool
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				dummy := testhelpers.CountriesDummyMap[5]
				mock.ExpectQuery("^" + queryCountryFindById).
					WithArgs(dummy.Name).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(dummy.Id, dummy.Name))
			},
			arg:     "Paraguay",
			wantErr: false,
		},
		{
			name: "not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^" + queryCountryFindById).
					WithArgs("Atlantis").
					WillReturnError(sql.ErrNoRows)
			},
			arg:     "Atlantis",
			wantErr: true,
		},
		{
			name: "error - db failure",
			setup: func(mock sqlmock.Sqlmock) {
				dummy := testhelpers.CountriesDummyMap[5]
				mock.ExpectQuery("^" + queryCountryFindById).
					WithArgs(dummy.Name).
					WillReturnError(errors.New("db is down"))
			},
			arg:     "Paraguay",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			tt.setup(mock)
			repo := repository.NewGeographyRepository(db)

			got, err := repo.FindCountryByName(context.Background(), tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.arg, got.Name)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
