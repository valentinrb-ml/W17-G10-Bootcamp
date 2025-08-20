package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/section"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

func TestSectionRepository_UpdateSection(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type input struct {
		id  int
		sec *models.Section
	}
	type output struct {
		expected      *models.Section
		expectedError bool
		err           error
	}

	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	sec := testhelpers.DummySection(1)

	testCases := []testCase{
		{
			name: "success: section updated",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(`^UPDATE sections SET section_number = \?, current_capacity = \?, current_temperature = \? , maximum_capacity = \?, minimum_capacity = \?, minimum_temperature = \?, product_type_id = \?, warehouse_id = \?, updated_at = NOW\(\) WHERE id = \?$`).
						WithArgs(
							sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
							sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
							sec.ProductTypeId, sec.WarehouseId, 123,
						).
						WillReturnResult(sqlmock.NewResult(0, 1))
				},
			},
			input:  input{id: 123, sec: &sec},
			output: output{expected: &sec, expectedError: false, err: nil},
		},
		{
			name: "not found: section does not exist",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(`^UPDATE sections SET section_number = \?, current_capacity = \?, current_temperature = \? , maximum_capacity = \?, minimum_capacity = \?, minimum_temperature = \?, product_type_id = \?, warehouse_id = \?, updated_at = NOW\(\) WHERE id = \?$`).
						WithArgs(
							sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
							sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
							sec.ProductTypeId, sec.WarehouseId, 999,
						).
						WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
				},
			},
			input: input{id: 999, sec: &sec},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeNotFound, "The section you are trying to update does not exist.")},
		},
		{
			name: "unique constraint error (duplicate section_number)",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					myErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
					m.ExpectExec("UPDATE sections SET .*").WithArgs(
						sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
						sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
						sec.ProductTypeId, sec.WarehouseId, 124,
					).WillReturnError(myErr)
				},
			},
			input: input{id: 124, sec: &sec},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeConflict, "Section number already exists."),
			},
		},
		{
			name: "foreign key constraint error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					myErr := &mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row"}
					m.ExpectExec("UPDATE sections SET .*").WithArgs(
						sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
						sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
						sec.ProductTypeId, sec.WarehouseId, 125,
					).WillReturnError(myErr)
				},
			},
			input: input{id: 125, sec: &sec},
			output: output{
				expected: nil, expectedError: true,
				err: apperrors.NewAppError(apperrors.CodeBadRequest, "Warehouse id or product type id does not exist.")},
		},
		{
			name: "other db error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec("UPDATE sections SET .*").WithArgs(
						sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
						sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
						sec.ProductTypeId, sec.WarehouseId, 126,
					).WillReturnError(errors.New("unknown db error"))
				},
			},
			input: input{id: 126, sec: &sec},
			output: output{
				expected: nil, expectedError: true,
				err: apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while updating the section.")},
		},
		{
			name: "error in RowsAffected",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec("UPDATE sections SET .*").WithArgs(
						sec.SectionNumber, sec.CurrentCapacity, sec.CurrentTemperature,
						sec.MaximumCapacity, sec.MinimumCapacity, sec.MinimumTemperature,
						sec.ProductTypeId, sec.WarehouseId, 127,
					).WillReturnResult(sqlmock.NewErrorResult(errors.New("rowsAffected error")))
				},
			},
			input: input{id: 127, sec: &sec},
			output: output{
				expected: nil, expectedError: true,
				err: apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while updating the section.")},
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
			result, err := repo.UpdateSection(context.Background(), tc.input.id, tc.input.sec)

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
