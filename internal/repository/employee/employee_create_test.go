package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

func TestEmployeeRepository_Create(t *testing.T) {
	testCases := []struct {
		name      string
		input     *models.Employee
		mockSetup func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "inserta_ok",
			input: &models.Employee{
				CardNumberID: "C01", FirstName: "Lucas", LastName: "Test", WarehouseID: 5,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs("C01", "Lucas", "Test", 5).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectErr: false,
		},
		{
			name: "error_db_exec",
			input: &models.Employee{
				CardNumberID: "ERR", FirstName: "Lucas", LastName: "Test", WarehouseID: 5,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs("ERR", "Lucas", "Test", 5).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			name: "error_last_insert_id",
			input: &models.Employee{
				CardNumberID: "ID2", FirstName: "Mario", LastName: "Rojo", WarehouseID: 9,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs("ID2", "Mario", "Rojo", 9).
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))
			},
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
			emp, err := repo.Create(ctx, tc.input)

			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, emp)
			} else {
				require.NoError(t, err)
				require.NotZero(t, emp.ID)
				require.Equal(t, tc.input.FirstName, emp.FirstName)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
