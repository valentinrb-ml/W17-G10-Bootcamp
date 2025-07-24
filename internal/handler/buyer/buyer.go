package handler

import (
	"net/http"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
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
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"))
		return
	}

	//fmt.Printf("Request received: %+v\n", br)
	if err := validators.ValidateRequestBuyer(br); err != nil {
		response.Error(w, err)
		return
	}

	b, err := h.sv.Create(ctx, br)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, b)
}

func (h *BuyerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid ID parameter"))
		return
	}

	var br models.RequestBuyer
	if err := httputil.DecodeJSON(r, &br); err != nil {
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid request body"))
		return
	}

	if err := validators.ValidateUpdateBuyer(br); err != nil {
		response.Error(w, err)
		return
	}
	if err := validators.ValidateBuyerPatchNotEmpty(br); err != nil {
		response.Error(w, err)
		return
	}

	updated, err := h.sv.Update(ctx, id, br)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, updated)
}

func (h *BuyerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid ID parameter"))
		return
	}

	if err := h.sv.Delete(ctx, id); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *BuyerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.sv.FindAll(ctx)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *BuyerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid ID parameter"))
		return
	}

	b, err := h.sv.FindById(ctx, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, b)
}
