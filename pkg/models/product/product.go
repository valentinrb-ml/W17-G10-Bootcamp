package models

import "database/sql"

type Dimensions struct{ Width, Height, Length float64 }

type Expiration struct {
	Rate                    float64
	RecommendedFreezingTemp float64
	FreezingRate            float64
}

// Product Main model used in the Domain
type Product struct {
	ID          int
	Code        string
	Description string
	Dimensions  Dimensions
	NetWeight   float64
	Expiration  Expiration
	ProductType int
	SellerID    *int
}

// DTOs
// ProductData contains the fields shared by Request and Response
type ProductData struct {
	ProductCode                    string  `json:"product_code"`
	Description                    string  `json:"description"`
	Width                          float64 `json:"width"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 float64 `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   float64 `json:"freezing_rate"`
	ProductTypeID                  int     `json:"product_type_id"`
	SellerID                       *int    `json:"seller_id,omitempty"`
}

type ProductRequest struct {
	ProductData
}

type ProductResponse struct {
	ID int `json:"id"`
	ProductData
}

type ProductPatchRequest struct {
	ProductCode                    *string  `json:"product_code,omitempty"`
	Description                    *string  `json:"description,omitempty"`
	Width                          *float64 `json:"width,omitempty"`
	Height                         *float64 `json:"height,omitempty"`
	Length                         *float64 `json:"length,omitempty"`
	NetWeight                      *float64 `json:"net_weight,omitempty"`
	ExpirationRate                 *float64 `json:"expiration_rate,omitempty"`
	RecommendedFreezingTemperature *float64 `json:"recommended_freezing_temperature,omitempty"`
	FreezingRate                   *float64 `json:"freezing_rate,omitempty"`
	ProductTypeID                  *int     `json:"product_type_id,omitempty"`
	SellerID                       *int     `json:"seller_id,omitempty"`
}

// Structure that matches 1-to-1 with table `products`
type ProductDb struct {
	ID          int           `db:"id"`
	Code        string        `db:"product_code"`
	Description string        `db:"description"`
	Width       float64       `db:"width"`
	Height      float64       `db:"height"`
	Length      float64       `db:"length"`
	NetWeight   float64       `db:"net_weight"`
	ExpRate     float64       `db:"expiration_rate"`
	RecFreeze   float64       `db:"recommended_freezing_temperature"`
	FreezeRate  float64       `db:"freezing_rate"`
	TypeID      int           `db:"product_type_id"`
	SellerID    sql.NullInt64 `db:"seller_id"`
}
