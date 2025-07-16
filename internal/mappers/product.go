package mappers

import (
	"database/sql"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

// Request DTO -> Domain
func ToDomain(r product.ProductRequest) product.Product {
	return product.Product{
		Code:        r.ProductCode,
		Description: r.Description,
		Dimensions:  product.Dimensions{Width: r.Width, Height: r.Height, Length: r.Length},
		NetWeight:   r.NetWeight,
		Expiration: product.Expiration{
			Rate:                    r.ExpirationRate,
			RecommendedFreezingTemp: r.RecommendedFreezingTemperature,
			FreezingRate:            r.FreezingRate,
		},
		ProductType: r.ProductTypeID,
		SellerID:    r.SellerID,
	}
}

// Domain -> Response DTO.
func FromDomain(p product.Product) product.ProductResponse {
	return product.ProductResponse{
		ID: p.ID,
		ProductData: product.ProductData{
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
func FromDomainList(list []product.Product) []product.ProductResponse {
	out := make([]product.ProductResponse, 0, len(list))
	for _, p := range list {
		out = append(out, FromDomain(p))
	}
	return out
}

// Response DTO -> Domain (used by the Loader)
func ResponseToDomain(r product.ProductResponse) product.Product {
	return product.Product{
		ID:          r.ID,
		Code:        r.ProductCode,
		Description: r.Description,
		Dimensions:  product.Dimensions{Width: r.Width, Height: r.Height, Length: r.Length},
		NetWeight:   r.NetWeight,
		Expiration: product.Expiration{
			Rate:                    r.ExpirationRate,
			RecommendedFreezingTemp: r.RecommendedFreezingTemperature,
			FreezingRate:            r.FreezingRate,
		},
		ProductType: r.ProductTypeID,
		SellerID:    r.SellerID,
	}
}

// DB -> Domain
func DbToDomain(d product.ProductDb) product.Product {
	var seller *int
	if d.SellerID.Valid {
		v := int(d.SellerID.Int64)
		seller = &v
	}

	return product.Product{
		ID:          d.ID,
		Code:        d.Code,
		Description: d.Description,
		Dimensions:  product.Dimensions{Width: d.Width, Height: d.Height, Length: d.Length},
		NetWeight:   d.NetWeight,
		Expiration: product.Expiration{
			Rate:                    d.ExpRate,
			RecommendedFreezingTemp: d.RecFreeze,
			FreezingRate:            d.FreezeRate,
		},
		ProductType: d.TypeID,
		SellerID:    seller,
	}
}

// Domain -> DB
func FromDomainToDb(p product.Product) product.ProductDb {
	var seller sql.NullInt64
	if p.SellerID != nil {
		seller = sql.NullInt64{Int64: int64(*p.SellerID), Valid: true}
	}

	return product.ProductDb{
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
