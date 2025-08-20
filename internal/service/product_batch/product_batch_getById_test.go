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

func TestProductBatchesService_GetReportProductById(t *testing.T) {
	type arrange struct {
		repoMock func() *mocks.ProductBatchServiceMock
	}
	type output struct {
		expected      *models.ReportProduct
		expectedError bool
		err           error
	}
	type input struct {
		sectionNumber int
	}
	type testCase struct {
		name string
		arrange
		output
		input
	}

	report := testhelpers.DummyReportProduct()

	testCases := []testCase{
		{
			name: "returns report product on success",
			arrange: arrange{
				repoMock: func() *mocks.ProductBatchServiceMock {
					return &mocks.ProductBatchServiceMock{
						FuncGetReportById: func(ctx context.Context, sectionNumber int) (*models.ReportProduct, error) {
							dummy := testhelpers.DummyReportProduct()
							return &dummy, nil
						},
					}
				},
			},
			input: input{sectionNumber: report.SectionId},
			output: output{
				expected:      &report,
				expectedError: false,
				err:           nil,
			},
		},
		{
			name: "returns error when repo fails",
			arrange: arrange{
				repoMock: func() *mocks.ProductBatchServiceMock {
					return &mocks.ProductBatchServiceMock{
						FuncGetReportById: func(ctx context.Context, sectionNumber int) (*models.ReportProduct, error) {
							return nil, context.DeadlineExceeded
						},
					}
				},
			},
			input: input{sectionNumber: 9999},
			output: output{
				expected:      nil,
				expectedError: true,
				err:           context.DeadlineExceeded,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewProductBatchesService(tc.arrange.repoMock())
			svc.SetLogger(testhelpers.NewTestLogger())

			result, err := svc.GetReportProductById(context.Background(), tc.input.sectionNumber)

			if tc.output.expectedError {
				require.Error(t, err)
				require.Equal(t, tc.output.err, err)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.output.expected, result)
		})
	}
}
