package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	repo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_batch"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestProductBatchesRepository_GetReportProductById(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type input struct {
		id int
	}
	type output struct {
		expected      *models.ReportProduct
		expectedError bool
		err           error
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	dummy := testhelpers.DummyReportProduct()

	const query = `SELECT s.id, s.section_number, SUM\(p.current_quantity\) FROM product_batches p INNER JOIN sections s on p.section_id = s.id\s+WHERE p.section_id = \?\s+GROUP BY p.section_id`

	testCases := []testCase{
		{
			name: "success - returns product report",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{
						"id", "section_number", "sum",
					}).AddRow(dummy.SectionId, dummy.SectionNumber, dummy.ProductsCount)
					m.ExpectQuery(query).WithArgs(10).WillReturnRows(rows)
				},
			},
			input: input{id: 10},
			output: output{
				expected:      &dummy,
				expectedError: false,
				err:           nil,
			},
		},
		{
			name: "not found - returns custom not found error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectQuery(query).WithArgs(99).WillReturnError(sql.ErrNoRows)
				},
			},
			input: input{id: 99},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeNotFound, "The section you are looking for does not exist."),
			},
		},
		{
			name: "internal error - db error",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectQuery(query).WithArgs(20).WillReturnError(errors.New("db error"))
				},
			},
			input: input{id: 20},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while retrieving the product report."),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repository := repo.NewProductBatchesRepository(db)

			tc.arrange.dbMock(mock)

			result, err := repository.GetReportProductById(context.Background(), tc.input.id)

			if tc.output.expectedError {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output.expected, result)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
