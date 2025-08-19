package request_test

import (
	"github.com/stretchr/testify/require"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type dummyStruct struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		expectErr   error
		wantVal     *dummyStruct
	}{
		{
			name:        "fails if not application/json",
			contentType: "text/plain",
			body:        `{"name":"Go","count":42}`,
			expectErr:   request.ErrRequestContentTypeNotJSON,
			wantVal:     nil,
		},
		{
			name:        "fails if bad JSON",
			contentType: "application/json",
			body:        `not valid json!`,
			expectErr:   request.ErrRequestJSONInvalid,
			wantVal:     nil,
		},
		{
			name:        "decodes valid JSON",
			contentType: "application/json",
			body:        `{"name":"Team","count":99}`,
			expectErr:   nil,
			wantVal:     &dummyStruct{Name: "Team", Count: 99},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/endpoint", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)

			var got dummyStruct
			err := request.JSON(req, &got)
			if tc.expectErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantVal, &got)
			}
		})
	}
}
