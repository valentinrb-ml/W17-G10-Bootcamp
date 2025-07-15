package handler

import (
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func NewSellerHandler(sv service.SellerService) *SellerHandler {
	return &SellerHandler{
		sv: sv,
	}
}

type SellerHandler struct {
	sv service.SellerService
}

func (h *SellerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sr models.RequestSeller
	if err := httputil.DecodeJSON(r, &sr); err != nil {
		response.Error(w, err)
		return
	}

	err := validators.ValidateSellerPost(sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	s, err := h.sv.Create(ctx, sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, s)
}

func (h *SellerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var sr models.RequestSeller
	if err := httputil.DecodeJSON(r, &sr); err != nil {
		response.Error(w, err)
		return
	}

	err = validators.ValidateSellerPatchNotEmpty(sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	err = validators.ValidateSellerPatch(sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	s, err := h.sv.Update(ctx, id, sr)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, s)
}

func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	err = h.sv.Delete(ctx, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *SellerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s, err := h.sv.FindAll(ctx)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, s)
}

func (h *SellerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	s, err := h.sv.FindById(ctx, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, s)
}

func handleApiError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if errorResp, ok := err.(*api.ServiceError); ok {
		// response.Error(w, errorResp.ResponseCode, errorResp.Message)
		// response.Error(w, errorResp.ResponseCode, errorResp.Message)
		response.Error(w, errorResp)
	} else {
		// response.Error(w, http.StatusInternalServerError, err.Error())
		// response.Error(w, http.StatusInternalServerError, err.Error())
		response.Error(w, err)
	}

	return true
}
