package requests

type ContractFilter struct {
	MarketID   int    `form:"market_id"`
	CreatedBy  int    `form:"created_by"`
	StallIdent string `form:"stall_ident"`
	Status     string `form:"status"`
	Search     string `form:"search"`
	Page       string `form:"page"`
	Limit      string `form:"limit"`
}

type CreateContractRequest struct {
	ApplicationID   int     `json:"application_id"`
	StallIdent      string  `json:"stall_ident"`
	StallNumber     string  `json:"stall_number"`
	StallID         int     `json:"stall_id"`
	TenantName      string  `json:"tenant_name"`
	TenantPhone     string  `json:"tenant_phone"`
	TenantIDCard    string  `json:"tenant_id_card"`
	TenantAddress   string  `json:"tenant_address"`
	ShopName        string  `json:"shop_name"`
	BusinessName    string  `json:"business_name"`
	BusinessType    string  `json:"business_type"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	MonthlyRental   float64 `json:"monthly_rental"`
	RentalFee       float64 `json:"rental_fee"`
	SecurityDeposit float64 `json:"security_deposit"`
	AdvancePayment  float64 `json:"advance_payment"`
	Status          string  `json:"status"`
}

type ReviewContractRequest struct {
	Reason string `json:"reason"`
}
