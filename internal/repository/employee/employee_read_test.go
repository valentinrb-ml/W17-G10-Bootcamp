package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestEmployeeRepository_Read(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(sqlmock.Sqlmock, *models.Employee)
		method    string // "FindByID" o "FindByCardNumberID"
		id        int
		cardID    string
		wantNil   bool
		expectErr bool
	}{
		{
			name:   "byID_ok",
			method: "FindByID",
			id:     1,
			setup: func(mock sqlmock.Sqlmock, expected *models.Employee) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"}).
					AddRow(expected.ID, expected.CardNumberID, expected.FirstName, expected.LastName, expected.WarehouseID)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?").
					WithArgs(expected.ID).
					WillReturnRows(rows)
			},
			wantNil:   false,
			expectErr: false,
		},
		{
			name:   "byID_not_found",
			method: "FindByID",
			id:     20,
			setup: func(mock sqlmock.Sqlmock, _ *models.Employee) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?").
					WithArgs(20).
					WillReturnError(sql.ErrNoRows)
			},
			wantNil:   true,
			expectErr: false,
		},
		{
			name:   "byID_query_err",
			method: "FindByID",
			id:     21,
			setup: func(mock sqlmock.Sqlmock, _ *models.Employee) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id=?").
					WithArgs(21).
					WillReturnError(sql.ErrConnDone)
			},
			wantNil:   true,
			expectErr: true,
		},
		// ByCardNumberID tests
		{
			name:   "byCard_ok",
			method: "FindByCardNumberID",
			cardID: "EMP001",
			setup: func(mock sqlmock.Sqlmock, expected *models.Employee) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"}).
					AddRow(expected.ID, expected.CardNumberID, expected.FirstName, expected.LastName, expected.WarehouseID)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?").
					WithArgs(expected.CardNumberID).
					WillReturnRows(rows)
			},
			wantNil:   false,
			expectErr: false,
		},
		{
			name:   "byCard_not_found",
			method: "FindByCardNumberID",
			cardID: "QX",
			setup: func(mock sqlmock.Sqlmock, _ *models.Employee) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?").
					WithArgs("QX").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil:   true,
			expectErr: false,
		},
		{
			name:   "byCard_query_err",
			method: "FindByCardNumberID",
			cardID: "ERRCARD",
			setup: func(mock sqlmock.Sqlmock, _ *models.Employee) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees WHERE id_card_number=?").
					WithArgs("ERRCARD").
					WillReturnError(sql.ErrConnDone)
			},
			wantNil:   true,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, db := testhelpers.CreateMockDB()
			defer db.Close()

			var expected *models.Employee
			// Para los casos exitosos usa el helper según id o card_id
			if !tc.wantNil {
				if tc.method == "FindByID" {
					expected = testhelpers.CreateExpectedEmployee(tc.id)
				} else {
					expected = testhelpers.CreateExpectedEmployee(1) // o busca el id con el cardID
					expected.CardNumberID = tc.cardID
				}
			}
			tc.setup(mock, expected)
			repo := repo.NewEmployeeRepository(db)
			ctx := context.Background()

			switch tc.method {
			case "FindByID":
				emp, err := repo.FindByID(ctx, tc.id)
				if tc.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					if tc.wantNil {
						require.Nil(t, emp)
					} else {
						require.Equal(t, expected.FirstName, emp.FirstName)
						require.Equal(t, expected.CardNumberID, emp.CardNumberID)
					}
				}
			case "FindByCardNumberID":
				emp, err := repo.FindByCardNumberID(ctx, tc.cardID)
				if tc.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					if tc.wantNil {
						require.Nil(t, emp)
					} else {
						require.Equal(t, expected.FirstName, emp.FirstName)
						require.Equal(t, expected.CardNumberID, emp.CardNumberID)
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
		setup     func(sqlmock.Sqlmock, []models.Employee)
		wantLen   int
		expectErr bool
	}{
		{
			name: "findall_ok",
			setup: func(mock sqlmock.Sqlmock, expected []models.Employee) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"})
				for _, e := range expected {
					rows.AddRow(e.ID, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID)
				}
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees").
					WillReturnRows(rows)
			},
			wantLen:   2,
			expectErr: false,
		},
		{
			name: "findall_error_query",
			setup: func(mock sqlmock.Sqlmock, _ []models.Employee) {
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees").
					WillReturnError(sql.ErrConnDone)
			},
			wantLen:   0,
			expectErr: true,
		},
		{
			name: "findall_scan_error",
			setup: func(mock sqlmock.Sqlmock, _ []models.Employee) {
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
			setup: func(mock sqlmock.Sqlmock, expected []models.Employee) {
				rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "wareHouse_id"})
				for _, e := range expected {
					rows.AddRow(e.ID, e.CardNumberID, e.FirstName, e.LastName, e.WarehouseID)
				}
				// Aplica un error en la primera fila
				rows.RowError(0, sql.ErrConnDone)
				mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, wareHouse_id FROM employees").
					WillReturnRows(rows)
			},
			wantLen:   0,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, db := testhelpers.CreateMockDB()
			defer db.Close()

			expectedEmployees := testhelpers.CreateTestEmployees() // Usa el helper para datos
			tc.setup(mock, expectedEmployees)
			repo := repo.NewEmployeeRepository(db)
			ctx := context.Background()
			list, err := repo.FindAll(ctx)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, list, tc.wantLen)
				if tc.wantLen > 0 {
					require.Equal(t, expectedEmployees[0].FirstName, list[0].FirstName)
				}
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
