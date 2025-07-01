package mappers

import models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"

func RequestSellerToSeller(rs models.RequestSeller) models.Seller {
	return models.Seller{
		Id:          0,
		Cid:         *rs.Cid,
		CompanyName: *rs.CompanyName,
		Address:     *rs.Address,
		Telephone:   *rs.Telephone,
	}
}
