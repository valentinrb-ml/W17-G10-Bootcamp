package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/inbound_order"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/inbound_order"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

// Test unitario para el método Create del repositorio de inbound_order usando sqlmock y helpers DRY
func TestInboundOrderRepository_Create(t *testing.T) {
	testCases := []struct {
		name      string
		input     *models.InboundOrder
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
		wantCode  string // para validar el tipo/código del error si aplica
	}{
		{
			name: "insert_ok",
			// Datos dummy con helper
			input: testhelpers.CreateExpectedInboundOrder(0),
			// Simula un insert exitoso con ret id 10
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO inbound_orders").
					WithArgs("2024-06-01", "INV001", 1, 10, 1).
					WillReturnResult(sqlmock.NewResult(10, 1))
			},
			expectErr: false,
		},
		{
			name: "unique_violation",
			// Simula conflicto de unique/duplicate en order_number (MySQL error 1062)
			input: testhelpers.CreateExpectedInboundOrder(0),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO inbound_orders").
					WithArgs("2024-06-01", "INV001", 1, 10, 1).
					WillReturnError(&mysql.MySQLError{Number: 1062, Message: "dup"})
			},
			expectErr: true,
			wantCode:  "CONFLICT",
		},
		{
			name: "fk_violation",
			// Simula violación de clave foránea (MySQL error 1452) en cualquier FK referenciada
			input: testhelpers.CreateExpectedInboundOrder(0),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO inbound_orders").
					WithArgs("2024-06-01", "INV001", 1, 10, 1).
					WillReturnError(&mysql.MySQLError{Number: 1452, Message: "fk fail"})
			},
			expectErr: true,
			wantCode:  "NOT_FOUND",
		},
		{
			name: "generic_error",
			// Simula un error genérico de la bd no manejado especialmente
			input: testhelpers.CreateExpectedInboundOrder(0),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO inbound_orders").
					WithArgs("2024-06-01", "INV001", 1, 10, 1).
					WillReturnError(errors.New("generic"))
			},
			expectErr: true,
		},
		{
			name: "error_last_insert_id",
			// Simula error al obtener el último insert id después del insert (mal driver)
			input: testhelpers.CreateExpectedInboundOrder(0),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO inbound_orders").
					WithArgs("2024-06-01", "INV001", 1, 10, 1).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("err")))
			},
			expectErr: true,
			wantCode:  "INTERNAL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Instancia una DB mockeada con sqlmock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			// Configura el mock según el caso
			tc.mockSetup(mock)
			repository := repo.NewInboundOrderRepository(db)
			repository.SetLogger(testhelpers.NewTestLogger())
			ctx := context.Background()
			// Ejecuta el método Create
			res, err := repository.Create(ctx, tc.input)
			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, res)
				// Si aplica, chequea el código en el mensaje de error
				if tc.wantCode != "" {
					require.Contains(t, err.Error(), tc.wantCode)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotZero(t, res.ID)
			}
			// Valida que todas las expectativas del mock se hayan cumplido
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
