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

func TestSectionDefault_CreateSection(t *testing.T) {
	type arrange struct {
		repoMock func() *mocks.SectionRepositoryMock
	}
	type output struct {
		expected      *models.Section
		expectedError bool
		err           error
	}
	type input struct {
		sec models.Section
	}
	type testCase struct {
		name string
		arrange
		output
		input
	}
	sec := testhelpers.DummySection(1)

	testCases := []testCase{
		{
			name: "returns new section on successful creation",
			arrange: arrange{
				func() *mocks.SectionRepositoryMock {
					return &mocks.SectionRepositoryMock{
						FuncCreate: func(ctx context.Context, sec models.Section) (*models.Section, error) {
							dummySec := testhelpers.DummySection(1)
							return &dummySec, nil
						},
					}
				},
			},
			input: input{sec},
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

			result, err := svc.CreateSection(context.Background(), tc.input.sec)

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
