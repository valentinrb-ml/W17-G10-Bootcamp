package mappers_test

import (
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestRequestSectionToSection(t *testing.T) {
	t.Run("maps all fields from PostSection (request) to Section (model)", func(t *testing.T) {
		// Given
		in := testhelpers.DummySectionPost(1)
		want := models.Section{
			SectionNumber:      in.SectionNumber,
			CurrentTemperature: *in.CurrentTemperature,
			MinimumTemperature: *in.MinimumTemperature,
			CurrentCapacity:    in.CurrentCapacity,
			MinimumCapacity:    in.MinimumCapacity,
			MaximumCapacity:    in.MaximumCapacity,
			WarehouseId:        in.WarehouseId,
			ProductTypeId:      in.ProductTypeId,
		}

		// When
		got := mappers.RequestSectionToSection(in)

		// Then
		require.Equal(t, want, got)
	})

	t.Run("panics if CurrentTemperature is nil", func(t *testing.T) {
		in := testhelpers.DummySectionPost(1)
		in.CurrentTemperature = nil

		require.Panics(t, func() {
			_ = mappers.RequestSectionToSection(in)
		})
	})
}

func TestSectionToResponseSection(t *testing.T) {
	t.Run("maps all fields from Section (model) to ResponseSection", func(t *testing.T) {
		in := testhelpers.DummySection(1)
		want := testhelpers.DummyResponseSection(1)

		got := mappers.SectionToResponseSection(in)

		require.Equal(t, want, got)
	})
}

func TestApplySectionPatch(t *testing.T) {
	type args struct {
		patch    models.PatchSection
		original models.Section
		expect   models.Section
	}

	tests := []args{
		{
			patch: models.PatchSection{
				SectionNumber:      testhelpers.IntPtr(900),
				CurrentTemperature: testhelpers.Float64Ptr(99.5),
				CurrentCapacity:    testhelpers.IntPtr(888),
				MaximumCapacity:    testhelpers.IntPtr(5000),
				MinimumCapacity:    testhelpers.IntPtr(111),
				MinimumTemperature: testhelpers.Float64Ptr(42.2),
				ProductTypeId:      testhelpers.IntPtr(77),
				WarehouseId:        testhelpers.IntPtr(66),
			},
			original: testhelpers.DummySection(1),
			expect: models.Section{
				Id:                 1,
				SectionNumber:      900,
				CurrentTemperature: 99.5,
				MinimumTemperature: 42.2,
				CurrentCapacity:    888,
				MinimumCapacity:    111,
				MaximumCapacity:    5000,
				ProductTypeId:      77,
				WarehouseId:        66,
			},
		},
		{
			patch: func() models.PatchSection {
				ps := models.PatchSection{}
				ps.CurrentTemperature = testhelpers.Float64Ptr(1234)
				return ps
			}(),
			original: testhelpers.DummySection(1),
			expect:   func() models.Section { s := testhelpers.DummySection(1); s.CurrentTemperature = 1234; return s }(),
		},
		// Caso: patch vacío no cambia nada
		{
			patch:    models.PatchSection{},
			original: testhelpers.DummySection(1),
			expect:   testhelpers.DummySection(1),
		},
	}

	for _, tt := range tests {
		t.Run("patch should update correct fields", func(t *testing.T) {
			orig := tt.original
			mappers.ApplySectionPatch(tt.patch, &orig)
			require.Equal(t, tt.expect, orig)
		})
	}
}
