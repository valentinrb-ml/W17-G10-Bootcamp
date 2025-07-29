package repository_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestSectionRepository_CreateSection(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type input struct {
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

	inputSec := testhelpers.DummySection(1)

	expSec := testhelpers.DummySection(1)

	testCases := []testCase{
		{
			name: "create a new section",
			arrange: arrange{
				func(m sqlmock.Sqlmock) {
					m.ExpectExec(`^INSERT INTO sections .*`).
						WithArgs(
							inputSec.SectionNumber, inputSec.CurrentCapacity, inputSec.CurrentTemperature,
							inputSec.MaximumCapacity, inputSec.MinimumCapacity, inputSec.MinimumTemperature,
							inputSec.ProductTypeId, inputSec.WarehouseId,
						).WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			input:  input{sec: &inputSec},
			output: output{expectedError: false, expected: &expSec, err: nil},
		},
		{
			name: "foreign key constraint error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					myErr := &mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row"}
					m.ExpectExec(`^INSERT INTO sections .*`).WithArgs(
						inputSec.SectionNumber, inputSec.CurrentCapacity, inputSec.CurrentTemperature,
						inputSec.MaximumCapacity, inputSec.MinimumCapacity, inputSec.MinimumTemperature,
						inputSec.ProductTypeId, inputSec.WarehouseId,
					).WillReturnError(myErr)
				},
			},
			input: input{&inputSec},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeBadRequest, "Warehouse id or product type id does not exist."),
			},
		},
		{
			name: "unique constraint error (duplicate section_number)",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					myErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
					m.ExpectExec(`^INSERT INTO sections .*`).WithArgs(
						inputSec.SectionNumber, inputSec.CurrentCapacity, inputSec.CurrentTemperature,
						inputSec.MaximumCapacity, inputSec.MinimumCapacity, inputSec.MinimumTemperature,
						inputSec.ProductTypeId, inputSec.WarehouseId,
					).WillReturnError(myErr)
				},
			},
			input: input{&inputSec},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeConflict, "Section number already exists."),
			},
		},

		{
			name: "other db error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(`^INSERT INTO sections .*`).WithArgs(
						inputSec.SectionNumber, inputSec.CurrentCapacity, inputSec.CurrentTemperature,
						inputSec.MaximumCapacity, inputSec.MinimumCapacity, inputSec.MinimumTemperature,
						inputSec.ProductTypeId, inputSec.WarehouseId,
					).WillReturnError(errors.New("unknown db error"))
				},
			},
			input: input{&inputSec},
			output: output{
				expected: nil, expectedError: true,
				err: apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while creating the section."),
			},
		},

		{
			name: "error on LastInsertId",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(`INSERT INTO sections`).WithArgs(
						inputSec.SectionNumber, inputSec.CurrentCapacity, inputSec.CurrentTemperature,
						inputSec.MaximumCapacity, inputSec.MinimumCapacity, inputSec.MinimumTemperature,
						inputSec.ProductTypeId, inputSec.WarehouseId,
					).WillReturnResult(sqlmock.NewErrorResult(errors.New("lastInsertId error")))
				},
			},
			input: input{&inputSec},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           errors.New("lastInsertId error"),
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

			result, err := repo.CreateSection(context.Background(), *tc.input.sec)
			fmt.Println(err, result)
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
