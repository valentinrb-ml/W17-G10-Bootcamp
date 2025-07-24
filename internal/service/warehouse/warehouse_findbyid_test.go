package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mocks"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseDefault_FindById(t *testing.T) {
    type arrange struct {
        mockRepo func() *mocks.WarehouseRepositoryMock
    }
    type input struct {
        id      int
        context context.Context
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
            name: "success - warehouse found",
            arrange: arrange{
                mockRepo: func() *mocks.WarehouseRepositoryMock {
                    mock := &mocks.WarehouseRepositoryMock{}
                    
                    mock.FuncFindById = func(ctx context.Context, id int) (*warehouse.Warehouse, error) {
                        return &warehouse.Warehouse{
                            Id:                 1,
                            WarehouseCode:     "WH-001",
                            Address:           "Test Address",
                            Telephone:         "123456789",
                            MinimumCapacity:   100,
                            MinimumTemperature: -10.5,
                            LocalityId:        "1900",
                        }, nil
                    }
                    
                    return mock
                },
            },
            input: input{
                id:      1,
                context: context.Background(),
            },
            output: output{
                result: &warehouse.Warehouse{
                    Id:                 1,
                    WarehouseCode:     "WH-001",
                    Address:           "Test Address",
                    Telephone:         "123456789",
                    MinimumCapacity:   100,
                    MinimumTemperature: -10.5,
                    LocalityId:        "1900",
                },
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
            result, err := srv.FindById(tc.input.context, tc.input.id)

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

