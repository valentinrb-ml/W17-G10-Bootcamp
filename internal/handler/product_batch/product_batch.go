package handler

import (
	svsProductBatch "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product_batch"
	"net/http"
	"strconv"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

// ProductBatchesHandler handles HTTP requests for product batches endpoints.
type ProductBatchesHandler struct {
	sv svsProductBatch.ProductBatchesService
}

// NewProductBatchesHandler creates a new ProductBatchesHandler with provided service.
func NewProductBatchesHandler(sv svsProductBatch.ProductBatchesService) *ProductBatchesHandler {
	return &ProductBatchesHandler{
		sv,
	}
}

// CreateProductBatches handles POST requests to create a new product batch.
// - Decodes the JSON body, validates input, and calls service to persist.
// - Responds with proper error or the created product batch.
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

// GetReportProduct handles GET requests for product batch reports.
// - If 'id' query param is present, returns report for specific section; else returns all.
// - Returns 400 if 'id' is not a valid integer.
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
}
