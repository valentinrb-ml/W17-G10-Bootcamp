package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/carry"
	carryMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/carry"
	geographyMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
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

func TestCarryDefault_GetCarriesReport_AllLocalities_WithLogger(t *testing.T) {
	// arrange - all localities with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	mockCarryRepo.FuncGetCarriesCountByAllLocalities = func(ctx context.Context) ([]carry.CarriesReport, error) {
		return testhelpers.CreateTestCarriesReportSlice(), nil
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	// act
	result, err := srv.GetCarriesReport(context.Background(), nil)

	// assert
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestCarryDefault_GetCarriesReport_AllLocalities_Error_WithLogger(t *testing.T) {
	// arrange - all localities error with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	mockCarryRepo.FuncGetCarriesCountByAllLocalities = func(ctx context.Context) ([]carry.CarriesReport, error) {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "repository error")
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	// act
	result, err := srv.GetCarriesReport(context.Background(), nil)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
}

func TestCarryDefault_GetCarriesReport_SpecificLocality_Success_WithLogger(t *testing.T) {
	// arrange - specific locality success with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	localityID := "LOC001"
	locality := testhelpers.CreateTestLocality("LOC001")

	mockGeographyRepo.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
		return locality, nil
	}

	mockCarryRepo.FuncGetCarriesCountByLocalityID = func(ctx context.Context, localityID string) (*carry.CarriesReport, error) {
		return testhelpers.CreateTestCarriesReport("LOC001", "Test Locality", 5), nil
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	// act
	result, err := srv.GetCarriesReport(context.Background(), &localityID)

	// assert
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestCarryDefault_GetCarriesReport_LocalityNotFound_WithLogger(t *testing.T) {
	// arrange - locality not found with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	localityID := "LOC999"

	mockGeographyRepo.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
		return nil, nil // Not found
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	// act
	result, err := srv.GetCarriesReport(context.Background(), &localityID)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "locality not found")
}

func TestCarryDefault_GetCarriesReport_GeographyError_WithLogger(t *testing.T) {
	// arrange - geography repository error with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	localityID := "LOC001"

	mockGeographyRepo.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "geography repository error")
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	// act
	result, err := srv.GetCarriesReport(context.Background(), &localityID)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
}

func TestCarryDefault_GetCarriesReport_CarryRepoError_WithLogger(t *testing.T) {
	// arrange - carry repository error with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	localityID := "LOC001"
	locality := testhelpers.CreateTestLocality("LOC001")

	mockGeographyRepo.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
		return locality, nil
	}

	mockCarryRepo.FuncGetCarriesCountByLocalityID = func(ctx context.Context, localityID string) (*carry.CarriesReport, error) {
		return nil, apperrors.NewAppError(apperrors.CodeInternal, "carry repository error")
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	// act
	result, err := srv.GetCarriesReport(context.Background(), &localityID)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
}
