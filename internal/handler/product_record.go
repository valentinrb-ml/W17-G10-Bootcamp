package handler

import (
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"
)

type ProductRecordHandler struct {
	svc service.ProductRecordService
}

func NewProductRecordHandler(svc service.ProductRecordService) *ProductRecordHandler {
	return &ProductRecordHandler{svc: svc}
}

func (h *ProductRecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.ProductRecordRequest

	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validators.ValidateProductRecordCreateRequest(req); err != nil {
		response.Error(w, err)
		return
	}

	record := mappers.ProductRecordRequestToDomain(req)

	result, err := h.svc.Create(r.Context(), record)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *ProductRecordHandler) GetRecordsReport(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter ‘id’ (optional)
	productID, err := httputil.ParseOptionalIntParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	// Validate only if the ID was provided
	if productID != 0 {
		if err = validators.ValidateID(productID, "id"); err != nil {
			response.Error(w, err)
			return
		}
	}

	report, err := h.svc.GetRecordsReport(r.Context(), productID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, report)
}
