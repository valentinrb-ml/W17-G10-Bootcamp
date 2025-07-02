package mappers

import models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"

func RequestBuyerToBuyer(rb models.RequestBuyer) models.Buyer {
	b := models.Buyer{
		Id: 0,
	}

	if rb.CardNumberId != nil {
		b.CardNumberId = *rb.CardNumberId
	}
	if rb.FirstName != nil {
		b.FirstName = *rb.FirstName
	}
	if rb.LastName != nil {
		b.LastName = *rb.LastName
	}

	return b
}

func RequestBuyerToBuyerUpadate(rb models.RequestBuyer, existing models.Buyer) models.Buyer {

	updated := existing

	if rb.CardNumberId != nil {
		updated.CardNumberId = *rb.CardNumberId
	}
	if rb.FirstName != nil {
		updated.FirstName = *rb.FirstName
	}
	if rb.LastName != nil {
		updated.LastName = *rb.LastName
	}

	return updated
}
