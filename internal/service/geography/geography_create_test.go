package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestGeographyService_Create(t *testing.T) {
	tests := []struct {
		name     string
		mockRepo func() *mocks.GeographyRepositoryMock
		wantErr  bool
		wantMsg  string
		wantResp *models.ResponseGeography
	}{
		{
			name: "success",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncBeginTx = func(ctx context.Context) (*sql.Tx, error) { return &sql.Tx{}, nil }
				mock.FuncCommitTx = func(tx *sql.Tx) error { return nil }
				mock.FuncRollbackTx = func(tx *sql.Tx) error { return nil }
				mock.FuncFindCountryByName = func(ctx context.Context, name string) (*models.Country, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateCountry = func(ctx context.Context, exec repository.Executor, c models.Country) (*models.Country, error) {
					return &models.Country{Id: 1, Name: c.Name}, nil
				}
				mock.FuncFindProvinceByName = func(ctx context.Context, name string, countryId int) (*models.Province, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateProvince = func(ctx context.Context, exec repository.Executor, p models.Province) (*models.Province, error) {
					return &models.Province{Id: 2, Name: p.Name, CountryId: 1}, nil
				}
				mock.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateLocality = func(ctx context.Context, exec repository.Executor, l models.Locality) (*models.Locality, error) {
					return &models.Locality{Id: "5194", Name: l.Name, ProvinceId: l.ProvinceId}, nil
				}
				return mock
			},
			wantErr:  false,
			wantResp: &models.ResponseGeography{LocalityId: "5194", LocalityName: "Villa General Belgrano", ProvinceName: "Cordoba", CountryName: "Argentina"},
		},
		{
			name: "error on begin tx",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncBeginTx = func(ctx context.Context) (*sql.Tx, error) { return nil, errors.New("cannot begin tx") }
				return mock
			},
			wantErr: true,
			wantMsg: "failed to start transaction",
		},
		{
			name: "error on create country",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncBeginTx = func(ctx context.Context) (*sql.Tx, error) { return &sql.Tx{}, nil }
				mock.FuncRollbackTx = func(tx *sql.Tx) error { return nil }
				mock.FuncFindCountryByName = func(ctx context.Context, name string) (*models.Country, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateCountry = func(ctx context.Context, exec repository.Executor, c models.Country) (*models.Country, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "error creating country")
				}
				return mock
			},
			wantErr: true,
			wantMsg: "error creating country",
		},
		{
			name: "error on create province",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncBeginTx = func(ctx context.Context) (*sql.Tx, error) { return &sql.Tx{}, nil }
				mock.FuncRollbackTx = func(tx *sql.Tx) error { return nil }
				mock.FuncFindCountryByName = func(ctx context.Context, name string) (*models.Country, error) {
					return &models.Country{Id: 1, Name: name}, nil
				}
				mock.FuncFindProvinceByName = func(ctx context.Context, name string, countryId int) (*models.Province, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateProvince = func(ctx context.Context, exec repository.Executor, p models.Province) (*models.Province, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "error creating province")
				}
				return mock
			},
			wantErr: true,
			wantMsg: "error creating province",
		},
		{
			name: "locality already exists",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncBeginTx = func(ctx context.Context) (*sql.Tx, error) { return &sql.Tx{}, nil }
				mock.FuncRollbackTx = func(tx *sql.Tx) error { return nil }
				mock.FuncFindCountryByName = func(ctx context.Context, name string) (*models.Country, error) {
					return &models.Country{Id: 1, Name: name}, nil
				}
				mock.FuncFindProvinceByName = func(ctx context.Context, name string, countryId int) (*models.Province, error) {
					return &models.Province{Id: 2, Name: name, CountryId: countryId}, nil
				}
				mock.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
					return &models.Locality{Id: id, Name: "Already", ProvinceId: 2}, nil
				}
				return mock
			},
			wantErr: true,
			wantMsg: "locality already exists",
		},
		{
			name: "error on create locality",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncBeginTx = func(ctx context.Context) (*sql.Tx, error) { return &sql.Tx{}, nil }
				mock.FuncRollbackTx = func(tx *sql.Tx) error { return nil }
				mock.FuncFindCountryByName = func(ctx context.Context, name string) (*models.Country, error) {
					return &models.Country{Id: 1, Name: name}, nil
				}
				mock.FuncFindProvinceByName = func(ctx context.Context, name string, countryId int) (*models.Province, error) {
					return &models.Province{Id: 2, Name: name, CountryId: countryId}, nil
				}
				mock.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateLocality = func(ctx context.Context, exec repository.Executor, l models.Locality) (*models.Locality, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "error creating locality")
				}
				return mock
			},
			wantErr: true,
			wantMsg: "error creating locality",
		},
		{
			name: "commit error",
			mockRepo: func() *mocks.GeographyRepositoryMock {
				mock := &mocks.GeographyRepositoryMock{}
				mock.FuncBeginTx = func(ctx context.Context) (*sql.Tx, error) { return &sql.Tx{}, nil }
				mock.FuncCommitTx = func(tx *sql.Tx) error { return errors.New("commit error") }
				mock.FuncRollbackTx = func(tx *sql.Tx) error { return nil }
				mock.FuncFindCountryByName = func(ctx context.Context, name string) (*models.Country, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateCountry = func(ctx context.Context, exec repository.Executor, c models.Country) (*models.Country, error) {
					return &models.Country{Id: 1, Name: c.Name}, nil
				}
				mock.FuncFindProvinceByName = func(ctx context.Context, name string, countryId int) (*models.Province, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateProvince = func(ctx context.Context, exec repository.Executor, p models.Province) (*models.Province, error) {
					return &models.Province{Id: 2, Name: p.Name, CountryId: 1}, nil
				}
				mock.FuncFindLocalityById = func(ctx context.Context, id string) (*models.Locality, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "not found")
				}
				mock.FuncCreateLocality = func(ctx context.Context, exec repository.Executor, l models.Locality) (*models.Locality, error) {
					return &models.Locality{Id: "5194", Name: l.Name, ProvinceId: l.ProvinceId}, nil
				}
				return mock
			},
			wantErr: true,
			wantMsg: "failed to commit transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			svc := service.NewGeographyService(repo)
			resp, err := svc.Create(context.Background(), testhelpers.DummyRequestGeography())
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantMsg != "" {
					require.Contains(t, err.Error(), tt.wantMsg)
				}
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
