package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/section"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/section"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSectionDefault_DeleteSection(t *testing.T) {
	type arrange struct {
		repoMock func() *mocks.SectionRepositoryMock
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
		arrange
		output
		input
	}

	testCases := []testCase{
		{
			name: "deletes section successfully",
			arrange: arrange{
				func() *mocks.SectionRepositoryMock {
					return &mocks.SectionRepositoryMock{
						FuncDelete: func(ctx context.Context, id int) error {
							return nil
						},
					}
				},
			},
			input: input{1},
			output: output{
				expected:      nil,
				expectedError: false,
				err:           nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewSectionService(tc.arrange.repoMock())
			svc.SetLogger(testhelpers.NewTestLogger())

			err := svc.DeleteSection(context.Background(), tc.input.id)

			if tc.output.expectedError {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
