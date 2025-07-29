package repository_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestSectionRepository_FindById(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type output struct {
		expected      *models.Section
		expectedError bool
		err           error
	}
	type input struct {
		id int
	}
	type testCase struct {
		name string
		input
		arrange
		output
	}
	expSec := testhelpers.DummySection(1)
	testCases := []testCase{
		{
			name: "returns sections by id",
			arrange: arrange{dbMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "section_number", "current_capacity", "current_temperature",
					"maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id",
				}).AddRow(expSec.Id, expSec.SectionNumber, expSec.CurrentCapacity, expSec.CurrentTemperature, expSec.MaximumCapacity, expSec.MinimumCapacity, expSec.MinimumTemperature, expSec.ProductTypeId, expSec.WarehouseId)
				m.ExpectQuery(`^SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections WHERE id = ?`).
					WillReturnRows(rows)
			}},
			input: input{id: 1},
			output: output{
				expected:      &expSec,
				expectedError: false,
				err:           nil,
			},
		},
		{
			name: "returns error when scan fails",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{
						"id", "section_number", "current_capacity", "current_temperature",
						"maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"},
					).AddRow("INVALID", 10, 20, 7, 100, 10, 5, 2, 3)
					m.ExpectQuery(`^SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections `).WillReturnRows(rows)
				},
			},
			input: input{id: 1},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the section."),
			},
		},
		{
			name: "returns empty array when no sections exist",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectQuery(`^SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections WHERE id = ?`).
						WillReturnRows(
							sqlmock.NewRows([]string{
								"id", "section_number", "current_capacity", "current_temperature",
								"maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id",
							}),
						)
				},
			},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeNotFound, "The section you are looking for does not exist."),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := repository.NewSectionRepository(db)

			tc.arrange.dbMock(mock)

			result, err := repo.FindById(context.Background(), tc.id)
			if tc.output.expectedError {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			require.NoError(t, err)
			require.Equal(t, tc.output.expected, result)
			require.NoError(t, mock.ExpectationsWereMet())
		})

	}
}
