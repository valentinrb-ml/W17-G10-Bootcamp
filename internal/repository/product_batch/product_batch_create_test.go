package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_batch"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductBatchesRepository_CreateProductBatches(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type input struct {
		batch *models.ProductBatches
	}
	type output struct {
		expected      *models.ProductBatches
		expectedError bool
		err           error
	}

	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	inputBatch := testhelpers.DummyProductBatch(1)
	expBatch := testhelpers.DummyProductBatch(1)
	expBatch.Id = 1 // lastInsertId simulado

	const insertRegex = `^INSERT INTO product_batches .*`

	testCases := []testCase{
		{
			name: "create a new product batch",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(insertRegex).
						WithArgs(
							inputBatch.BatchNumber,
							inputBatch.CurrentQuantity,
							inputBatch.CurrentTemperature,
							inputBatch.DueDate,
							inputBatch.InitialQuantity,
							inputBatch.ManufacturingDate,
							inputBatch.ManufacturingHour,
							inputBatch.MinimumTemperature,
							inputBatch.ProductId,
							inputBatch.SectionId,
						).WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			input:  input{batch: &inputBatch},
			output: output{expectedError: false, expected: &expBatch, err: nil},
		},
		{
			name: "foreign key constraint error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					myErr := &mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row"}
					m.ExpectExec(insertRegex).
						WithArgs(
							inputBatch.BatchNumber,
							inputBatch.CurrentQuantity,
							inputBatch.CurrentTemperature,
							inputBatch.DueDate,
							inputBatch.InitialQuantity,
							inputBatch.ManufacturingDate,
							inputBatch.ManufacturingHour,
							inputBatch.MinimumTemperature,
							inputBatch.ProductId,
							inputBatch.SectionId,
						).WillReturnError(myErr)
				},
			},
			input: input{batch: &inputBatch},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeBadRequest, "Section id or product id does not exist."),
			},
		},
		{
			name: "unique constraint error (duplicate batch_number)",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					myErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
					m.ExpectExec(insertRegex).
						WithArgs(
							inputBatch.BatchNumber,
							inputBatch.CurrentQuantity,
							inputBatch.CurrentTemperature,
							inputBatch.DueDate,
							inputBatch.InitialQuantity,
							inputBatch.ManufacturingDate,
							inputBatch.ManufacturingHour,
							inputBatch.MinimumTemperature,
							inputBatch.ProductId,
							inputBatch.SectionId,
						).WillReturnError(myErr)
				},
			},
			input: input{batch: &inputBatch},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeConflict, "Batch number already exists."),
			},
		},
		{
			name: "other db error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(insertRegex).
						WithArgs(
							inputBatch.BatchNumber,
							inputBatch.CurrentQuantity,
							inputBatch.CurrentTemperature,
							inputBatch.DueDate,
							inputBatch.InitialQuantity,
							inputBatch.ManufacturingDate,
							inputBatch.ManufacturingHour,
							inputBatch.MinimumTemperature,
							inputBatch.ProductId,
							inputBatch.SectionId,
						).WillReturnError(errors.New("unknown db error"))
				},
			},
			input: input{batch: &inputBatch},
			output: output{
				expected: nil, expectedError: true,
				err: apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while creating the Product Batch."),
			},
		},
		{
			name: "error on LastInsertId",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(insertRegex).
						WithArgs(
							inputBatch.BatchNumber,
							inputBatch.CurrentQuantity,
							inputBatch.CurrentTemperature,
							inputBatch.DueDate,
							inputBatch.InitialQuantity,
							inputBatch.ManufacturingDate,
							inputBatch.ManufacturingHour,
							inputBatch.MinimumTemperature,
							inputBatch.ProductId,
							inputBatch.SectionId,
						).WillReturnResult(sqlmock.NewErrorResult(errors.New("lastInsertId error")))
				},
			},
			input: input{batch: &inputBatch},
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
			repo := repository.NewProductBatchesRepository(db)
			repo.SetLogger(testhelpers.NewTestLogger())

			tc.arrange.dbMock(mock)

			result, err := repo.CreateProductBatches(context.Background(), *tc.input.batch)
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
