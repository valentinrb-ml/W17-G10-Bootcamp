package repository_test

import (
	"context"
	"testing"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

/*
Create() with an already-expired context must propagate
an INTERNAL_ERROR (wrapped deadline exceeded).
*/
func TestCreate_ContextTimeout(t *testing.T) {
	t.Parallel()

	repo, _, cleanup := testhelpers.NewProductRecordRepoMock(t)
	defer cleanup()

	// ctx expires instantly
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	_, err := repo.Create(ctx, testhelpers.BuildProductRecord())
	testhelpers.RequireAppErr(t, err, apperrors.CodeInternal)
}
