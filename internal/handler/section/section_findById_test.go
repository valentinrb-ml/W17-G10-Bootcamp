package handler_test

import (
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

func TestSectionHandler_FindById(t *testing.T) {
	sectionDummy := testhelpers.DummySection(1)
	respDummy := testhelpers.DummyResponseSection(1)

	tests := []struct {
		name             string
		inputID          string
		mockService      func() *mocks.SectionServiceMock
		wantStatus       int
		wantResponseBody any
		wantErrCode      string
		wantErrMsgSub    string
	}{
		{
			name:    "success: finds section by id",
			inputID: "1",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncFindById = func(ctx context.Context, id int) (*models.Section, error) {
					return &sectionDummy, nil
				}
				return mock
			},
			wantStatus:       http.StatusOK,
			wantResponseBody: respDummy,
		},
		{
			name:    "error: invalid id param",
			inputID: "abc",
			mockService: func() *mocks.SectionServiceMock {
				return &mocks.SectionServiceMock{}
			},
			wantStatus:    http.StatusBadRequest,
			wantErrCode:   apperrors.CodeBadRequest,
			wantErrMsgSub: "id must be a valid integer",
		},
		{
			name:    "error: section not found",
			inputID: "2",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncFindById = func(ctx context.Context, id int) (*models.Section, error) {
					return nil, apperrors.NewAppError(apperrors.CodeNotFound, "section not found")
				}
				return mock
			},
			wantStatus:    http.StatusNotFound,
			wantErrCode:   apperrors.CodeNotFound,
			wantErrMsgSub: "section not found",
		},
		{
			name:    "error: service error",
			inputID: "1",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncFindById = func(ctx context.Context, id int) (*models.Section, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "unexpected error")
				}
				return mock
			},
			wantStatus:    http.StatusInternalServerError,
			wantErrCode:   apperrors.CodeInternal,
			wantErrMsgSub: "unexpected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/"+tt.inputID, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			req = testhelpers.SetChiURLParam(req, "id", tt.inputID)

			h := handler.NewSectionHandler(tt.mockService())
			h.SetLogger(testhelpers.NewTestLogger())

			h.FindById(rec, req)

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
