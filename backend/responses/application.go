package responses

import "time"

type ApplicationResponse struct {
	ID                    int        `json:"id"`
	ApplicationNo         string     `json:"application_no"`
	ApplicationNumber     string     `json:"application_number"`
	ApplicationType       string     `json:"application_type"`
	Subject               string     `json:"subject,omitempty"`
	BookingID             *int       `json:"booking_id,omitempty"`
	StallID               *int       `json:"stall_id,omitempty"`
	MarketID              *int       `json:"market_id,omitempty"`
	MarketName            string     `json:"market_name,omitempty"`
	StallIdent            string     `json:"stall_ident"`
	StallNumber           string     `json:"stall_number"`
	ApplicantName         string     `json:"applicant_name"`
	TenantName            string     `json:"tenant_name"`
	ApplicantPhone        string     `json:"applicant_phone"`
	TenantPhone           string     `json:"tenant_phone"`
	ApplicantIDCard       string     `json:"applicant_id_card"`
	ApplicantAge          *int       `json:"applicant_age,omitempty"`
	HasIDCardCopy         bool       `json:"has_id_card_copy"`
	ApplicantEmail        string     `json:"applicant_email,omitempty"`
	ApplicantAddress      string     `json:"applicant_address,omitempty"`
	Details               string     `json:"details,omitempty"`
	ShopName              string     `json:"shop_name"`
	ProductCategory       string     `json:"product_category"`
	ProductType           string     `json:"product_type"`
	ProposedRentType      string     `json:"proposed_rent_type,omitempty"`
	RequestedStartDate    string     `json:"requested_start_date"`
	LeaseTermMonths       int        `json:"lease_term_months"`
	RequestedPeriodMonths int        `json:"requested_period_months"`
	MonthlyRate           float64    `json:"monthly_rate"`
	StallMonthlyRate      float64    `json:"stall_monthly_rate"`
	Status                string     `json:"status"`
	OfficerComment        string     `json:"officer_comment,omitempty"`
	MarketDirectorComment string     `json:"market_director_comment,omitempty"`
	BureauDirectorComment string     `json:"bureau_director_comment,omitempty"`
	ReviewedAt            *time.Time `json:"reviewed_at,omitempty"`
	ReviewNote            string     `json:"review_note,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}
