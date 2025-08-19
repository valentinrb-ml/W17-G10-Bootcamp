package httputil_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	httputil "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
)

func setChiParam(req *http.Request, name, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(name, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}

// Helper para decodificar body de error
func decodeErrBody(t *testing.T, rec *httptest.ResponseRecorder) response.ErrorResponse {
	var errResp response.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	return errResp
}

func TestMethodNotAllowedHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)
	httputil.MethodNotAllowedHandler(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	errResp := decodeErrBody(t, rec)
	require.Equal(t, apperrors.CodeMethodNotAllowed, errResp.Error.Code)
	require.Contains(t, errResp.Error.Message, "method not allowed")
}

func TestNotFoundHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/notfound", nil)
	httputil.NotFoundHandler(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
	errResp := decodeErrBody(t, rec)
	require.Equal(t, apperrors.CodeNotFound, errResp.Error.Code)
	require.Contains(t, errResp.Error.Message, "endpoint not found")
}

func TestDecodeJSON(t *testing.T) {
	type myJSON struct {
		Name string `json:"name"`
		Num  int    `json:"num"`
	}
	t.Run("fail if body is nil", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", nil)
		req.Body = nil
		var dest myJSON
		err := httputil.DecodeJSON(req, &dest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "request body is required")
	})
	t.Run("fail on invalid JSON", func(t *testing.T) {
		bad := strings.NewReader("not json!")
		req := httptest.NewRequest("POST", "/", bad)
		var dest myJSON
		err := httputil.DecodeJSON(req, &dest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid JSON format")
	})
	t.Run("successfully decodes good JSON", func(t *testing.T) {
		data := map[string]any{"name": "ok", "num": 42}
		buf, _ := json.Marshal(data)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(buf))
		var dest myJSON
		err := httputil.DecodeJSON(req, &dest)
		require.NoError(t, err)
		require.Equal(t, "ok", dest.Name)
		require.Equal(t, 42, dest.Num)
	})
}

func TestParseIDParam(t *testing.T) {
	t.Run("returns error if param is missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foo", nil)
		id, err := httputil.ParseIDParam(req, "id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "id parameter is required")
		require.Equal(t, 0, id)
	})

	t.Run("returns error if param is not an int", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foo/abc", nil)
		req = setChiParam(req, "id", "abc")
		id, err := httputil.ParseIDParam(req, "id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be a valid integer")
		require.Equal(t, 0, id)
	})

	t.Run("returns error if param is zero or negative", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foo/0", nil)
		req = setChiParam(req, "id", "0")
		id, err := httputil.ParseIDParam(req, "id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be a positive integer")
		require.Equal(t, 0, id)

		req = setChiParam(req, "id", "-5")
		id, err = httputil.ParseIDParam(req, "id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be a positive integer")
		require.Equal(t, 0, id)
	})

	t.Run("success with positive int param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foo/13", nil)
		req = setChiParam(req, "id", "13")
		id, err := httputil.ParseIDParam(req, "id")
		require.NoError(t, err)
		require.Equal(t, 13, id)
	})
}
