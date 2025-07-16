package mappers

import "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"

func RequestToCarry(req carry.CarryRequest) carry.Carry {
	return carry.Carry{
		Id:          0,
		Cid:         req.Cid,
		CompanyName: req.CompanyName,
		Address:     req.Address,
		Telephone:   req.Telephone,
		LocalityId:  req.LocalityId,
	}
}

func CarryToDoc(w *carry.Carry) carry.CarryDoc {
	return carry.CarryDoc{
		ID:          w.Id,
		Cid:         w.Cid,
		CompanyName: w.CompanyName,
		Address:     w.Address,
		Telephone:   w.Telephone,
		LocalityId:  w.LocalityId,
	}
}

func CarryToDocSlice(c []carry.Carry) []carry.CarryDoc {
	newCarriers := make([]carry.CarryDoc, 0, len(c))
	for _, ca := range c {
		cDoc := CarryToDoc(&ca)
		newCarriers = append(newCarriers, cDoc)
	}
	return newCarriers
}
