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

func TestEmployeeRepository_Update(t *testing.T) {
	testCases := []struct {
		name      string
		id        int
		input     *models.Employee
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:  "update_ok",
			id:    22,
			input: &models.Employee{CardNumberID: "A1", FirstName: "Daniela", LastName: "Zamora", WarehouseID: 99},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE employees").
					WithArgs("A1", "Daniela", "Zamora", 99, 22).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name:  "update_db_error",
			id:    24,
			input: &models.Employee{CardNumberID: "B1", FirstName: "Falla", LastName: "Falla", WarehouseID: 88},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE employees").
					WithArgs("B1", "Falla", "Falla", 88, 24).
					WillReturnError(sql.ErrConnDone)
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
			err = repo.Update(ctx, tc.id, tc.input)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
