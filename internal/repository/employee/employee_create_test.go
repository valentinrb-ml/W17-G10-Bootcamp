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

func TestEmployeeRepository_Create(t *testing.T) {
	testCases := []struct {
		name      string
		input     *models.Employee
		mockSetup func(mock sqlmock.Sqlmock, in models.Employee)
		expectErr bool
	}{
		{
			name: "inserta_ok",
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee() // struct
				return &e
			}(),
			mockSetup: func(mock sqlmock.Sqlmock, in models.Employee) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs(in.CardNumberID, in.FirstName, in.LastName, in.WarehouseID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectErr: false,
		},
		{
			name: "error_db_exec",
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				e.CardNumberID = "ERR"
				e.FirstName = "Lucas"
				e.LastName = "Test"
				e.WarehouseID = 5
				return &e
			}(),
			mockSetup: func(mock sqlmock.Sqlmock, in models.Employee) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs(in.CardNumberID, in.FirstName, in.LastName, in.WarehouseID).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			name: "error_last_insert_id",
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				e.CardNumberID = "ID2"
				e.FirstName = "Mario"
				e.LastName = "Rojo"
				e.WarehouseID = 9
				return &e
			}(),
			mockSetup: func(mock sqlmock.Sqlmock, in models.Employee) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs(in.CardNumberID, in.FirstName, in.LastName, in.WarehouseID).
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, db := testhelpers.CreateMockDB()
			defer db.Close()

			in := *tc.input // struct para mock
			tc.mockSetup(mock, in)

			repo := repo.NewEmployeeRepository(db)
			ctx := context.Background()
			emp, err := repo.Create(ctx, tc.input)

			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, emp)
			} else {
				require.NoError(t, err)
				require.NotZero(t, emp.ID)
				require.Equal(t, in.FirstName, emp.FirstName)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
