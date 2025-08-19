package mappers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func TestRequestSellerToSeller(t *testing.T) {
	cid := 42
	company := "DemoCompany"
	address := "Test 555"
	telephone := "12345"
	locality := "ABCD"

	tests := []struct {
		name string
		in   models.RequestSeller
		want models.Seller
	}{
		{
			name: "simple mapping",
			in: models.RequestSeller{
				Cid:         &cid,
				CompanyName: &company,
				Address:     &address,
				Telephone:   &telephone,
				LocalityId:  &locality,
			},
			want: models.Seller{
				Id:          0,
				Cid:         cid,
				CompanyName: company,
				Address:     address,
				Telephone:   telephone,
				LocalityId:  locality,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mappers.RequestSellerToSeller(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestToResponseSeller(t *testing.T) {
	input := models.Seller{
		Id:          5,
		Cid:         42,
		CompanyName: "Company",
		Address:     "Addr",
		Telephone:   "Tel",
		LocalityId:  "Local",
	}
	want := models.ResponseSeller{
		Id:          5,
		Cid:         42,
		CompanyName: "Company",
		Address:     "Addr",
		Telephone:   "Tel",
		LocalityId:  "Local",
	}

	t.Run("returns correct response", func(t *testing.T) {
		got := mappers.ToResponseSeller(&input)
		assert.Equal(t, want, got)
	})
}

func TestToResponseSellerList(t *testing.T) {
	input := []models.Seller{
		{Id: 1, Cid: 1, CompanyName: "A", Address: "B", Telephone: "C", LocalityId: "D"},
		{Id: 2, Cid: 2, CompanyName: "X", Address: "Y", Telephone: "Z", LocalityId: "W"},
	}
	want := []models.ResponseSeller{
		{Id: 1, Cid: 1, CompanyName: "A", Address: "B", Telephone: "C", LocalityId: "D"},
		{Id: 2, Cid: 2, CompanyName: "X", Address: "Y", Telephone: "Z", LocalityId: "W"},
	}

	t.Run("converts list", func(t *testing.T) {
		got := mappers.ToResponseSellerList(input)
		assert.Equal(t, want, got)
	})
}

func TestApplySellerPatch(t *testing.T) {
	base := models.Seller{
		Id:          100,
		Cid:         11,
		CompanyName: "OldComp",
		Address:     "OldStreet",
		Telephone:   "111",
		LocalityId:  "LOC",
	}
	patchCid := 22
	patchCompany := "NewComp"
	patchAddress := "NewStreet"
	patchTelephone := "222"
	patchLocality := "NEWLOC"

	tests := []struct {
		name   string
		patch  models.RequestSeller
		expect models.Seller
	}{
		{
			name:  "patch cid",
			patch: models.RequestSeller{Cid: &patchCid},
			expect: models.Seller{
				Id:          100,
				Cid:         patchCid,
				CompanyName: "OldComp",
				Address:     "OldStreet",
				Telephone:   "111",
				LocalityId:  "LOC",
			},
		},
		{
			name:  "patch company",
			patch: models.RequestSeller{CompanyName: &patchCompany},
			expect: models.Seller{
				Id:          100,
				Cid:         11,
				CompanyName: patchCompany,
				Address:     "OldStreet",
				Telephone:   "111",
				LocalityId:  "LOC",
			},
		},
		{
			name: "patch all",
			patch: models.RequestSeller{
				Cid:         &patchCid,
				CompanyName: &patchCompany,
				Address:     &patchAddress,
				Telephone:   &patchTelephone,
				LocalityId:  &patchLocality,
			},
			expect: models.Seller{
				Id:          100,
				Cid:         patchCid,
				CompanyName: patchCompany,
				Address:     patchAddress,
				Telephone:   patchTelephone,
				LocalityId:  patchLocality,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset to base value for each run
			original := base
			mappers.ApplySellerPatch(&original, &tc.patch)
			assert.Equal(t, tc.expect, original)
		})
	}
}
