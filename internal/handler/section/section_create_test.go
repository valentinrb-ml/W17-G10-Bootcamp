package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/section"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	testhelpers "github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSectionHandler_CreateSection(t *testing.T) {
	type args struct {
		requestBody any
	}

	tests := []struct {
		name             string
		args             args
		mockService      func() *mocks.SectionServiceMock
		wantStatus       int
		wantResponseBody any
		wantErrorCode    string
		wantErrorMsgSub  string
	}{
		{
			name: "success",
			args: args{
				requestBody: testhelpers.DummySectionPost(1),
			},
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncCreate = func(ctx context.Context, sec models.Section) (*models.Section, error) {
					expected := testhelpers.DummySection(1)
					return &expected, nil
				}
				return mock
			},
			wantStatus:       http.StatusCreated,
			wantResponseBody: testhelpers.DummySection(1),
		},
		{
			name: "error - invalid request payload",
			args: args{
				requestBody: `{not-valid-json}`,
			},
			mockService: func() *mocks.SectionServiceMock {
				return &mocks.SectionServiceMock{}
			},
			wantStatus:      http.StatusBadRequest,
			wantErrorCode:   apperrors.CodeBadRequest,
			wantErrorMsgSub: "invalid JSON",
		},
		{
			name: "error - validation error",
			args: args{
				requestBody: models.PostSection{},
			},
			mockService: func() *mocks.SectionServiceMock {
				return &mocks.SectionServiceMock{}
			},
			wantStatus:      http.StatusUnprocessableEntity,
			wantErrorCode:   apperrors.CodeValidationError,
			wantErrorMsgSub: "required",
		},
		{
			name: "error - service error",
			args: args{
				requestBody: testhelpers.DummySectionPost(1),
			},
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncCreate = func(ctx context.Context, sec models.Section) (*models.Section, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "unexpected error")
				}
				return mock
			},
			wantStatus:      http.StatusInternalServerError,
			wantErrorCode:   apperrors.CodeInternal,
			wantErrorMsgSub: "unexpected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBodyBytes []byte
			switch val := tt.args.requestBody.(type) {
			case string:
				requestBodyBytes = []byte(val)
			default:
				b, err := json.Marshal(tt.args.requestBody)
				require.NoError(t, err)
				requestBodyBytes = b
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/sections", bytes.NewReader(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h := handler.NewSectionHandler(tt.mockService())
			// inject test logger to cover logger != nil branches
			h.SetLogger(testhelpers.NewTestLogger())

			// Call handler
			h.CreateSection(rec, req)

			// Check status
			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusCreated {
				var responseEnvelope struct {
					Data models.Section `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &responseEnvelope)
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseBody, responseEnvelope.Data)
			} else {
				var body struct {
					Error struct {
						Code    string         `json:"code"`
						Message string         `json:"message"`
						Details map[string]any `json:"details"`
					} `json:"error"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.wantErrorCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrorMsgSub)
			}
		})
	}
}
