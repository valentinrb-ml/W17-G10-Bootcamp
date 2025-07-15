package carry

type CarryRequest struct {
	Id          int    `json:"id,omitempty"`
	Cid         string `json:"cid,omitempty"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address,omitempty"`
	Telephone   string `json:"telephone,omitempty"`
	LocalityId  int    `json:"locality_id"`
}

type CarryDoc struct {
	ID          int    `json:"id,omitempty"`
	Cid         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  int    `json:"locality_id"`
}
