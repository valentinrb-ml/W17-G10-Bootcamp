package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseMySQL_FindById(t *testing.T) {
	type arrange struct {
		dbMock func() (sqlmock.Sqlmock, *sql.DB)
	}
	type input struct {
		id      int
		context context.Context
	}
	type output struct {
		warehouse *warehouse.Warehouse
		err       error
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	// test cases
	testCases := []testCase{
		{
			name: "success - warehouse found",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := testhelpers.CreateMockDB()

					rows := sqlmock.NewRows([]string{
						"id", "warehouse_code", "address", "minimum_temperature", "minimum_capacity",
						"telephone", "locality_id",
					}).AddRow(1, "WH001", "123 Main St", 10.5, 1000, "555-1234", "LOC001")

					mock.ExpectQuery("SELECT (.+) FROM warehouse WHERE id = ?").
						WithArgs(1).
						WillReturnRows(rows)

					return mock, db
				},
			},
			input: input{
				id:      1,
				context: context.Background(),
			},
			output: output{
				warehouse: testhelpers.CreateExpectedWarehouse(1),
				err:       nil,
			},
		},
        {
            name: "error - warehouse not found",
            arrange: arrange{
                dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
                    mock, db := testhelpers.CreateMockDB()

                    mock.ExpectQuery("SELECT (.+) FROM warehouse WHERE id = ?").
                        WithArgs(99).
                        WillReturnError(sql.ErrNoRows)

                    return mock, db
                },
            },
            input: input{
                id:      99,
                context: context.Background(),
            },
            output: output{
                warehouse: nil,
                err:       apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found"),
            },
        },
        {
            name: "error - database error",
            arrange: arrange{
                dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
                    mock, db := testhelpers.CreateMockDB()

                    mock.ExpectQuery("SELECT (.+) FROM warehouse WHERE id = ?").
                        WithArgs(99).
                        WillReturnError(sql.ErrConnDone)

                    return mock, db
                },
            },
            input: input{
                id:      99,
                context: context.Background(),
            },
            output: output{
                warehouse: nil,
                err:       apperrors.NewAppError(apperrors.CodeInternal, "error getting warehouse"),
            },
        },
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			mock, db := tc.arrange.dbMock()
			defer db.Close()
			repo := repository.NewWarehouseRepository(db)

			// act
			result, err := repo.FindById(tc.input.context, tc.input.id)

			// assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.output.warehouse.Id, result.Id)
				require.Equal(t, tc.output.warehouse.WarehouseCode, result.WarehouseCode)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
