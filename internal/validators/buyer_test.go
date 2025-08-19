package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// helpers para punteros
func strptr(s string) *string { return &s }

func TestValidateRequestBuyer_AllValid(t *testing.T) {
	r := models.RequestBuyer{
		CardNumberId: strptr("123"),
		FirstName:    strptr("Juan"),
		LastName:     strptr("Perez"),
	}
	err := validators.ValidateRequestBuyer(r)
	assert.Nil(t, err)
}

func TestValidateRequestBuyer_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		req  models.RequestBuyer
	}{
		{"NilCardNumberId", models.RequestBuyer{CardNumberId: nil, FirstName: strptr("a"), LastName: strptr("b")}},
		{"NilFirstName", models.RequestBuyer{CardNumberId: strptr("1"), FirstName: nil, LastName: strptr("b")}},
		{"NilLastName", models.RequestBuyer{CardNumberId: strptr("1"), FirstName: strptr("a"), LastName: nil}},
		{"EmptyCardNumberId", models.RequestBuyer{CardNumberId: strptr(""), FirstName: strptr("a"), LastName: strptr("b")}},
		{"EmptyFirstName", models.RequestBuyer{CardNumberId: strptr("1"), FirstName: strptr(""), LastName: strptr("b")}},
		{"EmptyLastName", models.RequestBuyer{CardNumberId: strptr("1"), FirstName: strptr("a"), LastName: strptr("")}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateRequestBuyer(tc.req)
			assert.NotNil(t, err)
			assert.Equal(t, apperrors.CodeValidationError, err.Code)
		})
	}
}

func TestValidateUpdateBuyer_Valid(t *testing.T) {
	r := models.RequestBuyer{
		CardNumberId: strptr("123"),
		FirstName:    strptr("Juan"),
		LastName:     strptr("Perez"),
	}
	err := validators.ValidateUpdateBuyer(r)
	assert.Nil(t, err)
}

func TestValidateUpdateBuyer_EmptyFields(t *testing.T) {
	type tc struct {
		name     string
		r        models.RequestBuyer
		expected string
	}
	tests := []tc{
		{"EmptyCardNumberId", models.RequestBuyer{CardNumberId: strptr("")}, "id_card_number cannot be empty"},
		{"EmptyFirstName", models.RequestBuyer{FirstName: strptr("")}, "first_name cannot be empty"},
		{"EmptyLastName", models.RequestBuyer{LastName: strptr("")}, "last_name cannot be empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validators.ValidateUpdateBuyer(tt.r)
			assert.NotNil(t, err)
			assert.Equal(t, apperrors.CodeValidationError, err.Code)
			assert.Equal(t, tt.expected, err.Message)
		})
	}
}

func TestValidateUpdateBuyer_AllNil(t *testing.T) {
	r := models.RequestBuyer{}
	err := validators.ValidateUpdateBuyer(r)
	assert.Nil(t, err)
}

func TestValidateBuyerPatchNotEmpty_AllNil(t *testing.T) {
	r := models.RequestBuyer{}
	err := validators.ValidateBuyerPatchNotEmpty(r)
	assert.NotNil(t, err)
	assert.Contains(t, err.Message, "At least one field is required")
}

func TestValidateBuyerPatchNotEmpty_AtLeastOnePresent(t *testing.T) {
	tests := []models.RequestBuyer{
		{CardNumberId: strptr("1")},
		{FirstName: strptr("a")},
		{LastName: strptr("b")},
		{CardNumberId: strptr("1"), FirstName: strptr("a")},
	}
	for i, r := range tests {
		t.Run("CaseOk#"+string(rune(i+'0')), func(t *testing.T) {
			err := validators.ValidateBuyerPatchNotEmpty(r)
			assert.Nil(t, err)
		})
	}
}
