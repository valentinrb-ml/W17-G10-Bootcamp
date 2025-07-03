package validators

import (
	"errors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

func ValidateCreateRequest(req product.ProductRequest) error {
	if req.ProductCode == "" {
		return errors.New("product_code is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.NetWeight == 0 {
		return errors.New("net_weight is required and must be provided")
	}
	if req.Width == 0 {
		return errors.New("width is required and must be provided")
	}
	if req.Height == 0 {
		return errors.New("height is required and must be provided")
	}
	if req.Length == 0 {
		return errors.New("length is required and must be provided")
	}
	if len(req.ProductCode) > 50 {
		return errors.New("product_code cannot exceed 50 characters")
	}
	if len(req.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}

	if req.Width < 0 || req.Height < 0 || req.Length < 0 {
		return errors.New("dimensions cannot be negative")
	}
	if req.NetWeight < 0 {
		return errors.New("net_weight cannot be negative")
	}

	return nil
}

func ValidateProductBusinessRules(p product.Product) error {
	if p.Dimensions.Width <= 0 || p.Dimensions.Height <= 0 || p.Dimensions.Length <= 0 {
		return serviceErr(api.ErrUnprocessableEntity, nil, "product dimensions must be positive values")
	}

	if p.NetWeight <= 0 {
		return serviceErr(api.ErrUnprocessableEntity, nil, "product net weight must be positive")
	}

	if p.NetWeight > 1000 {
		return serviceErr(api.ErrUnprocessableEntity, nil, "product exceeds maximum weight limit of 1000kg")
	}

	if p.Dimensions.Width < 0.1 || p.Dimensions.Height < 0.1 || p.Dimensions.Length < 0.1 {
		return serviceErr(api.ErrUnprocessableEntity, nil, "product dimensions must be at least 0.1cm")
	}

	maxDimension := 500.0
	if p.Dimensions.Width > maxDimension || p.Dimensions.Height > maxDimension || p.Dimensions.Length > maxDimension {
		return serviceErr(api.ErrUnprocessableEntity, nil, "product dimensions cannot exceed 500cm")
	}

	return nil
}

func ValidatePatchRequest(req product.ProductPatchRequest) error {
	if req.ProductCode == nil && req.Description == nil && req.Width == nil &&
		req.Height == nil && req.Length == nil && req.NetWeight == nil &&
		req.ExpirationRate == nil && req.FreezingRate == nil &&
		req.RecommendedFreezingTemperature == nil && req.ProductTypeID == nil &&
		req.SellerID == nil {
		return errors.New("at least one field must be provided for update")
	}
	return nil
}

func serviceErr(code int, internal error, overrideMsg string) error {
	e := api.ServiceErrors[code]
	if overrideMsg != "" {
		e.Message = overrideMsg
	}
	e.InternalError = internal
	return e
}
