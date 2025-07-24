package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
)

func TestEmployeeRepository_FindByID_and_ByCardNumberID(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		method    string // "byID" o "byCardNumber"
		id        int
		cardID    string
		wantNil   bool
		expectErr bool
	}{
		{
			name:   "byID_ok",
			method: "byID",
			id:     1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"}).
					AddRow(1, "C01", "Test", "Name", 91)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantNil:   false,
			expectErr: false,
		},
		{
			name:   "byID_not_found",
			method: "byID",
			id:     7,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?").
					WithArgs(7).
					WillReturnError(sql.ErrNoRows)
			},
			wantNil:   true,
			expectErr: false,
		},
		{
			name:   "byID_query_err",
			method: "byID",
			id:     9,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?").
					WithArgs(9).
					WillReturnError(sql.ErrConnDone)
			},
			wantNil:   true,
			expectErr: true,
		},
		// ByCardNumberID tests
		{
			name:   "byCard_ok",
			method: "byCardNumber",
			cardID: "C99",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"}).
					AddRow(2, "C99", "Eri", "Rey", 11)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?").
					WithArgs("C99").
					WillReturnRows(rows)
			},
			wantNil:   false,
			expectErr: false,
		},
		{
			name:   "byCard_not_found",
			method: "byCardNumber",
			cardID: "QX",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?").
					WithArgs("QX").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil:   true,
			expectErr: false,
		},
		{
			name:   "byCard_query_err",
			method: "byCardNumber",
			cardID: "EX",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?").
					WithArgs("EX").
					WillReturnError(sql.ErrConnDone)
			},
			wantNil:   true,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)
			repo := repo.NewEmployeeRepository(db)
			ctx := context.Background()

			switch tc.method {
			case "byID":
				e, err := repo.FindByID(ctx, tc.id)
				if tc.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					if tc.wantNil {
						require.Nil(t, e)
					} else {
						require.NotNil(t, e)
					}
				}
			case "byCardNumber":
				e, err := repo.FindByCardNumberID(ctx, tc.cardID)
				if tc.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					if tc.wantNil {
						require.Nil(t, e)
					} else {
						require.NotNil(t, e)
					}
				}
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmployeeRepository_FindAll(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantLen   int
		expectErr bool
	}{
		{
			name: "findall_ok",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"}).
					AddRow(1, "A1", "Test", "AA", 1).
					AddRow(2, "B2", "Test", "BB", 2)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees").
					WillReturnRows(rows)
			},
			wantLen:   2,
			expectErr: false,
		},
		{
			name: "findall_error_query",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees").
					WillReturnError(sql.ErrConnDone)
			},
			wantLen:   0,
			expectErr: true,
		},
		{
			name: "findall_scan_error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"}).
					AddRow(nil, "B2", "Test", "BB", 2)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees").
					WillReturnRows(rows)
			},
			wantLen:   0,
			expectErr: true,
		},
		{
			name: "findall_row_iter_error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"}).
					AddRow(1, "X", "Y", "Z", 5).
					RowError(0, sql.ErrConnDone)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees").
					WillReturnRows(rows)
			},
			wantLen:   0,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)
			repo := repo.NewEmployeeRepository(db)
			ctx := context.Background()
			list, err := repo.FindAll(ctx)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, list, tc.wantLen)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
