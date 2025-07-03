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

func ApplySellerPatch(seller *models.Seller, patch *models.RequestSeller) {
	if patch.Cid != nil {
		seller.Cid = *patch.Cid
	}
	if patch.CompanyName != nil {
		seller.CompanyName = *patch.CompanyName
	}
	if patch.Address != nil {
		seller.Address = *patch.Address
	}
	if patch.Telephone != nil {
		seller.Telephone = *patch.Telephone
	}
}
