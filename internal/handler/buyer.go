package handler

import (
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type BuyerHandler struct {
	sv service.BuyerService
}

func NewBuyerHandler(sv service.BuyerService) *BuyerHandler {
	return &BuyerHandler{sv: sv}
}

func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var br models.RequestBuyer
	if err := httputil.DecodeJSON(r, &br); err != nil {
		response.Error(w, err)
		return
	}

	if err := validators.ValidateRequestBuyer(br); err != nil {
		response.Error(w, convertServiceErrorToAppError(err))
		return
	}

	b, err := h.sv.Create(ctx, br)
	if err != nil {
		response.Error(w, convertServiceErrorToAppError(err))
		return
	}

	response.JSON(w, http.StatusCreated, b)
}

func (h *BuyerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var br models.RequestBuyer
	if err := httputil.DecodeJSON(r, &br); err != nil {
		response.Error(w, err)
		return
	}

	if err := validators.ValidateUpdateBuyer(br); err != nil {
		response.Error(w, convertServiceErrorToAppError(err))
		return
	}
	if err := validators.ValidateBuyerPatchNotEmpty(br); err != nil {
		response.Error(w, convertServiceErrorToAppError(err))
		return
	}

	updated, err := h.sv.Update(ctx, id, br)
	if err != nil {
		response.Error(w, convertServiceErrorToAppError(err))
		return
	}

	response.JSON(w, http.StatusOK, updated)
}

func (h *BuyerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if handleError(w, err) {
		return
	}

	if err := h.sv.Delete(ctx, id); handleError(w, err) {
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *BuyerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.sv.FindAll(ctx)
	if handleError(w, err) {
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *BuyerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	b, err := h.sv.FindById(ctx, id)
	if err != nil {
		response.Error(w, convertServiceErrorToAppError(err))
		return
	}

	response.JSON(w, http.StatusOK, b)
}

// --- Utilidades ---

func handleError(w http.ResponseWriter, err error) bool {
	if err != nil {
		response.Error(w, convertServiceErrorToAppError(err))
		return true
	}
	return false
}

func convertServiceErrorToAppError(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *apperrors.AppError:
		return e
	case *api.ServiceError:
		return apperrors.NewAppError(mapServiceErrorCode(e.Code), e.Message)
	default:
		return apperrors.Wrap(err, "internal server error")
	}
}

func mapServiceErrorCode(code int) string {
	switch code {
	case http.StatusBadRequest:
		return apperrors.CodeBadRequest
	case http.StatusUnauthorized:
		return apperrors.CodeUnauthorized
	case http.StatusForbidden:
		return apperrors.CodeForbidden
	case http.StatusNotFound:
		return apperrors.CodeNotFound
	case http.StatusConflict:
		return apperrors.CodeConflict
	case http.StatusUnprocessableEntity:
		return apperrors.CodeValidationError
	default:
		return apperrors.CodeBadRequest
	}
}
