package handler

import (
	"net/http"
	"strconv"

	svsProductBatch "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product_batch"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_batches"
)

// ProductBatchesHandler handles HTTP requests for product batches endpoints.
type ProductBatchesHandler struct {
	sv     svsProductBatch.ProductBatchesService
	logger logger.Logger
}

// SetLogger allows injecting the logger after creation
func (h *ProductBatchesHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

// NewProductBatchesHandler creates a new ProductBatchesHandler with provided service.
func NewProductBatchesHandler(sv svsProductBatch.ProductBatchesService) *ProductBatchesHandler {
	return &ProductBatchesHandler{
		sv: sv,
	}
}

// CreateProductBatches handles POST requests to create a new product batch.
// - Decodes the JSON body, validates input, and calls service to persist.
// - Responds with proper error or the created product batch.
func (h *ProductBatchesHandler) CreateProductBatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if h.logger != nil {
		h.logger.Info(ctx, "product-batches-handler", "Create product batch request received")
	}

	var req models.PostProductBatches
	if err := httputil.DecodeJSON(r, &req); err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "product-batches-handler", "Failed to decode JSON", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	if err := validators.ValidateProductBatchPost(req); err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "product-batches-handler", "Validation failed for create request", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	newProBa, err := h.sv.CreateProductBatches(ctx, mappers.RequestToProductBatch(req))
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "product-batches-handler", "Failed to create product batch", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, mappers.ProductBatchesToResponse(*newProBa))
}

// GetReportProduct handles GET requests for product batch reports.
// - If 'id' query param is present, returns report for specific section; else returns all.
// - Returns 400 if 'id' is not a valid integer.
func (h *ProductBatchesHandler) GetReportProduct(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "product-batches-handler", "Get report product request received")
	}
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		report, err := h.sv.GetReportProduct(ctx)
		if err != nil {
			if h.logger != nil {
				h.logger.Error(ctx, "product-batches-handler", "Failed to get report product", err)
			}
			response.ErrorWithRequest(w, r, err)
			return
		}
		response.JSON(w, http.StatusOK, report)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "product-batches-handler", "Failed to parse id parameter", err)
		}
		response.ErrorWithRequest(w, r, apperrors.NewAppError(apperrors.CodeBadRequest, "invalid integer"))
		return
	}
	report, err := h.sv.GetReportProductById(ctx, idInt)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(ctx, "product-batches-handler", "Failed to get report product by id", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, report)
}
