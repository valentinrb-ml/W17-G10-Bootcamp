package mappers

import (
	"database/sql"
	"testing"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

func TestToDomain(t *testing.T) {
	t.Parallel()

	iPtr := func(i int) *int { return &i }

	tests := []struct {
		name string
		in   models.ProductRequest
		want models.Product
	}{
		{
			name: "maps all fields with nil seller",
			in: models.ProductRequest{
				ProductData: models.ProductData{
					ProductCode:                    "P-001",
					Description:                    "Prod 1",
					Width:                          1.1,
					Height:                         2.2,
					Length:                         3.3,
					NetWeight:                      0.5,
					ExpirationRate:                 0.1,
					RecommendedFreezingTemperature: -18.0,
					FreezingRate:                   0.2,
					ProductTypeID:                  7,
					SellerID:                       nil,
				},
			},
			want: models.Product{
				Code:        "P-001",
				Description: "Prod 1",
				Dimensions:  models.Dimensions{Width: 1.1, Height: 2.2, Length: 3.3},
				NetWeight:   0.5,
				Expiration: models.Expiration{
					Rate:                    0.1,
					RecommendedFreezingTemp: -18.0,
					FreezingRate:            0.2,
				},
				ProductType: 7,
				SellerID:    nil,
			},
		},
		{
			name: "maps all fields with seller id",
			in: models.ProductRequest{
				ProductData: models.ProductData{
					ProductCode:                    "P-002",
					Description:                    "Prod 2",
					Width:                          10,
					Height:                         20,
					Length:                         30,
					NetWeight:                      5,
					ExpirationRate:                 1.5,
					RecommendedFreezingTemperature: -10.5,
					FreezingRate:                   2.5,
					ProductTypeID:                  9,
					SellerID:                       iPtr(42),
				},
			},
			want: models.Product{
				Code:        "P-002",
				Description: "Prod 2",
				Dimensions:  models.Dimensions{Width: 10, Height: 20, Length: 30},
				NetWeight:   5,
				Expiration: models.Expiration{
					Rate:                    1.5,
					RecommendedFreezingTemp: -10.5,
					FreezingRate:            2.5,
				},
				ProductType: 9,
				SellerID:    iPtr(42),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ToDomain(tt.in)
			assertEqualProduct(t, got, tt.want)
		})
	}
}

func TestFromDomain(t *testing.T) {
	t.Parallel()

	iPtr := func(i int) *int { return &i }

	tests := []struct {
		name string
		in   models.Product
		want models.ProductResponse
	}{
		{
			name: "maps product to response with nil seller",
			in: models.Product{
				ID:          5,
				Code:        "PX-1",
				Description: "desc",
				Dimensions:  models.Dimensions{Width: 1, Height: 2, Length: 3},
				NetWeight:   0.9,
				Expiration: models.Expiration{
					Rate:                    0.01,
					RecommendedFreezingTemp: -5,
					FreezingRate:            0.02,
				},
				ProductType: 3,
				SellerID:    nil,
			},
			want: models.ProductResponse{
				ID: 5,
				ProductData: models.ProductData{
					ProductCode:                    "PX-1",
					Description:                    "desc",
					Width:                          1,
					Height:                         2,
					Length:                         3,
					NetWeight:                      0.9,
					ExpirationRate:                 0.01,
					RecommendedFreezingTemperature: -5,
					FreezingRate:                   0.02,
					ProductTypeID:                  3,
					SellerID:                       nil,
				},
			},
		},
		{
			name: "maps product to response with seller",
			in: models.Product{
				ID:          7,
				Code:        "PX-2",
				Description: "desc 2",
				Dimensions:  models.Dimensions{Width: 10, Height: 20, Length: 30},
				NetWeight:   1.2,
				Expiration: models.Expiration{
					Rate:                    1,
					RecommendedFreezingTemp: -20,
					FreezingRate:            2,
				},
				ProductType: 4,
				SellerID:    iPtr(100),
			},
			want: models.ProductResponse{
				ID: 7,
				ProductData: models.ProductData{
					ProductCode:                    "PX-2",
					Description:                    "desc 2",
					Width:                          10,
					Height:                         20,
					Length:                         30,
					NetWeight:                      1.2,
					ExpirationRate:                 1,
					RecommendedFreezingTemperature: -20,
					FreezingRate:                   2,
					ProductTypeID:                  4,
					SellerID:                       iPtr(100),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := FromDomain(tt.in)
			assertEqualProductResponse(t, got, tt.want)
		})
	}
}

func TestFromDomainList(t *testing.T) {
	t.Parallel()

	iPtr := func(i int) *int { return &i }

	tests := []struct {
		name string
		in   []models.Product
	}{
		{
			name: "nil list => empty slice",
			in:   nil,
		},
		{
			name: "empty list => empty slice",
			in:   []models.Product{},
		},
		{
			name: "maps multiple items",
			in: []models.Product{
				{
					ID:          1,
					Code:        "C1",
					Description: "D1",
					Dimensions:  models.Dimensions{Width: 1, Height: 2, Length: 3},
					NetWeight:   0.1,
					Expiration:  models.Expiration{Rate: 0.2, RecommendedFreezingTemp: -1, FreezingRate: 0.3},
					ProductType: 10,
					SellerID:    nil,
				},
				{
					ID:          2,
					Code:        "C2",
					Description: "D2",
					Dimensions:  models.Dimensions{Width: 4, Height: 5, Length: 6},
					NetWeight:   0.4,
					Expiration:  models.Expiration{Rate: 1.2, RecommendedFreezingTemp: -2, FreezingRate: 1.3},
					ProductType: 20,
					SellerID:    iPtr(8),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FromDomainList(tt.in)
			if len(tt.in) == 0 {
				if len(got) != 0 {
					t.Fatalf("expected empty slice, got len=%d", len(got))
				}
				return
			}
			if len(got) != len(tt.in) {
				t.Fatalf("len mismatch: got %d, want %d", len(got), len(tt.in))
			}
			for i := range tt.in {
				want := FromDomain(tt.in[i])
				assertEqualProductResponse(t, got[i], want)
			}
		})
	}
}

func TestResponseToDomain(t *testing.T) {
	t.Parallel()

	iPtr := func(i int) *int { return &i }

	tests := []struct {
		name string
		in   models.ProductResponse
		want models.Product
	}{
		{
			name: "maps all fields with nil seller",
			in: models.ProductResponse{
				ID: 11,
				ProductData: models.ProductData{
					ProductCode:                    "R-1",
					Description:                    "Resp desc",
					Width:                          7,
					Height:                         8,
					Length:                         9,
					NetWeight:                      0.77,
					ExpirationRate:                 0.5,
					RecommendedFreezingTemperature: -7.7,
					FreezingRate:                   0.6,
					ProductTypeID:                  99,
					SellerID:                       nil,
				},
			},
			want: models.Product{
				ID:          11,
				Code:        "R-1",
				Description: "Resp desc",
				Dimensions:  models.Dimensions{Width: 7, Height: 8, Length: 9},
				NetWeight:   0.77,
				Expiration: models.Expiration{
					Rate:                    0.5,
					RecommendedFreezingTemp: -7.7,
					FreezingRate:            0.6,
				},
				ProductType: 99,
				SellerID:    nil,
			},
		},
		{
			name: "maps all fields with seller",
			in: models.ProductResponse{
				ID: 12,
				ProductData: models.ProductData{
					ProductCode:                    "R-2",
					Description:                    "Desc 2",
					Width:                          12,
					Height:                         13,
					Length:                         14,
					NetWeight:                      1.5,
					ExpirationRate:                 2.5,
					RecommendedFreezingTemperature: -12.5,
					FreezingRate:                   3.5,
					ProductTypeID:                  77,
					SellerID:                       iPtr(55),
				},
			},
			want: models.Product{
				ID:          12,
				Code:        "R-2",
				Description: "Desc 2",
				Dimensions:  models.Dimensions{Width: 12, Height: 13, Length: 14},
				NetWeight:   1.5,
				Expiration: models.Expiration{
					Rate:                    2.5,
					RecommendedFreezingTemp: -12.5,
					FreezingRate:            3.5,
				},
				ProductType: 77,
				SellerID:    iPtr(55),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ResponseToDomain(tt.in)
			assertEqualProduct(t, got, tt.want)
		})
	}
}

func TestDbToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   models.ProductDb
		want models.Product
	}{
		{
			name: "maps with null seller id",
			in: models.ProductDb{
				ID:          1,
				Code:        "DB1",
				Description: "From DB",
				Width:       1,
				Height:      2,
				Length:      3,
				NetWeight:   0.5,
				ExpRate:     0.7,
				RecFreeze:   -7,
				FreezeRate:  0.9,
				TypeID:      10,
				SellerID:    sql.NullInt64{Valid: false},
			},
			want: models.Product{
				ID:          1,
				Code:        "DB1",
				Description: "From DB",
				Dimensions:  models.Dimensions{Width: 1, Height: 2, Length: 3},
				NetWeight:   0.5,
				Expiration: models.Expiration{
					Rate:                    0.7,
					RecommendedFreezingTemp: -7,
					FreezingRate:            0.9,
				},
				ProductType: 10,
				SellerID:    nil,
			},
		},
		{
			name: "maps with seller id present",
			in: models.ProductDb{
				ID:          2,
				Code:        "DB2",
				Description: "From DB 2",
				Width:       10,
				Height:      20,
				Length:      30,
				NetWeight:   5,
				ExpRate:     1.1,
				RecFreeze:   -10,
				FreezeRate:  2.2,
				TypeID:      11,
				SellerID:    sql.NullInt64{Int64: 77, Valid: true},
			},
			want: models.Product{
				ID:          2,
				Code:        "DB2",
				Description: "From DB 2",
				Dimensions:  models.Dimensions{Width: 10, Height: 20, Length: 30},
				NetWeight:   5,
				Expiration: models.Expiration{
					Rate:                    1.1,
					RecommendedFreezingTemp: -10,
					FreezingRate:            2.2,
				},
				ProductType: 11,
				SellerID:    intPtr(77),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := DbToDomain(tt.in)
			assertEqualProduct(t, got, tt.want)
		})
	}
}

func TestFromDomainToDb(t *testing.T) {
	t.Parallel()

	iPtr := func(i int) *int { return &i }

	tests := []struct {
		name string
		in   models.Product
		want models.ProductDb
	}{
		{
			name: "maps with nil seller",
			in: models.Product{
				ID:          10,
				Code:        "D2DB",
				Description: "Dom 2 DB",
				Dimensions:  models.Dimensions{Width: 2, Height: 3, Length: 4},
				NetWeight:   0.2,
				Expiration: models.Expiration{
					Rate:                    0.01,
					RecommendedFreezingTemp: -2.2,
					FreezingRate:            0.02,
				},
				ProductType: 5,
				SellerID:    nil,
			},
			want: models.ProductDb{
				ID:          10,
				Code:        "D2DB",
				Description: "Dom 2 DB",
				Width:       2,
				Height:      3,
				Length:      4,
				NetWeight:   0.2,
				ExpRate:     0.01,
				RecFreeze:   -2.2,
				FreezeRate:  0.02,
				TypeID:      5,
				SellerID:    sql.NullInt64{Valid: false},
			},
		},
		{
			name: "maps with seller present",
			in: models.Product{
				ID:          11,
				Code:        "D2DB-2",
				Description: "Dom 2 DB 2",
				Dimensions:  models.Dimensions{Width: 12, Height: 13, Length: 14},
				NetWeight:   1.4,
				Expiration: models.Expiration{
					Rate:                    2.3,
					RecommendedFreezingTemp: -12,
					FreezingRate:            3.4,
				},
				ProductType: 6,
				SellerID:    iPtr(99),
			},
			want: models.ProductDb{
				ID:          11,
				Code:        "D2DB-2",
				Description: "Dom 2 DB 2",
				Width:       12,
				Height:      13,
				Length:      14,
				NetWeight:   1.4,
				ExpRate:     2.3,
				RecFreeze:   -12,
				FreezeRate:  3.4,
				TypeID:      6,
				SellerID:    sql.NullInt64{Int64: 99, Valid: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := FromDomainToDb(tt.in)

			// Compare all non-NullInt64 fields directly
			if got.ID != tt.want.ID ||
				got.Code != tt.want.Code ||
				got.Description != tt.want.Description ||
				got.Width != tt.want.Width ||
				got.Height != tt.want.Height ||
				got.Length != tt.want.Length ||
				got.NetWeight != tt.want.NetWeight ||
				got.ExpRate != tt.want.ExpRate ||
				got.RecFreeze != tt.want.RecFreeze ||
				got.FreezeRate != tt.want.FreezeRate ||
				got.TypeID != tt.want.TypeID {
				t.Fatalf("mismatch.\n got: %+v\nwant: %+v", got, tt.want)
			}

			// Compare NullInt64 manually
			if got.SellerID.Valid != tt.want.SellerID.Valid {
				t.Fatalf("SellerID.Valid mismatch: got %v, want %v", got.SellerID.Valid, tt.want.SellerID.Valid)
			}
			if got.SellerID.Valid && got.SellerID.Int64 != tt.want.SellerID.Int64 {
				t.Fatalf("SellerID.Int64 mismatch: got %d, want %d", got.SellerID.Int64, tt.want.SellerID.Int64)
			}
		})
	}
}

// Helpers

func intPtr(v int) *int { return &v }

func assertEqualProduct(t *testing.T, got, want models.Product) {
	t.Helper()

	if got.ID != want.ID ||
		got.Code != want.Code ||
		got.Description != want.Description ||
		got.Dimensions.Width != want.Dimensions.Width ||
		got.Dimensions.Height != want.Dimensions.Height ||
		got.Dimensions.Length != want.Dimensions.Length ||
		got.NetWeight != want.NetWeight ||
		got.Expiration.Rate != want.Expiration.Rate ||
		got.Expiration.RecommendedFreezingTemp != want.Expiration.RecommendedFreezingTemp ||
		got.Expiration.FreezingRate != want.Expiration.FreezingRate ||
		got.ProductType != want.ProductType {
		t.Fatalf("product mismatch.\n got: %+v\nwant: %+v", got, want)
	}

	switch {
	case got.SellerID == nil && want.SellerID != nil:
		t.Fatalf("SellerID mismatch: got nil, want %d", *want.SellerID)
	case got.SellerID != nil && want.SellerID == nil:
		t.Fatalf("SellerID mismatch: got %d, want nil", *got.SellerID)
	case got.SellerID != nil && want.SellerID != nil && *got.SellerID != *want.SellerID:
		t.Fatalf("SellerID value mismatch: got %d, want %d", *got.SellerID, *want.SellerID)
	}
}

func assertEqualProductResponse(t *testing.T, got, want models.ProductResponse) {
	t.Helper()

	if got.ID != want.ID ||
		got.ProductCode != want.ProductCode ||
		got.Description != want.Description ||
		got.Width != want.Width ||
		got.Height != want.Height ||
		got.Length != want.Length ||
		got.NetWeight != want.NetWeight ||
		got.ExpirationRate != want.ExpirationRate ||
		got.RecommendedFreezingTemperature != want.RecommendedFreezingTemperature ||
		got.FreezingRate != want.FreezingRate ||
		got.ProductTypeID != want.ProductTypeID {
		t.Fatalf("product response mismatch.\n got: %+v\nwant: %+v", got, want)
	}

	switch {
	case got.SellerID == nil && want.SellerID != nil:
		t.Fatalf("SellerID mismatch: got nil, want %d", *want.SellerID)
	case got.SellerID != nil && want.SellerID == nil:
		t.Fatalf("SellerID mismatch: got %d, want nil", *got.SellerID)
	case got.SellerID != nil && want.SellerID != nil && *got.SellerID != *want.SellerID:
		t.Fatalf("SellerID value mismatch: got %d, want %d", *got.SellerID, *want.SellerID)
	}
}
