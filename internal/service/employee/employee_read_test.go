package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
	employeeMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/employee"
	warehouseMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/employee"
)

func TestEmployeeService_Read(t *testing.T) {
	testCases := []struct {
		name        string
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		findAll     bool
		inputID     int // para find_by_id casos
		wantErr     bool
		wantErrCode string
		wantLen     int // Para find_all caso feliz
		wantID      int // Para find_by_id_existente
	}{
		{
			name: "find_all",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindAll: func(ctx context.Context) ([]*models.Employee, error) {
						return []*models.Employee{
							{ID: 1, CardNumberID: "E001", FirstName: "Lucas", LastName: "Martinez", WarehouseID: 1},
							{ID: 2, CardNumberID: "E002", FirstName: "Sonia", LastName: "Lopez", WarehouseID: 2},
						}, nil
					},
				}
			},
			findAll: true,
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "find_by_id_non_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return nil, apperrors.NewAppError(apperrors.CodeNotFound, "employee not found")
					},
				}
			},
			findAll:     false,
			inputID:     99,
			wantErr:     true,
			wantErrCode: apperrors.CodeNotFound,
		},
		{
			name: "find_by_id_existent",
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) {
						return &models.Employee{
							ID: id, CardNumberID: "E015", FirstName: "Maria", LastName: "Perez", WarehouseID: 3,
						}, nil
					},
				}
			},
			findAll: false,
			inputID: 15,
			wantErr: false,
			wantID:  15,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emRepo := tc.repoMock()
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)

			if tc.findAll {
				result, err := svc.FindAll(context.Background())
				require.False(t, tc.wantErr)
				require.NoError(t, err)
				require.Len(t, result, tc.wantLen)
			} else {
				res, err := svc.FindByID(context.Background(), tc.inputID)
				if tc.wantErr {
					require.Error(t, err)
					appErr, ok := err.(*apperrors.AppError)
					require.True(t, ok)
					require.Equal(t, tc.wantErrCode, appErr.Code)
				} else {
					require.NoError(t, err)
					require.NotNil(t, res)
					require.Equal(t, tc.wantID, res.ID)
				}
			}
		})
	}
}

func TestEmployeeService_Read_extraCases(t *testing.T) {
	testCases := []struct {
		name        string
		id          int
		findAll     bool
		repoMock    func() *employeeMocks.EmployeeRepositoryMock
		wantErrCode string
		checkWrap   string
	}{
		{
			name: "invalid id for FindByID",
			id:   0,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) { return nil, nil },
				}
			},
			wantErrCode: apperrors.CodeValidationError,
		},
		{
			name: "repo.FindByID returns error",
			id:   5,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) { return nil, context.Canceled },
				}
			},
			checkWrap: "failed fetching employee by id",
		},
		{
			name: "repo.FindByID returns nil, nil",
			id:   15,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindByID: func(ctx context.Context, id int) (*models.Employee, error) { return nil, nil },
				}
			},
			wantErrCode: apperrors.CodeNotFound,
		},
		{
			name:    "repo.FindAll returns error",
			findAll: true,
			repoMock: func() *employeeMocks.EmployeeRepositoryMock {
				return &employeeMocks.EmployeeRepositoryMock{
					MockFindAll: func(ctx context.Context) ([]*models.Employee, error) { return nil, context.DeadlineExceeded },
				}
			},
			checkWrap: "failed fetching all employees",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emRepo := tc.repoMock()
			whRepo := &warehouseMocks.WarehouseRepositoryMock{}
			svc := service.NewEmployeeDefault(emRepo, whRepo)
			if tc.findAll {
				res, err := svc.FindAll(context.Background())
				if tc.wantErrCode != "" {
					require.Error(t, err)
					appErr, ok := err.(*apperrors.AppError)
					require.True(t, ok)
					require.Equal(t, tc.wantErrCode, appErr.Code)
					require.Nil(t, res)
				} else if tc.checkWrap != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.checkWrap)
					require.Nil(t, res)
				} else {
					require.NoError(t, err)
					require.NotNil(t, res)
				}
			} else {
				res, err := svc.FindByID(context.Background(), tc.id)
				if tc.wantErrCode != "" {
					require.Error(t, err)
					appErr, ok := err.(*apperrors.AppError)
					require.True(t, ok)
					require.Equal(t, tc.wantErrCode, appErr.Code)
					require.Nil(t, res)
				} else if tc.checkWrap != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.checkWrap)
					require.Nil(t, res)
				} else {
					require.NoError(t, err)
					require.NotNil(t, res)
				}
			}
		})
	}
}
