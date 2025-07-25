package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"

    carryMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/carry"
    geographyMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
    service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/carry"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
)

func TestCarryDefault_GetCarriesReport(t *testing.T) {
    type arrange struct {
        mockCarryRepo     func() *carryMocks.CarryRepositoryMock
        mockGeographyRepo func() *geographyMocks.GeographyRepositoryMock
    }
    type input struct {
        localityID *string
        context    context.Context
    }
    type output struct {
        result interface{}
        err    error
    }
    type testCase struct {
        name    string
        arrange arrange
        input   input
        output  output
    }

    // test cases
    testCases := []testCase{
        {
            name: "success - all localities (localityID is nil)",
            arrange: arrange{
                mockCarryRepo: func() *carryMocks.CarryRepositoryMock {
                    mock := &carryMocks.CarryRepositoryMock{}

                    mock.FuncGetCarriesCountByAllLocalities = func(ctx context.Context) ([]carry.CarriesReport, error) {
                        return testhelpers.CreateTestCarriesReportSlice(), nil
                    }

                    return mock
                },
                mockGeographyRepo: func() *geographyMocks.GeographyRepositoryMock {
                    return &geographyMocks.GeographyRepositoryMock{}
                },
            },
            input: input{
                localityID: nil,
                context:    context.Background(),
            },
            output: output{
                result: testhelpers.CreateTestCarriesReportSlice(),
                err:    nil,
            },
        },
        {
            name: "success - specific locality (localityID provided and exists)",
            arrange: arrange{
                mockCarryRepo: func() *carryMocks.CarryRepositoryMock {
                    mock := &carryMocks.CarryRepositoryMock{}

                    mock.FuncGetCarriesCountByLocalityID = func(ctx context.Context, localityID string) (*carry.CarriesReport, error) {
                        return testhelpers.CreateTestCarriesReport("1", "Test Locality 1", 5), nil
                    }

                    return mock
                },
                mockGeographyRepo: func() *geographyMocks.GeographyRepositoryMock {
                    mock := &geographyMocks.GeographyRepositoryMock{}

                    mock.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
                        return testhelpers.CreateTestLocality("1"), nil
                    }

                    return mock
                },
            },
            input: input{
                localityID: testhelpers.StringPtr("1"),
                context:    context.Background(),
            },
            output: output{
                result: testhelpers.CreateTestCarriesReport("1", "Test Locality 1", 5),
                err:    nil,
            },
        },
        {
            name: "error - locality not found",
            arrange: arrange{
                mockCarryRepo: func() *carryMocks.CarryRepositoryMock {
                    return &carryMocks.CarryRepositoryMock{}
                },
                mockGeographyRepo: func() *geographyMocks.GeographyRepositoryMock {
                    mock := &geographyMocks.GeographyRepositoryMock{}

                    mock.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
                        return nil, nil // locality not found
                    }

                    return mock
                },
            },
            input: input{
                localityID: testhelpers.StringPtr("999"),
                context:    context.Background(),
            },
            output: output{
                result: nil,
                err:    apperrors.NewAppError(apperrors.CodeNotFound, "locality not found"),
            },
        },
    }

    // run test cases
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // arrange
            mockCarryRepo := tc.arrange.mockCarryRepo()
            mockGeographyRepo := tc.arrange.mockGeographyRepo()
            srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)

            // act
            result, err := srv.GetCarriesReport(tc.input.context, tc.input.localityID)

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