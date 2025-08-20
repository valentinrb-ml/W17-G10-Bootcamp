package handler

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	"net/http"

	productMappers "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers/product"
	productService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductHandler struct {
	svc    productService.ProductService
	logger logger.Logger
}

// SetLogger allows injecting the logger after creation
func (h *ProductHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

func NewProductHandler(s productService.ProductService) *ProductHandler {
	return &ProductHandler{svc: s}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "GetAll products request received") // [LOG]
	}

	list, err := h.svc.GetAll(r.Context())
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "product-handler", "Failed to list products", err) // [LOG]
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Products listed successfully", map[string]interface{}{ // [LOG]
			"count": len(list), // [LOG]
		})
	}

	response.JSON(w, http.StatusOK, list)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Create product request received") // [LOG]
	}

	var req models.ProductRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Invalid JSON in create request", map[string]interface{}{ // [LOG]
				"error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if err := validators.ValidateCreateRequest(req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Validation failed for create request", map[string]interface{}{ // [LOG]
				"product_code":     req.ProductCode, // [LOG]
				"validation_error": err.Error(),     // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	newProduct := productMappers.ToDomain(req)

	result, err := h.svc.Create(r.Context(), newProduct)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "product-handler", "Failed to create product", err, map[string]interface{}{ // [LOG]
				"product_code": req.ProductCode, // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Product created successfully", map[string]interface{}{ // [LOG]
			"product_code": req.ProductCode, // [LOG]
			"product_id":   result.ID,       // [LOG]
		})
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Invalid product ID in get-by-id request", map[string]interface{}{ // [LOG]
				"param": "id",        // [LOG]
				"error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "GetByID product request received", map[string]interface{}{ // [LOG]
			"product_id": id, // [LOG]
		})
	}

	if err = validators.ValidateID(id, "product id"); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Validation failed for product id", map[string]interface{}{ // [LOG]
				"product_id":       id,          // [LOG]
				"validation_error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	currentProduct, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "product-handler", "Failed to get product by id", err, map[string]interface{}{ // [LOG]
				"product_id": id, // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Product retrieved successfully", map[string]interface{}{ // [LOG]
			"product_id": id, // [LOG]
		})
	}

	response.JSON(w, http.StatusOK, currentProduct)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Invalid product ID in delete request", map[string]interface{}{ // [LOG]
				"param": "id",        // [LOG]
				"error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Delete product request received", map[string]interface{}{ // [LOG]
			"product_id": id, // [LOG]
		})
	}

	if err = validators.ValidateID(id, "product id"); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Validation failed for product id in delete", map[string]interface{}{ // [LOG]
				"product_id":       id,          // [LOG]
				"validation_error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	err = h.svc.Delete(r.Context(), id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "product-handler", "Failed to delete product", err, map[string]interface{}{ // [LOG]
				"product_id": id, // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Product deleted successfully", map[string]interface{}{ // [LOG]
			"product_id": id, // [LOG]
		})
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *ProductHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseIntParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Invalid product ID in patch request", map[string]interface{}{ // [LOG]
				"param": "id",        // [LOG]
				"error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Patch product request received", map[string]interface{}{ // [LOG]
			"product_id": id, // [LOG]
		})
	}

	if err = validators.ValidateID(id, "product id"); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Validation failed for product id in patch", map[string]interface{}{ // [LOG]
				"product_id":       id,          // [LOG]
				"validation_error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	var req models.ProductPatchRequest
	if err = httputil.DecodeJSON(r, &req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Invalid JSON in patch product request", map[string]interface{}{ // [LOG]
				"product_id": id,          // [LOG]
				"error":      err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if err = validators.ValidatePatchRequest(req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "product-handler", "Validation failed for patch product request", map[string]interface{}{ // [LOG]
				"product_id":       id,          // [LOG]
				"validation_error": err.Error(), // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	result, err := h.svc.Patch(r.Context(), id, req)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "product-handler", "Failed to patch product", err, map[string]interface{}{ // [LOG]
				"product_id": id, // [LOG]
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "product-handler", "Product patched successfully", map[string]interface{}{ // [LOG]
			"product_id": id, // [LOG]
		})
	}

	response.JSON(w, http.StatusOK, result)
}
