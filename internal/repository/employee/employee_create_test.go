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

// Test unitario del método Create del repositorio de empleados usando helpers y sqlmock.
func TestEmployeeRepository_Create(t *testing.T) {
	testCases := []struct {
		name      string
		input     *models.Employee
		mockSetup func(mock sqlmock.Sqlmock, in models.Employee)
		expectErr bool
	}{
		{
			name: "inserta_ok",
			// Usa el helper para generar un employee dummy válido
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee() // struct
				return &e
			}(),
			// Simula un insert exitoso en la base, retorna una fila insertada con id=1
			mockSetup: func(mock sqlmock.Sqlmock, in models.Employee) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs(in.CardNumberID, in.FirstName, in.LastName, in.WarehouseID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectErr: false,
		},
		{
			name: "error_db_exec",
			// Cambia un campo para simular una situación errónea o borde
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				e.CardNumberID = "ERR"
				e.FirstName = "Lucas"
				e.LastName = "Test"
				e.WarehouseID = 5
				return &e
			}(),
			// Simula que la BD falla al ejecutar el insert
			mockSetup: func(mock sqlmock.Sqlmock, in models.Employee) {
				mock.ExpectExec("INSERT INTO employees").
					WithArgs(in.CardNumberID, in.FirstName, in.LastName, in.WarehouseID).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			name: "error_last_insert_id",
			// Cambia los campos para otro caso y simula error al obtener el id
			input: func() *models.Employee {
				e := testhelpers.CreateTestEmployee()
				e.CardNumberID = "ID2"
				e.FirstName = "Mario"
				e.LastName = "Rojo"
				e.WarehouseID = 9
				return &e
			}(),
			// El insert ocurre pero falla al obtener el LastInsertId
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
			// Crea un mock de la DB y el sqlmock que interceptará consultas
			mock, db := testhelpers.CreateMockDB()
			defer db.Close()
			// Convierte el input de puntero a struct para pasárselo al mock setup
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
				require.Equal(t, in.FirstName, emp.FirstName) // verifica el nombre creado
			}
			// Valida que todas las queries esperadas se hayan ejecutado (sqlmock)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
