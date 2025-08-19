package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func TestValidateSellerPost(t *testing.T) {
	validCid := 1
	validCompany := "Company"
	validAddress := "Address"
	validTel := "123456"
	validLoc := "10"

	cases := []struct {
		name      string
		input     models.RequestSeller
		wantError bool
		wantMsg   string
	}{
		{
			name: "valid",
			input: models.RequestSeller{
				Cid:         &validCid,
				CompanyName: &validCompany,
				Address:     &validAddress,
				Telephone:   &validTel,
				LocalityId:  &validLoc,
			},
			wantError: false,
		},
		{
			name:      "all missing",
			input:     models.RequestSeller{},
			wantError: true,
			wantMsg:   "Cid is required and must be greater than 0.",
		},
		{
			name: "invalid Cid",
			input: models.RequestSeller{
				Cid:         intPtr(0),
				CompanyName: &validCompany,
				Address:     &validAddress,
				Telephone:   &validTel,
				LocalityId:  &validLoc,
			},
			wantError: true,
			wantMsg:   "Cid is required and must be greater than 0.",
		},
		{
			name: "empty company",
			input: models.RequestSeller{
				Cid:         &validCid,
				CompanyName: strPtr(""),
				Address:     &validAddress,
				Telephone:   &validTel,
				LocalityId:  &validLoc,
			},
			wantError: true,
			wantMsg:   "CompanyName is required and cannot be empty.",
		},
		{
			name: "empty address",
			input: models.RequestSeller{
				Cid:         &validCid,
				CompanyName: &validCompany,
				Address:     strPtr(""),
				Telephone:   &validTel,
				LocalityId:  &validLoc,
			},
			wantError: true,
			wantMsg:   "Address is required and cannot be empty.",
		},
		{
			name: "empty telephone",
			input: models.RequestSeller{
				Cid:         &validCid,
				CompanyName: &validCompany,
				Address:     &validAddress,
				Telephone:   strPtr(""),
				LocalityId:  &validLoc,
			},
			wantError: true,
			wantMsg:   "Telephone is required and cannot be empty.",
		},
		{
			name: "empty locality",
			input: models.RequestSeller{
				Cid:         &validCid,
				CompanyName: &validCompany,
				Address:     &validAddress,
				Telephone:   &validTel,
				LocalityId:  strPtr(""),
			},
			wantError: true,
			wantMsg:   "Locality is required and must be greater than 0.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateSellerPost(tc.input)
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

func TestValidateSellerPatch(t *testing.T) {
	validCid := 1

	cases := []struct {
		name      string
		input     models.RequestSeller
		wantError bool
		wantMsg   string
	}{
		{
			name:      "all nil",
			input:     models.RequestSeller{},
			wantError: false,
		},
		{
			name: "valid one field",
			input: models.RequestSeller{
				Cid: &validCid,
			},
			wantError: false,
		},
		{
			name: "invalid cid",
			input: models.RequestSeller{
				Cid: intPtr(0),
			},
			wantError: true,
			wantMsg:   "Cid must be greater than 0.",
		},
		{
			name: "empty company",
			input: models.RequestSeller{
				CompanyName: strPtr(""),
			},
			wantError: true,
			wantMsg:   "CompanyName cannot be empty.",
		},
		{
			name: "empty address",
			input: models.RequestSeller{
				Address: strPtr(""),
			},
			wantError: true,
			wantMsg:   "Address cannot be empty.",
		},
		{
			name: "empty telephone",
			input: models.RequestSeller{
				Telephone: strPtr(""),
			},
			wantError: true,
			wantMsg:   "Telephone cannot be empty.",
		},
		{
			name: "empty locality",
			input: models.RequestSeller{
				LocalityId: strPtr(""),
			},
			wantError: true,
			wantMsg:   "Locality cannot be empty.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateSellerPatch(tc.input)
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

func TestValidateSellerPatchNotEmpty(t *testing.T) {
	validCid := 1
	validCompany := "Company"
	validAddress := "Some address"
	validTel := "1234"
	validLoc := "10"

	cases := []struct {
		name      string
		input     models.RequestSeller
		wantError bool
		wantMsg   string
	}{
		{
			name:      "all nil fields",
			input:     models.RequestSeller{},
			wantError: true,
			wantMsg:   "at least one field is required to be updated.",
		},
		{
			name:      "cid only",
			input:     models.RequestSeller{Cid: &validCid},
			wantError: false,
		},
		{
			name:      "company only",
			input:     models.RequestSeller{CompanyName: &validCompany},
			wantError: false,
		},
		{
			name:      "address only",
			input:     models.RequestSeller{Address: &validAddress},
			wantError: false,
		},
		{
			name:      "telephone only",
			input:     models.RequestSeller{Telephone: &validTel},
			wantError: false,
		},
		{
			name:      "locality only",
			input:     models.RequestSeller{LocalityId: &validLoc},
			wantError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateSellerPatchNotEmpty(tc.input)
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

// helpers
func strPtr(s string) *string {
	return &s
}
func intPtr(i int) *int {
	return &i
}
