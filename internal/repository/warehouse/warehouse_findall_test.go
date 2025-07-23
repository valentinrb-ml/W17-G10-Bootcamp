package repository_test

import (
    "context"
    "database/sql"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/require"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseMySQL_FindAll(t *testing.T) {
	type arrange struct {
        dbMock func() (sqlmock.Sqlmock, *sql.DB)
    }
    type input struct {
        context context.Context
    }
    type output struct {
        warehouses []warehouse.Warehouse
        err        error
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
			name: "success - warehouses found",
			arrange: arrange{
				dbMock: func() (sqlmock.Sqlmock, *sql.DB) {
					mock, db := createMockDB()

					rows := sqlmock.NewRows([]string{
						"id", "warehouse_code", "address", "minimum_temperature", "minimum_capacity",
						"telephone", "locality_id",
					}).
						AddRow(1, "WH001", "123 Main St", 10.5, 1000, "555-1234", "LOC001").
						AddRow(2, "WH002", "456 Elm St", 15.5, 2000, "555-5678", "LOC002")

					mock.ExpectQuery("SELECT (.+) FROM warehouse").
						WillReturnRows(rows)

					return mock, db
				},
			},
			input: input {
				context: context.Background(),
			},
			output: output {
				warehouses: createExpectedWarehouses(),
				err:       nil,
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
            result, err := repo.FindAll(tc.input.context)

            // assert
            if tc.output.err != nil {
                require.Error(t, err)
                require.Equal(t, tc.output.err.Error(), err.Error())
                require.Nil(t, result)
            } else {
                require.NoError(t, err)
                require.Equal(t, len(tc.output.warehouses), len(result))
                if len(tc.output.warehouses) > 0 {
                    require.Equal(t, tc.output.warehouses[0].Id, result[0].Id)
                    require.Equal(t, tc.output.warehouses[0].WarehouseCode, result[0].WarehouseCode)
                }
            }
            require.NoError(t, mock.ExpectationsWereMet())
        })
    }
}
