package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	productmock "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
)

func withIDPatch(r *http.Request, id string) *http.Request {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestProductHandler_Patch(t *testing.T) {
	t.Parallel()

	code := "ZZ"
	reqPayload := models.ProductPatchRequest{ProductCode: &code}
	bodyOk, _ := json.Marshal(reqPayload)
	resp := models.ProductResponse{ID: 4}

	tests := []struct {
		name      string
		body      []byte
		param     string
		mockSetup func(*productmock.MockService)
		status    int
		appCode   string
	}{
		{
			name:  "success",
			body:  bodyOk,
			param: "4",
			mockSetup: func(s *productmock.MockService) {
				s.On("Patch", mock.Anything, 4, reqPayload).Return(resp, nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name:      "bad json",
			body:      []byte(`{"bad":`),
			param:     "4",
			status:    http.StatusBadRequest,
			appCode:   apperrors.CodeBadRequest,
			mockSetup: func(s *productmock.MockService) {},
		},
		{
			name:    "service conflict",
			body:    bodyOk,
			param:   "4",
			status:  http.StatusConflict,
			appCode: apperrors.CodeConflict,
			mockSetup: func(s *productmock.MockService) {
				s.On("Patch", mock.Anything, 4, reqPayload).
					Return(models.ProductResponse{}, apperrors.NewAppError(apperrors.CodeConflict, "x")).Once()
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := &productmock.MockService{}
			tc.mockSetup(svc)

			h := handler.NewProductHandler(svc)
			h.SetLogger(testhelpers.NewTestLogger())
			req := testhelpers.NewRequest(t, http.MethodPatch, "/products/"+tc.param, bytes.NewReader(tc.body))
			req = withIDPatch(req, tc.param)

			rec := testhelpers.DoRawRequest(t, req, http.HandlerFunc(h.Patch))

			require.Equal(t, tc.status, rec.Code)
			if tc.appCode != "" {
				app, _ := testhelpers.DecodeAppErr(rec.Body)
				testhelpers.RequireAppErr(t, app, tc.appCode)
			}
			svc.AssertExpectations(t)
		})
	}
}
