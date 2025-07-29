package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"testing"
)

func TestSectionRepository_DeleteSection(t *testing.T) {
	type arrange struct {
		dbMock func(sqlmock.Sqlmock)
	}
	type input struct {
		id int
	}
	type output struct {
		expectedError bool
		err           error
	}

	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	testCases := []testCase{
		{
			name: "deletes section successfully",
			arrange: arrange{
				dbMock: func(m sqlmock.Sqlmock) {
					m.ExpectExec(`^DELETE FROM sections WHERE id =\?`).
						WithArgs(1).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			input:  input{1},
			output: output{expectedError: false, err: nil},
		},
		{
			name: "returns conflict error when section has associated batches",
			arrange: arrange{dbMock: func(m sqlmock.Sqlmock) {
				myErr := &mysql.MySQLError{Number: 1451, Message: "Cannot delete or update a parent row"}
				m.ExpectExec(`^DELETE FROM sections WHERE id =\?`).
					WithArgs(2).
					WillReturnError(myErr)
			}},
			input:  input{id: 2},
			output: output{expectedError: true, err: apperrors.NewAppError(apperrors.CodeConflict, "Cannot delete section: there are products batches associated with this section.")},
		},
		{
			name: "returns not found error when section does not exist",
			arrange: arrange{dbMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(`^DELETE FROM sections WHERE id =\?`).
					WithArgs(3).
					WillReturnResult(sqlmock.NewResult(0, 0))
			}},
			input:  input{id: 3},
			output: output{expectedError: true, err: apperrors.NewAppError(apperrors.CodeNotFound, "The section you are trying to delete does not exist.")},
		},
		{
			name: "returns internal error when ExecContext fails",
			arrange: arrange{dbMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(`^DELETE FROM sections WHERE id =\?`).
					WithArgs(4).
					WillReturnError(errors.New("db error"))
			}},
			input:  input{id: 4},
			output: output{expectedError: true, err: apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while deleting the section.")},
		},
		{
			name: "returns internal error when RowsAffected returns error",
			arrange: arrange{dbMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(`^DELETE FROM sections WHERE id =\?`).
					WithArgs(5).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rowsAffected error")))
			}},
			input:  input{id: 5},
			output: output{expectedError: true, err: apperrors.NewAppError(apperrors.CodeInternal, "An internal server error occurred while deleting the section.")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := NewSectionRepository(db)

			tc.arrange.dbMock(mock)
			err = repo.DeleteSection(context.Background(), tc.input.id)

			if tc.output.expectedError {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}

}
