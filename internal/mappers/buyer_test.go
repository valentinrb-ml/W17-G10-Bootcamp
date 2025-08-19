package mappers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

func strPtr(s string) *string { return &s }

func TestRequestBuyerToBuyer_AllFieldsPresent(t *testing.T) {
	req := models.RequestBuyer{
		CardNumberId: strPtr("CARD-777"),
		FirstName:    strPtr("Fede"),
		LastName:     strPtr("Laguardia"),
	}
	got := mappers.RequestBuyerToBuyer(req)
	assert.Equal(t, "CARD-777", got.CardNumberId)
	assert.Equal(t, "Fede", got.FirstName)
	assert.Equal(t, "Laguardia", got.LastName)
}

func TestRequestBuyerToBuyer_FieldsNil(t *testing.T) {
	req := models.RequestBuyer{
		CardNumberId: nil,
		FirstName:    nil,
		LastName:     nil,
	}
	got := mappers.RequestBuyerToBuyer(req)
	assert.Equal(t, "", got.CardNumberId)
	assert.Equal(t, "", got.FirstName)
	assert.Equal(t, "", got.LastName)
}

func TestApplyBuyerPatch(t *testing.T) {
	original := &models.Buyer{Id: 2, CardNumberId: "ORIG", FirstName: "OLD", LastName: "SURNAME"}
	patch := models.RequestBuyer{CardNumberId: strPtr("PATCHED-ID")}
	mappers.ApplyBuyerPatch(original, &patch)
	assert.Equal(t, "PATCHED-ID", original.CardNumberId)
	assert.Equal(t, "OLD", original.FirstName)
	assert.Equal(t, "SURNAME", original.LastName)

	patch2 := models.RequestBuyer{FirstName: strPtr("NEWNAME")}
	mappers.ApplyBuyerPatch(original, &patch2)
	assert.Equal(t, "NEWNAME", original.FirstName)

	patch3 := models.RequestBuyer{LastName: strPtr("NEWLAST")}
	mappers.ApplyBuyerPatch(original, &patch3)
	assert.Equal(t, "NEWLAST", original.LastName)

	patchNone := models.RequestBuyer{}
	mappers.ApplyBuyerPatch(original, &patchNone)
	assert.Equal(t, "PATCHED-ID", original.CardNumberId)
	assert.Equal(t, "NEWNAME", original.FirstName)
	assert.Equal(t, "NEWLAST", original.LastName)
}

func TestToResponseBuyer(t *testing.T) {
	b := &models.Buyer{
		Id:           7,
		CardNumberId: "FOO",
		FirstName:    "X",
		LastName:     "Y",
	}
	r := mappers.ToResponseBuyer(b)
	expected := models.ResponseBuyer{
		Id:           7,
		CardNumberId: "FOO",
		FirstName:    "X",
		LastName:     "Y",
	}
	assert.Equal(t, expected, r)
}

func TestToResponseBuyerList(t *testing.T) {
	buyers := []models.Buyer{
		{Id: 1, CardNumberId: "A1", FirstName: "A", LastName: "L"},
		{Id: 2, CardNumberId: "B2", FirstName: "B", LastName: "M"},
	}
	got := mappers.ToResponseBuyerList(buyers)
	expected := []models.ResponseBuyer{
		{Id: 1, CardNumberId: "A1", FirstName: "A", LastName: "L"},
		{Id: 2, CardNumberId: "B2", FirstName: "B", LastName: "M"},
	}
	assert.Equal(t, expected, got)
}
