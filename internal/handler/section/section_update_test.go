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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSectionHandler_UpdateSection(t *testing.T) {

	dummyPatch := testhelpers.DummySectionPatch(1)
	dummy := testhelpers.DummySection(1)
	dummyUpdated := testhelpers.DummySection(1)

	tests := []struct {
		name             string
		inputID          string
		requestBody      any
		mockService      func() *mocks.SectionServiceMock
		wantStatus       int
		wantResponseBody any
		wantErrCode      string
		wantErrMsgSub    string
	}{
		{
			name:        "success: updates section",
			inputID:     "1",
			requestBody: dummyPatch,
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncUpdate = func(ctx context.Context, id int, patch models.PatchSection) (*models.Section, error) {
					s := dummyUpdated
					return &s, nil
				}
				mock.FuncFindById = func(ctx context.Context, id int) (*models.Section, error) {
					s := dummy
					return &s, nil
				}
				return mock
			},
			wantStatus:       http.StatusOK,
			wantResponseBody: testhelpers.DummyResponseSection(1),
		},
		{
			name:        "error: invalid id param",
			inputID:     "abc",
			requestBody: dummyPatch,
			mockService: func() *mocks.SectionServiceMock {
				return &mocks.SectionServiceMock{}
			},
			wantStatus:    http.StatusBadRequest,
			wantErrCode:   apperrors.CodeBadRequest,
			wantErrMsgSub: "id must be a valid integer",
		},
		{
			name:        "error: invalid JSON body",
			inputID:     "123",
			requestBody: "{notajson",
			mockService: func() *mocks.SectionServiceMock {
				return &mocks.SectionServiceMock{}
			},
			wantStatus:    http.StatusBadRequest,
			wantErrCode:   apperrors.CodeBadRequest,
			wantErrMsgSub: "invalid JSON",
		},
		{
			name:        "error: validation error",
			inputID:     "123",
			requestBody: models.PatchSection{},
			mockService: func() *mocks.SectionServiceMock {
				return &mocks.SectionServiceMock{}
			},
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrCode:   apperrors.CodeValidationError,
			wantErrMsgSub: "At least one field must be provided to update the section.",
		},
		{
			name:        "error: service error",
			inputID:     "123",
			requestBody: dummyPatch,
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncUpdate = func(ctx context.Context, id int, patch models.PatchSection) (*models.Section, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "internal error")
				}
				return mock
			},
			wantStatus:    http.StatusInternalServerError,
			wantErrCode:   apperrors.CodeInternal,
			wantErrMsgSub: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBodyBytes []byte
			switch val := tt.requestBody.(type) {
			case string:
				requestBodyBytes = []byte(val)
			default:
				b, err := json.Marshal(tt.requestBody)
				require.NoError(t, err)
				requestBodyBytes = b
			}
			req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/"+tt.inputID, bytes.NewReader(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h := handler.NewSectionHandler(tt.mockService())

			req = testhelpers.SetChiURLParam(req, "id", tt.inputID)

			h.UpdateSection(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var envelope struct {
					Data models.ResponseSection `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &envelope)
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseBody, envelope.Data)
			} else {
				var body struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.wantErrCode, body.Error.Code)
				require.Contains(t, body.Error.Message, tt.wantErrMsgSub)
			}
		})
	}
}
