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

func TestSectionDefault_FindById(t *testing.T) {
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
		name    string
		arrange arrange
		output  output
		input
	}
	sec := testhelpers.DummySection(1)

	testCases := []testCase{
		{
			name: "finds section by id successfully",
			arrange: arrange{
				func() *mocks.SectionRepositoryMock {
					return &mocks.SectionRepositoryMock{
						FuncFindById: func(ctx context.Context, id int) (*models.Section, error) {
							dummySec := testhelpers.DummySection(1)
							return &dummySec, nil
						},
					}
				},
			},
			input: input{id: 1},
			output: output{
				expected:      &sec,
				expectedError: false,
				err:           nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewSectionService(tc.arrange.repoMock())
			svc.SetLogger(testhelpers.NewTestLogger())

			result, err := svc.FindById(context.Background(), tc.input.id)

			if tc.output.expectedError {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.output.expected, result)
		})
	}
}
