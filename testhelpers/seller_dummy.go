package testhelpers

import (
	"sort"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

func DummyRequestSeller() models.RequestSeller {
	cid := 101
	companyName := "Frutas del Sur"
	address := "Calle 1"
	telephone := "221-111"
	localityId := "1900"

	return models.RequestSeller{
		Cid:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		LocalityId:  &localityId,
	}
}

func DummyResponseSeller() models.ResponseSeller {
	return models.ResponseSeller(SellersDummyMap[1])
}

var SellersDummyMap = map[int]models.Seller{
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

func FindAllSellersDummy() []models.Seller {
	out := make([]models.Seller, 0, len(SellersDummyMap))
	for i := 1; i <= len(SellersDummyMap); i++ {
		if seller, ok := SellersDummyMap[i]; ok {
			out = append(out, seller)
		}
	}

	return out
}

func FindAllSellersResponseDummy() []models.ResponseSeller {
	var keys []int
	for k := range SellersDummyMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var sellersList []models.Seller
	for _, k := range keys {
		sellersList = append(sellersList, SellersDummyMap[k])
	}

	ms := mappers.ToResponseSellerList(sellersList)
	return ms
}

// Ptr helper for expected pointer struct
func Ptr[T any](v T) *T { return &v }
