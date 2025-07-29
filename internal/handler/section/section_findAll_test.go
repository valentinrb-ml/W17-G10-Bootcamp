package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func TestSectionHandler_FindAllSections(t *testing.T) {
	sections := []models.Section{
		testhelpers.DummySection(1),
		testhelpers.DummySection(2),
	}
	expectedResp := []models.ResponseSection{
		mappers.SectionToResponseSection(sections[0]),
		mappers.SectionToResponseSection(sections[1]),
	}

	tests := []struct {
		name        string
		mockService func() *mocks.SectionServiceMock
		wantStatus  int
		wantBody    []models.ResponseSection
		wantErrCode string
		wantErrSub  string
	}{
		{
			name: "success: returns all sections",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncFindAll = func(ctx context.Context) ([]models.Section, error) {
					return sections, nil
				}
				return mock
			},
			wantStatus: http.StatusOK,
			wantBody:   expectedResp,
		},
		{
			name: "error: service failure",
			mockService: func() *mocks.SectionServiceMock {
				mock := &mocks.SectionServiceMock{}
				mock.FuncFindAll = func(ctx context.Context) ([]models.Section, error) {
					return nil, apperrors.NewAppError(apperrors.CodeInternal, "db failure")
				}
				return mock
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: apperrors.CodeInternal,
			wantErrSub:  "db failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/sections", nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			h := handler.NewSectionHandler(tt.mockService())

			h.FindAllSections(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var envelope struct {
					Data []models.ResponseSection `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &envelope)
				require.NoError(t, err)
				require.Equal(t, tt.wantBody, envelope.Data)
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
				require.Contains(t, body.Error.Message, tt.wantErrSub)
			}
		})
	}
}
