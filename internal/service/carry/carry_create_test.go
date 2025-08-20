package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/carry"
	carryMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/carry"
	geographyMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestCarryDefault_Create(t *testing.T) {
	type arrange struct {
		mockCarryRepo     func() *carryMocks.CarryRepositoryMock
		mockGeographyRepo func() *geographyMocks.GeographyRepositoryMock
	}
	type input struct {
		carry   carry.Carry
		context context.Context
	}
	type output struct {
		result *carry.Carry
		err    error
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	// test cases (only happy case for handrail method)
	testCases := []testCase{
		{
			name: "success - carry created",
			arrange: arrange{
				mockCarryRepo: func() *carryMocks.CarryRepositoryMock {
					mock := &carryMocks.CarryRepositoryMock{}

					mock.FuncCreate = func(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
						return testhelpers.CreateExpectedCarry(1), nil
					}

					return mock
				},
				mockGeographyRepo: func() *geographyMocks.GeographyRepositoryMock {
					return &geographyMocks.GeographyRepositoryMock{}
				},
			},
			input: input{
				carry:   testhelpers.CreateTestCarryForCreate(),
				context: context.Background(),
			},
			output: output{
				result: testhelpers.CreateExpectedCarry(1),
				err:    nil,
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
			result, err := srv.Create(tc.input.context, tc.input.carry)

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

func TestCarryDefault_Create_Success_WithLogger(t *testing.T) {
	// arrange - success case with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	expectedCarry := testhelpers.CreateExpectedCarry(1)

	mockCarryRepo.FuncCreate = func(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
		return expectedCarry, nil
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	testCarry := testhelpers.CreateTestCarry(0)

	// act
	result, err := srv.Create(context.Background(), *testCarry)

	// assert
	require.NoError(t, err)
	require.Equal(t, expectedCarry, result)
}

func TestCarryDefault_Create_Error_WithLogger(t *testing.T) {
	// arrange - error case with logger
	mockCarryRepo := &carryMocks.CarryRepositoryMock{}
	mockGeographyRepo := &geographyMocks.GeographyRepositoryMock{}

	mockCarryRepo.FuncCreate = func(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
		return nil, errors.New("repository error")
	}

	srv := service.NewCarryService(mockCarryRepo, mockGeographyRepo)
	srv.SetLogger(&SimpleTestLogger{})

	testCarry := testhelpers.CreateTestCarry(0)

	// act
	result, err := srv.Create(context.Background(), *testCarry)

	// assert
	require.Error(t, err)
	require.Nil(t, result)
}
