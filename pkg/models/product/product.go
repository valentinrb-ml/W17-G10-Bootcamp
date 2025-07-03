package product

type Dimensions struct{ Width, Height, Length float64 }

type Expiration struct {
	Rate                    int
	RecommendedFreezingTemp float64
	FreezingRate            int
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
// productBase contains the fields shared by Request and Response
type ProductData struct {
	ProductCode                    string  `json:"product_code"`
	Description                    string  `json:"description"`
	Width                          float64 `json:"width"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 int     `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   int     `json:"freezing_rate"`
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
	ExpirationRate                 *int     `json:"expiration_rate,omitempty"`
	RecommendedFreezingTemperature *float64 `json:"recommended_freezing_temperature,omitempty"`
	FreezingRate                   *int     `json:"freezing_rate,omitempty"`
	ProductTypeID                  *int     `json:"product_type_id,omitempty"`
	SellerID                       *int     `json:"seller_id,omitempty"`
}
