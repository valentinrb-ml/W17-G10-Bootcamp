package models

type Country struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Province struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CountryId int    `json:"country_id"`
}

type Locality struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ProvinceId int    `json:"province_id"`
}

type RequestGeography struct {
	Id           *string `json:"id"`
	LocalityName *string `json:"locality_name"`
	ProvinceName *string `json:"province_name"`
	CountryName  *string `json:"country_name"`
}

type ResponseGeography struct {
	LocalityId   string `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
}
