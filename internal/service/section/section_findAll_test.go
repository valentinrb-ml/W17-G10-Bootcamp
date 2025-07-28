package service_test

import (
	"context"
	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/section"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/section"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestSectionDefault_FindAllSections(t *testing.T) {
	type arrange struct {
		repoMock func() *mocks.SectionRepositoryMock
	}
	type output struct {
		expected      []models.Section
		expectedError bool
		err           error
	}
	type testCase struct {
		name string
		arrange
		output
	}
	sec := []models.Section{testhelpers.DummySection(1), testhelpers.DummySection(2)}

	testCases := []testCase{
		{
			name: "returns all sections successfully",
			arrange: arrange{
				func() *mocks.SectionRepositoryMock {
					return &mocks.SectionRepositoryMock{
						FuncFindAll: func(ctx context.Context) ([]models.Section, error) {
							dummySec := []models.Section{testhelpers.DummySection(1), testhelpers.DummySection(2)}
							return dummySec, nil
						},
					}
				},
			},
			output: output{
				expected:      sec,
				expectedError: false,
				err:           nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewSectionService(tc.arrange.repoMock())

			result, err := svc.FindAllSections(context.Background())

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
