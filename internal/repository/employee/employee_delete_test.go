package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Test unitario del método Delete del repositorio de employees usando sqlmock.
func TestEmployeeRepository_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		inputID   int
		mockSetup func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:    "delete_ok",
			inputID: 10,
			// Simula que la query DELETE ejecuta correctamente y borra 1 fila
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(10).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 fila borrada
			},
			expectErr: false,
		},
		{
			name:    "delete_db_error",
			inputID: 11,
			// Simula que la BD falla al ejecutar el delete
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(11).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			name:    "no_rows_affected",
			inputID: 12,
			// Simula que el DELETE ejecuta pero no borra ninguna fila (empleado no existe)
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(12).
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 filas borradas
			},
			expectErr: true,
		},
		{
			name:    "rows_affected_fails",
			inputID: 13,
			// Simula que hay un error al preguntar cuántas filas fueron afectadas (poco común)
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM employees").
					WithArgs(13).
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Usa el helper para armar DB/mock (DRY)
			mock, db := testhelpers.CreateMockDB()
			defer db.Close()
			// Configura el escenario de query esperado para el caso
			tc.mockSetup(mock)
			repo := repo.NewEmployeeRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())
			ctx := context.Background()
			// Ejecuta el método Delete
			err := repo.Delete(ctx, tc.inputID)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			// Valida que no quedó ningún EXPECT sin ejecutar en el mock (verifica cobertura total de queries)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
