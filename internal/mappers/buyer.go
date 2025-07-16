package mappers

import (
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// Convierte un request en un buyer DB
func RequestBuyerToBuyer(req models.RequestBuyer) models.Buyer {
	return models.Buyer{
		CardNumberId: coalesceString(req.CardNumberId),
		FirstName:    coalesceString(req.FirstName),
		LastName:     coalesceString(req.LastName),
	}
}

func ApplyBuyerPatch(b *models.Buyer, req *models.RequestBuyer) {
	if req.CardNumberId != nil {
		b.CardNumberId = *req.CardNumberId
	}
	if req.FirstName != nil {
		b.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		b.LastName = *req.LastName
	}
}

// Convierte de DB a respuesta para la API
func ToResponseBuyer(b *models.Buyer) models.ResponseBuyer {
	return models.ResponseBuyer{
		Id:           b.Id,
		CardNumberId: b.CardNumberId,
		FirstName:    b.FirstName,
		LastName:     b.LastName,
	}
}

// Convierte una lista de buyers a response para API
func ToResponseBuyerList(buyers []models.Buyer) []models.ResponseBuyer {
	res := make([]models.ResponseBuyer, 0, len(buyers))
	for _, b := range buyers {
		res = append(res, ToResponseBuyer(&b))
	}
	return res
}

// para campos opcionales en request
func coalesceString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
