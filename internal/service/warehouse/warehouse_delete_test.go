package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mocks"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"
)

func TestWarehouseDefault_Delete(t *testing.T) {
	type arrange struct {
		mockRepo func() *mocks.WarehouseRepositoryMock
	}
	type input struct {
		id      int
		context context.Context
	}
	type output struct {
		err error
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	// test cases (solo caso feliz para método pasamanos)
	testCases := []testCase{
		{
			name: "success - warehouse deleted",
			arrange: arrange{
				mockRepo: func() *mocks.WarehouseRepositoryMock {
					mock := &mocks.WarehouseRepositoryMock{}

					mock.FuncDelete = func(ctx context.Context, id int) error {
						return nil
					}

					return mock
				},
			},
			input: input{
				id:      1,
				context: context.Background(),
			},
			output: output{
				err: nil,
			},
		},
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			mockRepo := tc.arrange.mockRepo()
			srv := service.NewWarehouseService(mockRepo)

			// act
			err := srv.Delete(tc.input.context, tc.input.id)

			// assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
