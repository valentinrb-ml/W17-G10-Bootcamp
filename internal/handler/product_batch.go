package handler

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
	"net/http"
	"strconv"
)

type ProductBatchesHandler struct {
	sv service.ProductBatchesService
}

func NewProductBatchesHandler(sv service.ProductBatchesService) *ProductBatchesHandler {
	return &ProductBatchesHandler{
		sv,
	}
}

func (h *ProductBatchesHandler) CreateProductBatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.PostProductBatches
	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}
	if err := validators.ValidateProductBatchPost(req); err != nil {
		response.Error(w, err)
		return
	}
	newProBa, err := h.sv.CreateProductBatches(ctx, mappers.RequestToProductBatch(req))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, mappers.ProductBatchesToResponse(*newProBa))
}

func (h *ProductBatchesHandler) GetReportProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		report, err := h.sv.GetReportProduct(ctx)
		if err != nil {
			response.Error(w, err)
			return
		}
		response.JSON(w, http.StatusOK, report)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, apperrors.NewAppError(apperrors.CodeBadRequest, "id must be a valid integer"))
		return
	}
	report, err := h.sv.GetReportProductById(ctx, idInt)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, report)
	return
}
