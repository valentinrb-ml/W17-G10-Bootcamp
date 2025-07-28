package service_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/section"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/section"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestSectionDefault_UpdateSection(t *testing.T) {
	type arrange struct {
		repoMock func() *mocks.SectionRepositoryMock
	}
	type output struct {
		expected      *models.Section
		expectedError bool
		err           error
	}
	type input struct {
		sec models.PatchSection
		id  int
	}
	type testCase struct {
		name string
		arrange
		output
		input
	}
	sec := testhelpers.DummySection(1)
	dumSec := testhelpers.DummySectionPatch(1)

	testCases := []testCase{
		{
			name: "updates section successfully",
			arrange: arrange{
				func() *mocks.SectionRepositoryMock {
					return &mocks.SectionRepositoryMock{
						FuncFindById: func(ctx context.Context, id int) (*models.Section, error) {
							dummySec := testhelpers.DummySection(1)
							return &dummySec, nil
						},
						FuncUpdate: func(ctx context.Context, id int, sec *models.Section) (*models.Section, error) {
							dummySec := testhelpers.DummySection(1)
							return &dummySec, nil
						},
					}
				},
			},
			input: input{sec: dumSec, id: 1},
			output: output{
				expected:      &sec,
				expectedError: false,
				err:           nil,
			},
		},
		{
			name: "error: not found on FindById",
			arrange: arrange{
				repoMock: func() *mocks.SectionRepositoryMock {
					return &mocks.SectionRepositoryMock{
						FuncFindById: func(ctx context.Context, id int) (*models.Section, error) {
							return nil, errors.New("section not found")
						},
						FuncUpdate: nil,
					}
				},
			},
			input: input{sec: dumSec, id: 2},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           errors.New("section not found"),
			},
		},
		{
			name: "error: update fails",
			arrange: arrange{
				repoMock: func() *mocks.SectionRepositoryMock {
					return &mocks.SectionRepositoryMock{
						FuncFindById: func(ctx context.Context, id int) (*models.Section, error) {
							dummySec := testhelpers.DummySection(1)
							return &dummySec, nil
						},
						FuncUpdate: func(ctx context.Context, id int, sec *models.Section) (*models.Section, error) {
							return nil, errors.New("update error")
						},
					}
				},
			},
			input: input{sec: dumSec, id: 1},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           errors.New("update error"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewSectionService(tc.arrange.repoMock())

			result, err := svc.UpdateSection(context.Background(), tc.input.id, tc.input.sec)
			fmt.Println(err)

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
