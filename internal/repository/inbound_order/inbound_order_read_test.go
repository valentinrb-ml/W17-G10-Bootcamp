package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/inbound_order"
)

// Test unitario de ReportAll: repote agrupado de todos los empleados
func TestInboundOrderRepository_ReportAll(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantEmpty bool
		expectErr bool
	}{
		{
			name: "report_ok",
			// Devuelve dos filas, simulando dos empleados con inbound_orders
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(
					[]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id", "inbound_orders_count"},
				).
					AddRow(1, "CARDID", "Juan", "Tester", 1, 3).
					AddRow(2, "CARDID2", "Alma", "Otro", 2, 1)
				mock.ExpectQuery("SELECT e.id, e.id_card_number,").WillReturnRows(rows)
			},
			wantEmpty: false,
			expectErr: false,
		},
		{
			name: "report_empty",
			// Devuelve resultset vacío (sin empleados)
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(
					[]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id", "inbound_orders_count"},
				)
				mock.ExpectQuery("SELECT e.id, e.id_card_number,").WillReturnRows(rows)
			},
			wantEmpty: true,
			expectErr: false,
		},
		{
			name: "db_query_error",
			// Simula error de ejecución de query
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT e.id, e.id_card_number,").WillReturnError(errors.New("fail"))
			},
			wantEmpty: true,
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Crea un sqlmock y una db falsa
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			tc.mockSetup(mock)

			repository := repo.NewInboundOrderRepository(db)
			res, err := repository.ReportAll(context.Background())
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.wantEmpty {
					require.Empty(t, res)
				} else {
					require.NotEmpty(t, res)
				}
			}
			// Verifica que no quedó ninguna query/fila pendiente en el mock
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Test unitario para ReportByID: reporte de inbound orders de un solo empleado
func TestInboundOrderRepository_ReportByID(t *testing.T) {
	testCases := []struct {
		name       string
		employeeID int
		mockSetup  func(sqlmock.Sqlmock)
		wantNil    bool
		expectErr  bool
	}{
		{
			name:       "report_by_id_ok",
			employeeID: 1,
			// Devuelve datos para el empleado 1
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT e.id, e.id_card_number,").
					WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{
							"id", "id_card_number", "first_name", "last_name", "warehouse_id", "inbound_orders_count"}).
							AddRow(1, "CARDID", "Juan", "Tester", 1, 9))
			},
			wantNil:   false,
			expectErr: false,
		},
		{
			name:       "not_found",
			employeeID: 2,
			// Simula no encontrar ese empleado (sql.ErrNoRows)
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT e.id, e.id_card_number,").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			wantNil:   true,
			expectErr: true,
		},
		{
			name:       "scan_error",
			employeeID: 3,
			// Simula un error genérico
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT e.id, e.id_card_number,").
					WithArgs(3).
					WillReturnError(errors.New("fail"))
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
			repository := repo.NewInboundOrderRepository(db)
			res, err := repository.ReportByID(context.Background(), tc.employeeID)
			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
