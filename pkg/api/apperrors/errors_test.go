package apperrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
)

func TestNewAppError_KnownAndUnknownCode(t *testing.T) {
	// Debe encontrar el status del code conocido
	appErr := apperrors.NewAppError(apperrors.CodeInternal, "err")
	require.Equal(t, apperrors.CodeInternal, appErr.Code)
	require.Equal(t, "err", appErr.Message)
	require.Greater(t, appErr.HTTPStatus, 0)
	require.Empty(t, appErr.Details)

	// Code desconocido: debe quedarse con status 500
	unknown := apperrors.NewAppError("FAKE_CODE_X", "nope")
	require.Equal(t, 500, unknown.HTTPStatus)
	require.Equal(t, "FAKE_CODE_X", unknown.Code)
	require.Equal(t, "nope", unknown.Message)
}

func TestAppError_ErrorMethod(t *testing.T) {
	appErr := apperrors.NewAppError(apperrors.CodeInternal, "msg here")
	require.Contains(t, appErr.Error(), appErr.Code)
	require.Contains(t, appErr.Error(), "msg here")
}

func TestAppError_WithDetail_IsImmutable(t *testing.T) {
	e1 := apperrors.NewAppError("SOME", "test with details")
	e2 := e1.WithDetail("field", "value")

	require.NotSame(t, e1, e2, "should return a new instance")
	require.Empty(t, e1.Details)
	require.Equal(t, "value", e2.Details["field"])

	// Con más detalles, viejo debe seguir sin cambios
	e3 := e2.WithDetail("other", 123)
	require.NotSame(t, e2, e3)
	require.Equal(t, "value", e3.Details["field"])
	require.Equal(t, 123, e3.Details["other"])
	require.Empty(t, e1.Details)
}

func TestAppError_Wrap(t *testing.T) {
	// Nil retorna nil
	require.Nil(t, apperrors.Wrap(nil, "ignore"))

	// Ya es app error: retorna el mismo
	orig := apperrors.NewAppError(apperrors.CodeInternal, "msg")
	result := apperrors.Wrap(orig, "x")
	require.Equal(t, orig, result)

	// Error de Go: debe devolver AppError nuevo
	ext := errors.New("foo error")
	wrapped := apperrors.Wrap(ext, "wrapped!")
	require.NotNil(t, wrapped)
	require.Equal(t, apperrors.CodeInternal, wrapped.Code)
	require.Equal(t, "wrapped!", wrapped.Message)
}

func TestIsAppError(t *testing.T) {
	appErr := apperrors.NewAppError("CODE1", "x")
	normal := errors.New("not app error")

	require.True(t, apperrors.IsAppError(appErr, "CODE1"))
	require.False(t, apperrors.IsAppError(normal, "CODE1"))
	require.False(t, apperrors.IsAppError(appErr, "CODE_ERR"))
}
