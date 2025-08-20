package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

func NewWarehouseHandler(sv service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{sv: sv}
}

type WarehouseHandler struct {
	sv     service.WarehouseService
	logger logger.Logger
}

// SetLogger allows injecting the logger after creation
func (h *WarehouseHandler) SetLogger(l logger.Logger) {
	h.logger = l
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "warehouse-handler", "Create warehouse request received")
	}

	var req = warehouse.WarehouseRequest{}
	if err := request.JSON(r, &req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "warehouse-handler", "Invalid JSON in create request", map[string]interface{}{
				"error": err.Error(),
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if err := validators.ValidateWarehouseCreateRequest(req); err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "warehouse-handler", "Validation failed for create request", map[string]interface{}{
				"warehouse_code":   req.WarehouseCode,
				"validation_error": err.Error(),
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	wh := mappers.RequestToWarehouse(req)

	newW, err := h.sv.Create(r.Context(), wh)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "warehouse-handler", "Failed to create warehouse", err, map[string]interface{}{
				"warehouse_code": req.WarehouseCode,
			})
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "warehouse-handler", "Warehouse created successfully", map[string]interface{}{
			"warehouse_id":   newW.Id,
			"warehouse_code": newW.WarehouseCode,
		})
	}

	// Get request ID from context to include in response
	requestID, _ := r.Context().Value("request_id").(string)

	responseData := mappers.WarehouseToDoc(newW)
	if requestID != "" {
		// You could add request_id to the response if needed
		// responseData["request_id"] = requestID
	}

	response.JSON(w, http.StatusCreated, responseData)
}

func (h *WarehouseHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	whs, err := h.sv.FindAll(r.Context())
	if err != nil {
		response.ErrorWithRequest(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDocSlice(whs))
}

func (h *WarehouseHandler) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid id")
		response.ErrorWithRequest(w, r, err)
		return
	}

	wh, er := h.sv.FindById(r.Context(), id)
	if er != nil {
		response.ErrorWithRequest(w, r, er)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDoc(wh))

}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid id")
		response.ErrorWithRequest(w, r, err)
		return
	}

	var req warehouse.WarehousePatchDTO
	if err := request.JSON(r, &req); err != nil {
		response.ErrorWithRequest(w, r, err)
		return
	}

	updated, serviceErr := h.sv.Update(r.Context(), id, req)
	if serviceErr != nil {
		response.ErrorWithRequest(w, r, serviceErr)
		return
	}

	response.JSON(w, http.StatusOK, mappers.WarehouseToDoc(updated))
}

func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if h.logger != nil {
			h.logger.Warning(r.Context(), "warehouse-handler", "Invalid warehouse ID in delete request", map[string]interface{}{
				"id_param": idStr,
			})
		}
		err := apperrors.NewAppError(apperrors.CodeBadRequest, "Invalid id")
		response.ErrorWithRequest(w, r, err)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "warehouse-handler", "Delete warehouse request received", map[string]interface{}{
			"warehouse_id": id,
		})
	}

	serviceErr := h.sv.Delete(r.Context(), id)
	if serviceErr != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "warehouse-handler", "Failed to delete warehouse", serviceErr, map[string]interface{}{
				"warehouse_id": id,
			})
		}
		response.ErrorWithRequest(w, r, serviceErr)
		return
	}

	if h.logger != nil {
		h.logger.Info(r.Context(), "warehouse-handler", "Warehouse deleted successfully", map[string]interface{}{
			"warehouse_id": id,
		})
	}

	response.JSON(w, http.StatusNoContent, nil)
}
