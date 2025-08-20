package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestWarehouseDefault_Create(t *testing.T) {
	type arrange struct {
		mockRepo func() *mocks.WarehouseRepositoryMock
	}
	type input struct {
		warehouse warehouse.Warehouse
		context   context.Context
	}
	type output struct {
		result *warehouse.Warehouse
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
			name: "success - warehouse created",
			arrange: arrange{
				mockRepo: func() *mocks.WarehouseRepositoryMock {
					mock := &mocks.WarehouseRepositoryMock{}

					mock.FuncCreate = func(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
						return testhelpers.CreateExpectedWarehouse(1), nil
					}

					return mock
				},
			},
			input: input{
				warehouse: testhelpers.CreateTestWarehouse(),
				context:   context.Background(),
			},
			output: output{
				result: testhelpers.CreateExpectedWarehouse(1),
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
			srv.SetLogger(testhelpers.NewTestLogger())

			// act
			result, err := srv.Create(tc.input.context, tc.input.warehouse)

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
