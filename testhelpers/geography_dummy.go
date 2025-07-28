package testhelpers

import (
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

var CountriesDummyMap = map[int]models.Country{
	1: {Id: 1, Name: "Argentina"},
	2: {Id: 2, Name: "Brasil"},
	3: {Id: 3, Name: "Chile"},
	4: {Id: 4, Name: "Uruguay"},
	5: {Id: 5, Name: "Paraguay"},
}

var ProvincesDummyMap = map[int]models.Province{
	1: {Id: 1, Name: "Buenos Aires", CountryId: 1},
	2: {Id: 2, Name: "Córdoba", CountryId: 1},
	3: {Id: 3, Name: "Santa Fe", CountryId: 1},
	4: {Id: 4, Name: "Mendoza", CountryId: 1},
	5: {Id: 5, Name: "San Pablo", CountryId: 2},
}

var LocalitiesDummyMap = map[string]models.Locality{
	"1900":     {Id: "1900", Name: "La Plata", ProvinceId: 1},
	"5000":     {Id: "5000", Name: "Córdoba Capital", ProvinceId: 2},
	"2000":     {Id: "2000", Name: "Rosario", ProvinceId: 3},
	"5501":     {Id: "5501", Name: "Godoy Cruz", ProvinceId: 4},
	"13001970": {Id: "13001970", Name: "Campinas", ProvinceId: 5},
}

// CreateTestLocality creates a test locality
func CreateTestLocality(id string) *models.Locality {
	return &models.Locality{
		Id:         id,
		Name:       "Test Locality",
		ProvinceId: 1,
	}
}

func DummyRequestGeography() models.RequestGeography {
	id := "5194"
	country_name := "Argentina"
	province_name := "Cordoba"
	locality_name := "Villa General Belgrano"

	return models.RequestGeography{
		Id:           &id,
		CountryName:  &country_name,
		ProvinceName: &province_name,
		LocalityName: &locality_name,
	}
}

func DummyResponseGeography() models.ResponseGeography {
	return models.ResponseGeography{
		LocalityId:   "5194",
		CountryName:  "Argentina",
		ProvinceName: "Cordoba",
		LocalityName: "Villa General Belgrano",
	}
}

func DummyResponseLocalitySellers() models.ResponseLocalitySellers {
	return models.ResponseLocalitySellers{
		LocalityId:   "2000",
		LocalityName: "Rosario",
		SellersCount: 3,
	}
}

func DummySliceResponseLocalitySellers() []models.ResponseLocalitySellers {
	return []models.ResponseLocalitySellers{
		{LocalityId: "1900", LocalityName: "La Plata", SellersCount: 5},
		{LocalityId: "2000", LocalityName: "Rosario", SellersCount: 3},
	}
}
