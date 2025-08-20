package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product_batch"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product_batch"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestProductBatchesService_CreateProductBatches(t *testing.T) {
	type arrange struct {
		repoMock func() *mocks.ProductBatchServiceMock
	}
	type output struct {
		expected      *models.ProductBatches
		expectedError bool
		err           error
	}
	type input struct {
		batch models.ProductBatches
	}
	type testCase struct {
		name string
		arrange
		output
		input
	}

	prodBatch := testhelpers.DummyProductBatch(1)

	testCases := []testCase{
		{
			name: "returns new product batch on successful creation",
			arrange: arrange{
				repoMock: func() *mocks.ProductBatchServiceMock {
					return &mocks.ProductBatchServiceMock{
						FuncCreate: func(ctx context.Context, proBa models.ProductBatches) (*models.ProductBatches, error) {
							dummy := testhelpers.DummyProductBatch(1)
							return &dummy, nil
						},
					}
				},
			},
			input: input{batch: prodBatch},
			output: output{
				expected:      &prodBatch,
				expectedError: false,
				err:           nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewProductBatchesService(tc.arrange.repoMock())
			svc.SetLogger(testhelpers.NewTestLogger())

			result, err := svc.CreateProductBatches(context.Background(), tc.input.batch)

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
