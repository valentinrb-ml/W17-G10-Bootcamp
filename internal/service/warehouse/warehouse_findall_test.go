package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mocks"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/testhelpers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseDefault_FindAll(t *testing.T) {
	type arrange struct {
		mockRepo func() *mocks.WarehouseRepositoryMock
	}
	type input struct {
		context context.Context
	}
	type output struct {
		result []warehouse.Warehouse
		err    error
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
			name: "success - warehouses found",
			arrange: arrange{
				mockRepo: func() *mocks.WarehouseRepositoryMock {
					mock := &mocks.WarehouseRepositoryMock{}

					mock.FuncFindAll = func(ctx context.Context) ([]warehouse.Warehouse, error) {
						return testhelpers.CreateTestWarehouses(), nil
					}

					return mock
				},
			},
			input: input{
				context: context.Background(),
			},
			output: output{
				result: testhelpers.CreateTestWarehouses(),
				err:    nil,
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
			result, err := srv.FindAll(tc.input.context)

			// assert
			if tc.output.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.output.err.Error(), err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output.result, result)
			}
		})
	}
}
