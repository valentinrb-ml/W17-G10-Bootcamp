package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/inbound_order"
)

// Test para ExistsByOrderNumber: verifica consulta de unicidad de número de orden
func TestInboundOrderRepository_ExistsByOrderNumber(t *testing.T) {
	testCases := []struct {
		name      string
		orderNum  string
		mockSetup func(sqlmock.Sqlmock)
		expected  bool
		expectErr bool
	}{
		{
			name:     "exists_true",
			orderNum: "INV001",
			// Simula que la query retorna un count=1 (número ya existe)
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(1\\) FROM inbound_orders WHERE order_number = ?").
					WithArgs("INV001").
					WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(1))
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:     "exists_false",
			orderNum: "INV002",
			// Simula que la query retorna un count=0 (número no existe)
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(1\\) FROM inbound_orders WHERE order_number = ?").
					WithArgs("INV002").
					WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(0))
			},
			expected:  false,
			expectErr: false,
		},
		{
			name:     "scan_error",
			orderNum: "INV003",
			// Simula un error de ejecución/scan en la query
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(1\\) FROM inbound_orders WHERE order_number = ?").
					WithArgs("INV003").
					WillReturnError(errors.New("db error"))
			},
			expected:  false,
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// sqlmock crea la "db" en memoria y un mock para interceptar queries
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock) // Setup para este test en concreto

			repository := repo.NewInboundOrderRepository(db)
			res, err := repository.ExistsByOrderNumber(context.Background(), tc.orderNum)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, res) // True/False según test
			}
			require.NoError(t, mock.ExpectationsWereMet()) // Checa que no hay queries faltantes
		})
	}
}
