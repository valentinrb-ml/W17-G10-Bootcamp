package mappers

import models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"

func RequestSellerToSeller(rs models.RequestSeller) models.Seller {
	return models.Seller{
		Id:          0,
		Cid:         *rs.Cid,
		CompanyName: *rs.CompanyName,
		Address:     *rs.Address,
		Telephone:   *rs.Telephone,
		LocalityId:  *rs.LocalityId,
	}
}

func ToResponseSeller(s *models.Seller) models.ResponseSeller {
	return models.ResponseSeller{
		Id:          s.Id,
		Cid:         s.Cid,
		CompanyName: s.CompanyName,
		Address:     s.Address,
		Telephone:   s.Telephone,
	}
}

func ToResponseSellerList(ss []models.Seller) []models.ResponseSeller {
	res := make([]models.ResponseSeller, len(ss))
	for i, s := range ss {
		res[i] = ToResponseSeller(&s)
	}
	return res
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
	if patch.LocalityId != nil {
		seller.LocalityId = *patch.LocalityId
	}
}
