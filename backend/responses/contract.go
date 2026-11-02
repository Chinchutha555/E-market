package responses

import "time"

type ContractResponse struct {
	TenantEmail         string    `json:"tenant_email"`
	TenantAddress       string    `json:"tenant_address"`
	BusinessType        string    `json:"business_type"`
	BillingCycle        string    `json:"billing_cycle"`
	TermsAndConditions  string    `json:"terms_and_conditions"`
	AdvancePayment      float64   `json:"advance_payment"`
	HasTenantSignature  bool      `json:"has_tenant_signature"`
	HasOfficerSignature bool      `json:"has_officer_signature"`
	HasMarketSignature  bool      `json:"has_market_signature"`
	HasBureauSignature  bool      `json:"has_bureau_signature"`
	HasDocument         bool      `json:"has_document"`
	ReviewReason        string    `json:"review_reason"`
	ID                  int       `json:"id"`
	ContractNumber      string    `json:"contract_number"`
	ApplicationID       *int      `json:"application_id,omitempty"`
	StallID             *int      `json:"stall_id,omitempty"`
	MarketID            *int      `json:"market_id,omitempty"`
	TenantUserID        *int      `json:"tenant_user_id,omitempty"`
	StallIdent          string    `json:"stall_ident"`
	TenantName          string    `json:"tenant_name"`
	TenantPhone         string    `json:"tenant_phone"`
	TenantIDCard        string    `json:"tenant_id_card,omitempty"`
	ShopName            string    `json:"shop_name"`
	StartDate           string    `json:"start_date"`
	EndDate             string    `json:"end_date"`
	MonthlyRental       float64   `json:"monthly_rental"`
	SecurityDeposit     float64   `json:"security_deposit"`
	Status              string    `json:"status"`
	RenewalStatus       string    `json:"renewal_status"`
	ReturnStatus        string    `json:"return_status"`
	LegalStatus         string    `json:"legal_status"`
	CreatedAt           time.Time `json:"created_at"`
}
