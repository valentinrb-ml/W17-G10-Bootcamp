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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

const (
	queryProvinceCreate   = "INSERT INTO provinces"
	queryProvinceFindById = "SELECT id, name, country_id FROM provinces WHERE LOWER\\(name\\) = LOWER\\(\\?\\) AND country_id = \\?"
)

func TestGeographyRepository_CreateProvince(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(sqlmock.Sqlmock)
		arg            models.Province
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				p := testhelpers.ProvincesDummyMap[1]
				mock.ExpectExec("^"+queryProvinceCreate).
					WithArgs(p.Name, p.CountryId).
					WillReturnResult(sqlmock.NewResult(31, 1))
			},
			arg:     testhelpers.ProvincesDummyMap[1],
			wantErr: false,
		},
		{
			name: "db error",
			setup: func(mock sqlmock.Sqlmock) {
				p := testhelpers.ProvincesDummyMap[2]
				mock.ExpectExec("^"+queryProvinceCreate).
					WithArgs(p.Name, p.CountryId).
					WillReturnError(errors.New("db error"))
			},
			arg:            testhelpers.ProvincesDummyMap[2],
			wantErr:        true,
			expectedErrMsg: "failed to create province",
		},
		{
			name: "LastInsertId error",
			setup: func(mock sqlmock.Sqlmock) {
				p := testhelpers.ProvincesDummyMap[3]
				mock.ExpectExec("^"+queryProvinceCreate).
					WithArgs(p.Name, p.CountryId).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert error")))
			},
			arg:            testhelpers.ProvincesDummyMap[3],
			wantErr:        true,
			expectedErrMsg: "failed to get province ID after creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewGeographyRepository(db)

			got, err := repo.CreateProvince(context.Background(), db, tt.arg)
			if !tt.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotZero(t, got.Id)
				require.Equal(t, tt.arg.Name, got.Name)
				require.Equal(t, tt.arg.CountryId, got.CountryId)
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

func TestGeographyRepository_FindProvinceByName(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(sqlmock.Sqlmock)
		argName        string
		argCountryId   int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				d := testhelpers.ProvincesDummyMap[1]
				mock.ExpectQuery(queryProvinceFindById).
					WithArgs(d.Name, d.CountryId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "country_id"}).
						AddRow(d.Id, d.Name, d.CountryId))
			},
			argName:      testhelpers.ProvincesDummyMap[1].Name,
			argCountryId: testhelpers.ProvincesDummyMap[1].CountryId,
			wantErr:      false,
		},
		{
			name: "not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryProvinceFindById).
					WithArgs("NoExiste", 99).
					WillReturnError(sql.ErrNoRows)
			},
			argName:        "NoExiste",
			argCountryId:   99,
			wantErr:        true,
			expectedErrMsg: "province not found",
		},
		{
			name: "other sql error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryProvinceFindById).
					WithArgs("Mendoza", 1).
					WillReturnError(errors.New("db down"))
			},
			argName:        "Mendoza",
			argCountryId:   1,
			wantErr:        true,
			expectedErrMsg: "failed to find province",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewGeographyRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			got, err := repo.FindProvinceByName(context.Background(), tt.argName, tt.argCountryId)
			if !tt.wantErr {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.argName, got.Name)
				require.Equal(t, tt.argCountryId, got.CountryId)
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
