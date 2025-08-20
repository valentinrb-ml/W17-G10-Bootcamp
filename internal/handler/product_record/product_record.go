package handler

import (
	productRecordMappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product_record"
	productRecordService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product_record"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

type ProductRecordHandler struct {
	svc    productRecordService.ProductRecordService
	logger logger.Logger // [LOG]
}

func NewProductRecordHandler(svc productRecordService.ProductRecordService) *ProductRecordHandler {
	return &ProductRecordHandler{svc: svc}
}

// SetLogger allows injecting the logger after creation
func (h *ProductRecordHandler) SetLogger(l logger.Logger) { // [LOG]
	h.logger = l // [LOG]
}

func (h *ProductRecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "product-record-handler", "Create product record request received") // [LOG]
	}

	var req models.ProductRecordRequest

	if err := httputil.DecodeJSON(r, &req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-record-handler", "Invalid JSON in create product record request", map[string]interface{}{ // [LOG]
				"error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if err := validators.ValidateProductRecordCreateRequest(req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-record-handler", "Validation failed for product record create request", map[string]interface{}{ // [LOG]
				"product_id":       req.Data.ProductID, // [LOG]
				"validation_error": err.Error(),        // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	record := productRecordMappers.ProductRecordRequestToDomain(req)

	result, err := h.svc.Create(r.Context(), record)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "product-record-handler", "Failed to create product record", err, map[string]interface{}{ // [LOG]
				"product_id": req.Data.ProductID, // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-record-handler", "Product record created successfully", map[string]interface{}{ // [LOG]
			"product_id":        result.ProductID, // [LOG]
			"product_record_id": result.ID,        // [LOG]
		})
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *ProductRecordHandler) GetRecordsReport(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "product-record-handler", "Get product records report request received") // [LOG]
	}

	// Parse query parameter ‘id’ (optional)
	productID, err := httputil.ParseOptionalIntParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-record-handler", "Invalid product_id query parameter for report", map[string]interface{}{ // [LOG]
				"error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	// Validate only if the ID was provided
	if productID != 0 {
		if err = validators.ValidateID(productID, "id"); err != nil {
			if h.logger != nil {
				h.logger.Warning(r.Context(), "product-record-handler", "Validation failed for report product_id", map[string]interface{}{ // [LOG]
					"product_id":       productID,   // [LOG]
					"validation_error": err.Error(), // [LOG]
				})
			}
			response.ErrorWithRequest(w, r, err)
			return
		}
	}

	report, err := h.svc.GetRecordsReport(r.Context(), productID)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "product-record-handler", "Failed to get product records report", err, map[string]interface{}{ // [LOG]
				"product_id": productID, // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		if productID != 0 {
			h.logger.Info(r.Context(), "product-record-handler", "Product records report generated (filtered)", map[string]interface{}{ // [LOG]
				"product_id": productID,   // [LOG]
				"count":      len(report), // [LOG]
			})
		} else {
			h.logger.Info(r.Context(), "product-record-handler", "Product records report generated", map[string]interface{}{ // [LOG]
				"count": len(report), // [LOG]
			})
		}
	}

	response.JSON(w, http.StatusOK, report)
}
