package models

type Employee struct {
	ID           int
	CardNumberID string
	FirstName    string
	LastName     string
	WarehouseID  int
}

type EmployeeDoc struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}
type EmployeePatch struct {
	CardNumberID *string `json:"card_number_id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	WarehouseID  *int    `json:"warehouse_id,omitempty"`
}
