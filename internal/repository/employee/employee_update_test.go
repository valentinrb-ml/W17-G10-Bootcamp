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

// Test unitario del método Update del repositorio de empleados usando sqlmock y helpers DRY.
func TestEmployeeRepository_Update(t *testing.T) {
	testCases := []struct {
		name      string
		id        int
		input     *models.Employee
		mockSetup func(sqlmock.Sqlmock, models.Employee)
		expectErr bool
	}{
		{
			name: "update_ok",
			id:   22,
			// Usa el helper para crear un employee dummy y luego personaliza campos para el caso
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				e.CardNumberID = "A1"
				e.FirstName = "Daniela"
				e.LastName = "Zamora"
				e.WarehouseID = 99
				return &e
			}(),
			// Simula que la actualización en la DB es exitosa para esos campos
			mockSetup: func(mock sqlmock.Sqlmock, in models.Employee) {
				mock.ExpectExec("UPDATE employees").
					WithArgs("A1", "Daniela", "Zamora", 99, 22).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name: "update_db_error",
			id:   24,
			// Modifica los campos básicos para simular otro caso de update
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				e.CardNumberID = "B1"
				e.FirstName = "Falla"
				e.LastName = "Falla"
				e.WarehouseID = 88
				return &e
			}(),
			// Simula que la DB devuelve un error al ejecutar el update
			mockSetup: func(mock sqlmock.Sqlmock, in models.Employee) {
				mock.ExpectExec("UPDATE employees").
					WithArgs("B1", "Falla", "Falla", 88, 24).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Helper para crear conexión y mock sqlmock
			mock, db := testhelpers.CreateMockDB()
			defer db.Close()
			// Valor a usar en el setup del mock
			in := *tc.input
			tc.mockSetup(mock, in)

			repo := repo.NewEmployeeRepository(db)
			ctx := context.Background()
			// Ejecuta el método Update del repo
			err := repo.Update(ctx, tc.id, tc.input)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			// Chequea que todas las queries esperadas hayan sido ejecutadas
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
