package requests

type ContractFilter struct {
	MarketID   int    `form:"market_id"`
	CreatedBy  int    `form:"created_by"`
	StallIdent string `form:"stall_ident"`
	Status     string `form:"status"`
	Search     string `form:"search"`
	Page       string `form:"page"`
	Limit      string `form:"limit"`
	ViewerRole string `form:"-"`
}

type SaveContractDraftRequest struct {
	TenantEmail        string  `json:"tenant_email"`
	TenantAddress      string  `json:"tenant_address"`
	BusinessType       string  `json:"business_type"`
	BillingCycle       string  `json:"billing_cycle"`
	TermsAndConditions string  `json:"terms_and_conditions"`
	TenantName         string  `json:"tenant_name"`
	TenantPhone        string  `json:"tenant_phone"`
	TenantIDCard       string  `json:"tenant_id_card"`
	ShopName           string  `json:"shop_name"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	MonthlyRental      float64 `json:"monthly_rental"`
	SecurityDeposit    float64 `json:"security_deposit"`
	AdvancePayment     float64 `json:"advance_payment"`
}

type CreateContractRequest struct {
	TenantEmail        string  `json:"tenant_email"`
	BillingCycle       string  `json:"billing_cycle"`
	TermsAndConditions string  `json:"terms_and_conditions"`
	SignatureBase64    string  `json:"signature_base64"`
	DocumentBase64     string  `json:"document_base64"`
	ApplicationID      int     `json:"application_id"`
	StallIdent         string  `json:"stall_ident"`
	StallNumber        string  `json:"stall_number"`
	StallID            int     `json:"stall_id"`
	TenantName         string  `json:"tenant_name"`
	TenantPhone        string  `json:"tenant_phone"`
	TenantIDCard       string  `json:"tenant_id_card"`
	TenantAddress      string  `json:"tenant_address"`
	ShopName           string  `json:"shop_name"`
	BusinessName       string  `json:"business_name"`
	BusinessType       string  `json:"business_type"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	MonthlyRental      float64 `json:"monthly_rental"`
	RentalFee          float64 `json:"rental_fee"`
	SecurityDeposit    float64 `json:"security_deposit"`
	AdvancePayment     float64 `json:"advance_payment"`
	Status             string  `json:"status"`
}

type ReviewContractRequest struct {
	Reason string `json:"reason"`
}

// UpdateContractWorkflowRequest อัปเดตสถานะติดตามสัญญา (ส่งเฉพาะฟิลด์ที่ต้องการเปลี่ยน; ส่งสตริงว่างเพื่อล้างสถานะ)
type UpdateContractWorkflowRequest struct {
	RenewalStatus *string `json:"renewal_status"`
	ReturnStatus  *string `json:"return_status"`
	LegalStatus   *string `json:"legal_status"`
}
