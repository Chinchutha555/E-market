package requests

type ApplicationFilter struct {
	MarketID        int    `form:"market_id"`
	CreatedBy       int    `form:"created_by"`
	StallIdent      string `form:"stall_ident"`
	Status          string `form:"status"`
	ApplicationType string `form:"application_type"`
	Search          string `form:"search"`
	Page            string `form:"page"`
	Limit           string `form:"limit"`
}

type CreateApplicationRequest struct {
	ApplicationType       string      `json:"application_type"`
	Subject               string      `json:"subject"`
	BookingID             interface{} `json:"booking_id"`
	StallID               interface{} `json:"stall_id"`
	StallIdent            string      `json:"stall_ident"`
	StallNumber           string      `json:"stall_number"`
	MarketID              interface{} `json:"market_id"`
	MarketName            string      `json:"market_name"`
	ApplicantName         string      `json:"applicant_name"`
	TenantName            string      `json:"tenant_name"`
	ApplicantPhone        string      `json:"applicant_phone"`
	TenantPhone           string      `json:"tenant_phone"`
	ApplicantIDCard       string      `json:"applicant_id_card"`
	ApplicantAge          *int        `json:"applicant_age"`
	IDCardCopyBase64      string      `json:"id_card_copy_base64"`
	IDCard                string      `json:"id_card"`
	ApplicantEmail        string      `json:"applicant_email"`
	ApplicantAddress      string      `json:"applicant_address"`
	Address               string      `json:"address"`
	ShopName              string      `json:"shop_name"`
	ProductCategory       string      `json:"product_category"`
	ProductType           string      `json:"product_type"`
	RequestedStartDate    string      `json:"requested_start_date"`
	LeaseTermMonths       int         `json:"lease_term_months"`
	RequestedPeriodMonths int         `json:"requested_period_months"`
	MonthlyRate           float64     `json:"monthly_rate"`
	Fee                   float64     `json:"fee"`
	ProposedRentType      string      `json:"proposed_rent_type"`
	ReviewNote            string      `json:"review_note"`
	Details               string      `json:"details"`
}

type UpdateApplicationStatusRequest struct {
	OfficerSignature      []byte `json:"-"`
	OfficerSignatureType  string `json:"-"`
	Status                string `json:"status" validate:"nonzero"`
	Reason                string `json:"reason"`
	OfficerComment        string `json:"officer_comment"`
	MarketDirectorComment string `json:"market_director_comment"`
	BureauDirectorComment string `json:"bureau_director_comment"`
}

type ApplicationActionRequest struct {
	SignatureBase64 string `json:"signature_base64"`
	Action          string `json:"action"` // "approve" or "reject"
	Status          string `json:"status"`
	Comment         string `json:"comment"`
	Reason          string `json:"reason"`
}
