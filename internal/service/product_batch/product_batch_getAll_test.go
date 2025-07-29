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

func TestProductBatchesService_GetReportProduct(t *testing.T) {
	type arrange struct {
		repoMock func() *mocks.ProductBatchServiceMock
	}
	type output struct {
		expected      []models.ReportProduct
		expectedError bool
		err           error
	}
	type testCase struct {
		name string
		arrange
		output
	}

	expectedReport := testhelpers.DummyReportProductsList()

	testCases := []testCase{
		{
			name: "returns product batch report successfully",
			arrange: arrange{
				repoMock: func() *mocks.ProductBatchServiceMock {
					return &mocks.ProductBatchServiceMock{
						FuncGetReport: func(ctx context.Context) ([]models.ReportProduct, error) {
							return expectedReport, nil
						},
					}
				},
			},
			output: output{
				expected:      expectedReport,
				expectedError: false,
				err:           nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewProductBatchesService(tc.arrange.repoMock())

			result, err := svc.GetReportProduct(context.Background())

			if tc.output.expectedError {
				require.Error(t, err)
				require.EqualError(t, err, tc.output.err.Error())
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.output.expected, result)
		})
	}
}
