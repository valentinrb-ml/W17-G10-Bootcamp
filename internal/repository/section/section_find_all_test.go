package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSectionRepository_FindAllSections(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type output struct {
		expected      []models.Section
		expectedError bool
		err           error
	}
	type testCase struct {
		name    string
		arrange arrange
		output  output
	}
	expSec := testhelpers.DummySection(1)
	testCases := []testCase{
		{
			name: "returns all sections on success",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{
						"id", "section_number", "current_capacity", "current_temperature",
						"maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id",
					}).AddRow(expSec.Id, expSec.SectionNumber, expSec.CurrentCapacity, expSec.CurrentTemperature, expSec.MaximumCapacity, expSec.MinimumCapacity, expSec.MinimumTemperature, expSec.ProductTypeId, expSec.WarehouseId)
					m.ExpectQuery(`^SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections `).
						WillReturnRows(rows)
				},
			},
			output: output{expected: []models.Section{expSec},
				expectedError: false,
				err:           nil,
			},
		},
		{
			name: "returns error when db fails",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectQuery(`^SELECT .* FROM sections`).WillReturnError(errors.New("connection error"))
				},
			},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the sections."),
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
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the section."),
			},
		}, {
			name: "returns error if rows.Err is not nil",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{
						"id", "section_number", "current_capacity", "current_temperature",
						"maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id",
					}).
						AddRow(1, 1, 1, 1, 1, 1, 1, 1, 1).
						RowError(0, fmt.Errorf("error row"))
					m.ExpectQuery(`^SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections\s*$`).
						WillReturnRows(rows)
				},
			},
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
					m.ExpectQuery(`^SELECT id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id FROM sections\s*$`).
						WillReturnRows(
							sqlmock.NewRows([]string{
								"id", "section_number", "current_capacity", "current_temperature",
								"maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id",
							}),
						)
				},
			},
			output: output{
				expected:      []models.Section{},
				expectedError: false,
				err:           nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := repository.NewSectionRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			tc.arrange.dbMock(mock)

			result, err := repo.FindAllSections(context.Background())
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
