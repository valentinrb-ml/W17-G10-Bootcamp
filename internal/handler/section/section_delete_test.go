package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/section"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSectionHandler_DeleteSection(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		mockService   func() *mocks.SectionServiceMock
		wantStatus    int
		wantErrCode   string
		wantErrMsgSub string
	}{
		{
			name:    "success: deletes section",
			inputID: "1",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncDelete = func(ctx context.Context, id int) error {
					return nil
				}
				return mock
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:    "error: invalid id param",
			inputID: "bad",
			mockService: func() *mocks.SectionServiceMock {
				return &mocks.SectionServiceMock{}
			},
			wantStatus:    http.StatusBadRequest,
			wantErrCode:   apperrors.CodeBadRequest,
			wantErrMsgSub: "id must be a valid integer",
		},
		{
			name:    "error: not found",
			inputID: "2",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncDelete = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeNotFound, "section not found")
				}
				return mock
			},
			wantStatus:    http.StatusNotFound,
			wantErrCode:   apperrors.CodeNotFound,
			wantErrMsgSub: "section not found",
		},
		{
			name:    "error: service error",
			inputID: "7",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncDelete = func(ctx context.Context, id int) error {
					return apperrors.NewAppError(apperrors.CodeInternal, "internal error")
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
			req, err := http.NewRequest(http.MethodDelete, "/api/v1/sections/"+tt.inputID, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			req = testhelpers.SetChiURLParam(req, "id", tt.inputID)
			h := handler.NewSectionHandler(tt.mockService())

			h.DeleteSection(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus != http.StatusNoContent {
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
			} else {
				require.Empty(t, rec.Body.String())
			}
		})
	}
}
