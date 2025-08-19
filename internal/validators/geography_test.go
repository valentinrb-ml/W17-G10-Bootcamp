package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

func TestValidateGeographyPost(t *testing.T) {
	validPostal := "1234"
	validCountry := "Argentina"
	validProvince := "Buenos Aires"
	validLocality := "CABA"

	cases := []struct {
		name      string
		input     models.RequestGeography
		wantError bool
		wantMsg   string
	}{
		{
			name: "all valid",
			input: models.RequestGeography{
				Id:           strPtr(validPostal),
				CountryName:  strPtr(validCountry),
				ProvinceName: strPtr(validProvince),
				LocalityName: strPtr(validLocality),
			},
			wantError: false,
		},
		{
			name:      "all nil",
			input:     models.RequestGeography{},
			wantError: true,
			wantMsg:   "Id (Postal Code) is required and cannot be empty.",
		},
		{
			name: "postal empty",
			input: models.RequestGeography{
				Id:           strPtr(""),
				CountryName:  strPtr(validCountry),
				ProvinceName: strPtr(validProvince),
				LocalityName: strPtr(validLocality),
			},
			wantError: true,
			wantMsg:   "Id (Postal Code) is required and cannot be empty.",
		},
		{
			name: "country empty",
			input: models.RequestGeography{
				Id:           strPtr(validPostal),
				CountryName:  strPtr(""),
				ProvinceName: strPtr(validProvince),
				LocalityName: strPtr(validLocality),
			},
			wantError: true,
			wantMsg:   "Country Name is required and cannot be empty.",
		},
		{
			name: "province empty",
			input: models.RequestGeography{
				Id:           strPtr(validPostal),
				CountryName:  strPtr(validCountry),
				ProvinceName: strPtr(""),
				LocalityName: strPtr(validLocality),
			},
			wantError: true,
			wantMsg:   "Pronvice Name is required and cannot be empty.",
		},
		{
			name: "locality empty",
			input: models.RequestGeography{
				Id:           strPtr(validPostal),
				CountryName:  strPtr(validCountry),
				ProvinceName: strPtr(validProvince),
				LocalityName: strPtr(""),
			},
			wantError: true,
			wantMsg:   "Locality Name is required and cannot be empty.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateGeographyPost(tc.input)
			if tc.wantError {
				assert.Error(t, err)
				if appErr, ok := err.(*apperrors.AppError); ok && tc.wantMsg != "" {
					assert.Equal(t, tc.wantMsg, appErr.Message)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
