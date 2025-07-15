package carry

type CarryRequest struct {
	Id          int    `json:"id,omitempty"`
	Cid         string `json:"cid,omitempty"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address,omitempty"`
	Telephone   string `json:"telephone,omitempty"`
	LocalityId  string    `json:"locality_id"`
}

type CarryDoc struct {
	ID          int    `json:"id,omitempty"`
	Cid         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  string    `json:"locality_id"`
}

type CarriesReport struct {
	LocalityID   string `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	CarriesCount int `json:"carries_count"`
}