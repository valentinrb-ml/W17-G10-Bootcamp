package validators

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductValidator struct {
	err *apperrors.AppError
}

func NewProductValidator() *ProductValidator {
	return &ProductValidator{}
}

func (v *ProductValidator) addError(field, message string) {
	if v.err == nil {
		v.err = apperrors.NewAppError(apperrors.CodeValidationError, "validation failed")
	}
	v.err = v.err.WithDetail(field, message)
}

func (v *ProductValidator) Error() error {
	return v.err
}

// Reusable validation functions
func (v *ProductValidator) validateProductCode(code string, required bool) {
	if required && code == "" {
		v.addError("product_code", "is required")
		return
	}
	if code != "" && len(code) > 50 {
		v.addError("product_code", "cannot exceed 50 characters")
	}
}

func (v *ProductValidator) validateDescription(desc string, required bool) {
	if required && desc == "" {
		v.addError("description", "is required")
		return
	}
	if desc != "" && len(desc) > 500 {
		v.addError("description", "cannot exceed 500 characters")
	}
}

func (v *ProductValidator) validateWidth(width float64, required bool) {
	if required && width == 0 {
		v.addError("width", "is required")
		return
	}
	if width < 0 {
		v.addError("width", "cannot be negative")
	}
}

func (v *ProductValidator) validateHeight(height float64, required bool) {
	if required && height == 0 {
		v.addError("height", "is required")
		return
	}
	if height < 0 {
		v.addError("height", "cannot be negative")
	}
}

func (v *ProductValidator) validateLength(length float64, required bool) {
	if required && length == 0 {
		v.addError("length", "is required")
		return
	}
	if length < 0 {
		v.addError("length", "cannot be negative")
	}
}

func (v *ProductValidator) validateNetWeight(weight float64, required bool) {
	if required && weight == 0 {
		v.addError("net_weight", "is required")
		return
	}
	if weight < 0 {
		v.addError("net_weight", "cannot be negative")
	}
}

// Validators
func ValidateCreateRequest(req product.ProductRequest) error {
	validator := NewProductValidator()

	validator.validateProductCode(req.ProductCode, true)
	validator.validateDescription(req.Description, true)
	validator.validateWidth(req.Width, true)
	validator.validateHeight(req.Height, true)
	validator.validateLength(req.Length, true)
	validator.validateNetWeight(req.NetWeight, true)

	return validator.Error()
}

func ValidatePatchRequest(req product.ProductPatchRequest) error {
	if !hasAnyPatchField(req) {
		return apperrors.BadRequest("at least one field must be provided for update")
	}

	validator := NewProductValidator()

	if req.ProductCode != nil {
		validator.validateProductCode(*req.ProductCode, false)
	}
	if req.Description != nil {
		validator.validateDescription(*req.Description, false)
	}
	if req.Width != nil {
		validator.validateWidth(*req.Width, false)
	}
	if req.Height != nil {
		validator.validateHeight(*req.Height, false)
	}
	if req.Length != nil {
		validator.validateLength(*req.Length, false)
	}
	if req.NetWeight != nil {
		validator.validateNetWeight(*req.NetWeight, false)
	}

	return validator.Error()
}

func ValidateProductBusinessRules(p product.Product) error {
	validator := NewProductValidator()

	if p.Dimensions.Width <= 0 || p.Dimensions.Height <= 0 || p.Dimensions.Length <= 0 {
		validator.addError("dimensions", "must be positive values")
	}

	if p.NetWeight <= 0 {
		validator.addError("net_weight", "must be positive")
	}

	return validator.Error()
}

func hasAnyPatchField(req product.ProductPatchRequest) bool {
	return req.ProductCode != nil || req.Description != nil || req.Width != nil ||
		req.Height != nil || req.Length != nil || req.NetWeight != nil ||
		req.ExpirationRate != nil || req.FreezingRate != nil ||
		req.RecommendedFreezingTemperature != nil || req.ProductTypeID != nil ||
		req.SellerID != nil
}
