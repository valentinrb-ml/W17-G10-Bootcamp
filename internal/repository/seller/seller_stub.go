package repository

import models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"

var SellersMapStub = map[int]models.Seller{
	1: {
		Id:          1,
		Cid:         101,
		CompanyName: "Frutas del Sur",
		Address:     "Calle 1",
		Telephone:   "221-111",
		LocalityId:  "1900",
	},
	2: {
		Id:          2,
		Cid:         102,
		CompanyName: "Verdulería Norte",
		Address:     "Calle 2",
		Telephone:   "221-112",
		LocalityId:  "5000",
	},
	3: {
		Id:          3,
		Cid:         103,
		CompanyName: "Carnes Argentinas",
		Address:     "Calle 3",
		Telephone:   "221-113",
		LocalityId:  "2000",
	},
	4: {
		Id:          4,
		Cid:         104,
		CompanyName: "Almacén Cordobés",
		Address:     "Calle 4",
		Telephone:   "221-114",
		LocalityId:  "5501",
	},
	5: {
		Id:          5,
		Cid:         105,
		CompanyName: "Exportadora Brasil",
		Address:     "Calle 5",
		Telephone:   "11-221",
		LocalityId:  "13001970",
	},
}

func FindAllSellersStub() []models.Seller {
	out := make([]models.Seller, 0, len(SellersMapStub))
	for i := 1; i <= len(SellersMapStub); i++ {
		if seller, ok := SellersMapStub[i]; ok {
			out = append(out, seller)
		}
	}

	return out
}
