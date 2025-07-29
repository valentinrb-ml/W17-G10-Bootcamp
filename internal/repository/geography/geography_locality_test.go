package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

const (
	queryLocalityCreate      = "INSERT INTO localities"
	queryLocalityFindById    = "SELECT id, name, province_id FROM localities WHERE id = \\?"
	queryLocalityWithSellers = "SELECT l.id, l.name, COUNT\\(s.id\\) FROM localities l"
)

func TestGeographyRepository_CreateLocality(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(sqlmock.Sqlmock)
		arg            models.Locality
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				l := testhelpers.LocalitiesDummyMap["1900"]
				mock.ExpectExec("^"+queryLocalityCreate).
					WithArgs(l.Id, l.Name, l.ProvinceId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			arg:     testhelpers.LocalitiesDummyMap["1900"],
			wantErr: false,
		},
		{
			name: "duplicate locality (mysql 1062)",
			setup: func(mock sqlmock.Sqlmock) {
				l := testhelpers.LocalitiesDummyMap["5000"]
				mysqlErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
				mock.ExpectExec("^"+queryLocalityCreate).
					WithArgs(l.Id, l.Name, l.ProvinceId).
					WillReturnError(mysqlErr)
			},
			arg:            testhelpers.LocalitiesDummyMap["5000"],
			wantErr:        true,
			expectedErrMsg: "already exists",
		},
		{
			name: "internal db error",
			setup: func(mock sqlmock.Sqlmock) {
				l := testhelpers.LocalitiesDummyMap["2000"]
				mock.ExpectExec("^"+queryLocalityCreate).
					WithArgs(l.Id, l.Name, l.ProvinceId).
					WillReturnError(errors.New("db failure"))
			},
			arg:            testhelpers.LocalitiesDummyMap["2000"],
			wantErr:        true,
			expectedErrMsg: "failed to create locality",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewGeographyRepository(db)
			got, err := repo.CreateLocality(context.Background(), db, tt.arg)

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.arg.Id, got.Id)
				require.Equal(t, tt.arg.Name, got.Name)
				require.Equal(t, tt.arg.ProvinceId, got.ProvinceId)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGeographyRepository_FindLocalityById(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(sqlmock.Sqlmock)
		argId          string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				l := testhelpers.LocalitiesDummyMap["5501"]
				mock.ExpectQuery("^" + queryLocalityFindById).
					WithArgs(l.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "province_id"}).
						AddRow(l.Id, l.Name, l.ProvinceId))
			},
			argId:   "5501",
			wantErr: false,
		},
		{
			name: "not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^" + queryLocalityFindById).
					WithArgs("doesnotexist").
					WillReturnError(sql.ErrNoRows)
			},
			argId:          "doesnotexist",
			wantErr:        true,
			expectedErrMsg: "does not exist",
		},
		{
			name: "db error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^" + queryLocalityFindById).
					WithArgs("2000").
					WillReturnError(errors.New("db is down"))
			},
			argId:          "2000",
			wantErr:        true,
			expectedErrMsg: "failed to find locality",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewGeographyRepository(db)
			got, err := repo.FindLocalityById(context.Background(), tt.argId)
			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.argId, got.Id)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGeographyRepository_CountSellersByLocality(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		argId          string
		wantErr        bool
		expectedErrMsg string
		expectResp     *models.ResponseLocalitySellers
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^" + queryLocalityWithSellers).
					WithArgs("1900").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "count"}).
						AddRow("1900", "La Plata", 5))
			},
			argId:      "1900",
			wantErr:    false,
			expectResp: &models.ResponseLocalitySellers{LocalityId: "1900", LocalityName: "La Plata", SellersCount: 5},
		},
		{
			name: "not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^" + queryLocalityWithSellers).
					WithArgs("doesnotexist").
					WillReturnError(sql.ErrNoRows)
			},
			argId:          "doesnotexist",
			wantErr:        true,
			expectedErrMsg: "does not exist",
		},
		{
			name: "db error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^" + queryLocalityWithSellers).
					WithArgs("5000").
					WillReturnError(errors.New("db down"))
			},
			argId:          "5000",
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewGeographyRepository(db)
			got, err := repo.CountSellersByLocality(context.Background(), tt.argId)
			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.expectResp, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGeographyRepository_CountSellersGroupedByLocality(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		wantErr        bool
		expectedErrMsg string
		expectResp     []models.ResponseLocalitySellers
	}{
		{
			name: "success (multiple localities)",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "count"}).
					AddRow("1900", "La Plata", 5).
					AddRow("2000", "Rosario", 3)
				mock.ExpectQuery("^" + queryLocalityWithSellers).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectResp: []models.ResponseLocalitySellers{
				{LocalityId: "1900", LocalityName: "La Plata", SellersCount: 5},
				{LocalityId: "2000", LocalityName: "Rosario", SellersCount: 3},
			},
		},
		{
			name: "db error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^" + queryLocalityWithSellers).
					WillReturnError(errors.New("db down"))
			},
			wantErr:        true,
			expectedErrMsg: "internal server error",
		},
		{
			name: "row scan error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "count"}).
					AddRow("1900", "La Plata", "notAnInt")
				mock.ExpectQuery("^" + queryLocalityWithSellers).
					WillReturnRows(rows)
			},
			wantErr:        true,
			expectedErrMsg: "Failed to scan locality sellers count",
		},
		{
			name: "rows iteration error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "count"}).
					AddRow("1900", "La Plata", 2).
					RowError(0, errors.New("row error"))
				mock.ExpectQuery("^" + queryLocalityWithSellers).
					WillReturnRows(rows)
			},
			wantErr:        true,
			expectedErrMsg: "An error occurred while iterating over the localities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(mock)
			repo := repository.NewGeographyRepository(db)
			got, err := repo.CountSellersGroupedByLocality(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErrMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.expectResp, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
