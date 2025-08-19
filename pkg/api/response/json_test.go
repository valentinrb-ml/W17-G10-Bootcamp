package response_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name          string
		code          int
		input         any
		wantStatus    int
		wantData      any
		skipTypeCheck bool // para tipo no comparable, ej: map
	}{
		{
			name:       "returns json with data object",
			code:       http.StatusOK,
			input:      map[string]any{"foo": "bar"},
			wantStatus: http.StatusOK,
			wantData:   map[string]any{"foo": "bar"},
		},
		{
			name:       "response with slice",
			code:       http.StatusCreated,
			input:      []string{"a", "b", "c"},
			wantStatus: http.StatusCreated,
			wantData:   []any{"a", "b", "c"}, // por tipo JSON decode es []any, no []string
		},
		{
			name:       "returns only status if body is nil",
			code:       http.StatusAccepted,
			input:      nil,
			wantStatus: http.StatusAccepted,
			wantData:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			response.JSON(rec, tc.code, tc.input)
			require.Equal(t, tc.wantStatus, rec.Code)

			if tc.input == nil {
				// No JSON, no Content-Type
				require.Empty(t, rec.Body.String())
				return
			}

			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			// Desempaqueta envelope y compara Data
			var resp struct {
				Data any `json:"data"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			require.NoError(t, err)

			// Si es slice/map, compara con ElementsMatch o omite tipo
			if m, ok := tc.wantData.([]any); ok {
				require.ElementsMatch(t, m, resp.Data.([]any))
			} else if mp, ok := tc.wantData.(map[string]any); ok {
				require.Equal(t, mp, resp.Data)
			} else {
				require.Equal(t, tc.wantData, resp.Data)
			}
		})
	}
}
