package product

import (
	"database/sql"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

// Request DTO -> Domain
func ToDomain(r models.ProductRequest) models.Product {
	return models.Product{
		Code:        r.ProductCode,
		Description: r.Description,
		Dimensions:  models.Dimensions{Width: r.Width, Height: r.Height, Length: r.Length},
		NetWeight:   r.NetWeight,
		Expiration: models.Expiration{
			Rate:                    r.ExpirationRate,
			RecommendedFreezingTemp: r.RecommendedFreezingTemperature,
			FreezingRate:            r.FreezingRate,
		},
		ProductType: r.ProductTypeID,
		SellerID:    r.SellerID,
	}
}

// Domain -> Response DTO.
func FromDomain(p models.Product) models.ProductResponse {
	return models.ProductResponse{
		ID: p.ID,
		ProductData: models.ProductData{
			ProductCode:                    p.Code,
			Description:                    p.Description,
			Width:                          p.Dimensions.Width,
			Height:                         p.Dimensions.Height,
			Length:                         p.Dimensions.Length,
			NetWeight:                      p.NetWeight,
			ExpirationRate:                 p.Expiration.Rate,
			RecommendedFreezingTemperature: p.Expiration.RecommendedFreezingTemp,
			FreezingRate:                   p.Expiration.FreezingRate,
			ProductTypeID:                  p.ProductType,
			SellerID:                       p.SellerID,
		},
	}
}

// Helper for lists
func FromDomainList(list []models.Product) []models.ProductResponse {
	out := make([]models.ProductResponse, 0, len(list))
	for _, p := range list {
		out = append(out, FromDomain(p))
	}
	return out
}

// Response DTO -> Domain (used by the Loader)
func ResponseToDomain(r models.ProductResponse) models.Product {
	return models.Product{
		ID:          r.ID,
		Code:        r.ProductCode,
		Description: r.Description,
		Dimensions:  models.Dimensions{Width: r.Width, Height: r.Height, Length: r.Length},
		NetWeight:   r.NetWeight,
		Expiration: models.Expiration{
			Rate:                    r.ExpirationRate,
			RecommendedFreezingTemp: r.RecommendedFreezingTemperature,
			FreezingRate:            r.FreezingRate,
		},
		ProductType: r.ProductTypeID,
		SellerID:    r.SellerID,
	}
}

// DB -> Domain
func DbToDomain(d models.ProductDb) models.Product {
	var seller *int
	if d.SellerID.Valid {
		v := int(d.SellerID.Int64)
		seller = &v
	}

	return models.Product{
		ID:          d.ID,
		Code:        d.Code,
		Description: d.Description,
		Dimensions:  models.Dimensions{Width: d.Width, Height: d.Height, Length: d.Length},
		NetWeight:   d.NetWeight,
		Expiration: models.Expiration{
			Rate:                    d.ExpRate,
			RecommendedFreezingTemp: d.RecFreeze,
			FreezingRate:            d.FreezeRate,
		},
		ProductType: d.TypeID,
		SellerID:    seller,
	}
}

// Domain -> DB
func FromDomainToDb(p models.Product) models.ProductDb {
	var seller sql.NullInt64
	if p.SellerID != nil {
		seller = sql.NullInt64{Int64: int64(*p.SellerID), Valid: true}
	}

	return models.ProductDb{
		ID:          p.ID,
		Code:        p.Code,
		Description: p.Description,
		Width:       p.Dimensions.Width,
		Height:      p.Dimensions.Height,
		Length:      p.Dimensions.Length,
		NetWeight:   p.NetWeight,
		ExpRate:     p.Expiration.Rate,
		RecFreeze:   p.Expiration.RecommendedFreezingTemp,
		FreezeRate:  p.Expiration.FreezingRate,
		TypeID:      p.ProductType,
		SellerID:    seller,
	}
}
