package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mocks"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func TestWarehouseDefault_Update(t *testing.T) {
    type arrange struct {
        mockRepo func() *mocks.WarehouseRepositoryMock
    }
    type input struct {
        id      int
        patch   warehouse.WarehousePatchDTO
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

    // test cases (tests completos para método con lógica de negocio)
    testCases := []testCase{
        {
            name: "success - update with valid minimum capacity",
            arrange: arrange{
                mockRepo: func() *mocks.WarehouseRepositoryMock {
                    mock := &mocks.WarehouseRepositoryMock{}
                    existing := &warehouse.Warehouse{
                        Id:                 1,
                        WarehouseCode:     "WH-001",
                        Address:           "Test Address",
                        Telephone:         "123456789",
                        MinimumCapacity:   100,
                        MinimumTemperature: -10.5,
                        LocalityId:        "1900",
                    }
                    updated := &warehouse.Warehouse{
                        Id:                 1,
                        WarehouseCode:     "WH-001",
                        Address:           "Test Address",
                        Telephone:         "123456789",
                        MinimumCapacity:   200,
                        MinimumTemperature: -10.5,
                        LocalityId:        "1900",
                    }

                    mock.FuncFindById = func(ctx context.Context, id int) (*warehouse.Warehouse, error) {
                        return existing, nil
                    }
                    mock.FuncUpdate = func(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
                        return updated, nil
                    }
                    return mock
                },
            },
            input: input{
                id: 1,
                patch: warehouse.WarehousePatchDTO{
                    MinimumCapacity: intPtr(200),
                },
                context: context.Background(),
            },
            output: output{
                result: &warehouse.Warehouse{
                    Id:                 1,
                    WarehouseCode:     "WH-001",
                    Address:           "Test Address",
                    Telephone:         "123456789",
                    MinimumCapacity:   200,
                    MinimumTemperature: -10.5,
                    LocalityId:        "1900",
                },
                err: nil,
            },
        },
        {
            name: "error - warehouse not found",
            arrange: arrange{
                mockRepo: func() *mocks.WarehouseRepositoryMock {
                    mock := &mocks.WarehouseRepositoryMock{}
                    mock.FuncFindById = func(ctx context.Context, id int) (*warehouse.Warehouse, error) {
                        return nil, apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found")
                    }
                    return mock
                },
            },
            input: input{
                id: 99,
                patch: warehouse.WarehousePatchDTO{
                    MinimumCapacity: intPtr(200),
                },
                context: context.Background(),
            },
            output: output{
                result: nil,
                err:    apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found"),
            },
        },
        {
            name: "error - invalid minimum capacity",
            arrange: arrange{
                mockRepo: func() *mocks.WarehouseRepositoryMock {
                    mock := &mocks.WarehouseRepositoryMock{}
                    existing := &warehouse.Warehouse{
                        Id:                 1,
                        WarehouseCode:     "WH-001",
                        Address:           "Test Address",
                        Telephone:         "123456789",
                        MinimumCapacity:   100,
                        MinimumTemperature: -10.5,
                        LocalityId:        "1900",
                    }

                    mock.FuncFindById = func(ctx context.Context, id int) (*warehouse.Warehouse, error) {
                        return existing, nil
                    }
                    return mock
                },
            },
            input: input{
                id: 1,
                patch: warehouse.WarehousePatchDTO{
                    MinimumCapacity: intPtr(-10), // Invalid negative capacity
                },
                context: context.Background(),
            },
            output: output{
                result: nil,
                err:    apperrors.NewAppError(apperrors.CodeValidationError, "minimum capacity must be greater than 0"),
            },
        },
        {
            name: "error - repository update fails",
            arrange: arrange{
                mockRepo: func() *mocks.WarehouseRepositoryMock {
                    mock := &mocks.WarehouseRepositoryMock{}
                    existing := &warehouse.Warehouse{
                        Id:                 1,
                        WarehouseCode:     "WH-001",
                        Address:           "Test Address",
                        Telephone:         "123456789",
                        MinimumCapacity:   100,
                        MinimumTemperature: -10.5,
                        LocalityId:        "1900",
                    }

                    mock.FuncFindById = func(ctx context.Context, id int) (*warehouse.Warehouse, error) {
                        return existing, nil
                    }
                    mock.FuncUpdate = func(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
                        return nil, apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists")
                    }
                    return mock
                },
            },
            input: input{
                id: 1,
                patch: warehouse.WarehousePatchDTO{
                    WarehouseCode: stringPtr("EXISTING-CODE"),
                },
                context: context.Background(),
            },
            output: output{
                result: nil,
                err:    apperrors.NewAppError(apperrors.CodeConflict, "warehouse_code already exists"),
            },
        },
        {
            name: "success - update without minimum capacity validation",
            arrange: arrange{
                mockRepo: func() *mocks.WarehouseRepositoryMock {
                    mock := &mocks.WarehouseRepositoryMock{}
                    existing := &warehouse.Warehouse{
                        Id:                 1,
                        WarehouseCode:     "WH-001",
                        Address:           "Test Address",
                        Telephone:         "123456789",
                        MinimumCapacity:   100,
                        MinimumTemperature: -10.5,
                        LocalityId:        "1900",
                    }
                    updated := &warehouse.Warehouse{
                        Id:                 1,
                        WarehouseCode:     "WH-001",
                        Address:           "New Address",
                        Telephone:         "123456789",
                        MinimumCapacity:   100,
                        MinimumTemperature: -10.5,
                        LocalityId:        "1900",
                    }

                    mock.FuncFindById = func(ctx context.Context, id int) (*warehouse.Warehouse, error) {
                        return existing, nil
                    }
                    mock.FuncUpdate = func(ctx context.Context, id int, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
                        return updated, nil
                    }
                    return mock
                },
            },
            input: input{
                id: 1,
                patch: warehouse.WarehousePatchDTO{
                    Address: stringPtr("New Address"),
                },
                context: context.Background(),
            },
            output: output{
                result: &warehouse.Warehouse{
                    Id:                 1,
                    WarehouseCode:     "WH-001",
                    Address:           "New Address",
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
            result, err := srv.Update(tc.input.context, tc.input.id, tc.input.patch)

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

// Helper functions for pointer values
func intPtr(i int) *int {
    return &i
}

func stringPtr(s string) *string {
    return &s
}
