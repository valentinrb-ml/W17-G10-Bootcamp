package validators_test

import (
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	section "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
	"testing"
)

func TestValidateSectionRequest(t *testing.T) {
	testCases := []struct {
		name       string
		input      section.PostSection
		wantErrMsg string
	}{
		{
			name:       "valid section returns nil error",
			input:      testhelpers.DummySectionPost(1),
			wantErrMsg: "",
		},
		{
			name: "fails if a required int field is zero",
			input: func() section.PostSection {
				ps := testhelpers.DummySectionPost(1)
				ps.SectionNumber = 0
				return ps
			}(),
			wantErrMsg: "All fields are required",
		},
		{
			name: "fails if a required pointer field is nil",
			input: func() section.PostSection {
				ps := testhelpers.DummySectionPost(1)
				ps.MinimumTemperature = nil
				return ps
			}(),
			wantErrMsg: "All fields are required",
		},
		{
			name: "capacity negative (maximum)",
			input: func() section.PostSection {
				ps := testhelpers.DummySectionPost(1)
				ps.MaximumCapacity = -5
				return ps
			}(),
			wantErrMsg: "Capacity values cannot be negative",
		},
		{
			name: "max less than min",
			input: func() section.PostSection {
				ps := testhelpers.DummySectionPost(1)
				ps.MaximumCapacity = 5
				ps.MinimumCapacity = 10
				return ps
			}(),
			wantErrMsg: "Maximum capacity cannot be less than minimum capacity",
		},
		{
			name: "current exceeds max",
			input: func() section.PostSection {
				ps := testhelpers.DummySectionPost(1)
				ps.CurrentCapacity = 200
				return ps
			}(),
			wantErrMsg: "Current capacity cannot exceed maximum capacity",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateSectionRequest(tc.input)
			if tc.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Contains(t, appErr.Message, tc.wantErrMsg)
			}
		})
	}
}

func TestValidateSectionPatch(t *testing.T) {
	testCases := []struct {
		name       string
		input      section.PatchSection
		wantErrMsg string
	}{
		{
			name:       "valid patch returns nil error",
			input:      testhelpers.DummySectionPatch(1),
			wantErrMsg: "",
		},
		{
			name:       "all fields nil should fail",
			input:      section.PatchSection{},
			wantErrMsg: "At least one field must be provided",
		},
		{
			name: "negative maximum capacity",
			input: func() section.PatchSection {
				ps := testhelpers.DummySectionPatch(1)
				ps.MaximumCapacity = testhelpers.IntPtr(-5)
				return ps
			}(),
			wantErrMsg: "Maximum capacity cannot be negative",
		},
		{
			name: "negative minimum capacity",
			input: func() section.PatchSection {
				ps := testhelpers.DummySectionPatch(1)
				ps.MinimumCapacity = testhelpers.IntPtr(-10)
				return ps
			}(),
			wantErrMsg: "Minimum capacity cannot be negative",
		},
		{
			name: "negative current capacity",
			input: func() section.PatchSection {
				ps := testhelpers.DummySectionPatch(1)
				ps.CurrentCapacity = testhelpers.IntPtr(-1)
				return ps
			}(),
			wantErrMsg: "Current capacity cannot be negative",
		},
		{
			name: "max less than min",
			input: func() section.PatchSection {
				ps := testhelpers.DummySectionPatch(1)
				max, min := 5, 10
				ps.MaximumCapacity = &max
				ps.MinimumCapacity = &min
				return ps
			}(),
			wantErrMsg: "Maximum capacity cannot be less than minimum capacity",
		},
		{
			name: "current greater than max",
			input: func() section.PatchSection {
				ps := testhelpers.DummySectionPatch(1)
				max, curr := 50, 60
				ps.MaximumCapacity = &max
				ps.CurrentCapacity = &curr
				return ps
			}(),
			wantErrMsg: "Current capacity cannot exceed maximum capacity",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateSectionPatch(tc.input)
			if tc.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				appErr, ok := err.(*apperrors.AppError)
				require.True(t, ok)
				require.Contains(t, appErr.Message, tc.wantErrMsg)
			}
		})
	}
}
