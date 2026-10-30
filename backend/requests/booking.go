package requests

type BookingFilter struct {
	MarketID   int    `form:"market_id"`
	StallIdent string `form:"stall_ident"`
	Status     string `form:"status"`
	Search     string `form:"search"`
	Page       string `form:"page"`
	Limit      string `form:"limit"`
}

type CreateBookingRequest struct {
	// Stalls is the canonical multi-stall payload. Legacy single-stall fields
	// remain supported so older clients keep working during rollout.
	Stalls             []BookingStallRequest `json:"stalls"`
	StallID            interface{}           `json:"stall_id"`
	StallIdent         string                `json:"stall_ident"`
	StallNumber        string                `json:"stall_number"`
	MarketID           *int                  `json:"market_id"`
	MarketName         string                `json:"market_name"`
	ProjectName        string                `json:"project_name"`
	ApplicantName      string                `json:"applicant_name"`
	TenantName         string                `json:"tenant_name"`
	ApplicantPhone     string                `json:"applicant_phone"`
	TenantPhone        string                `json:"tenant_phone"`
	ApplicantEmail     string                `json:"applicant_email"`
	TenantEmail        string                `json:"tenant_email"`
	ApplicantIdCard    string                `json:"applicant_id_card"`
	TenantIdCard       string                `json:"tenant_id_card"`
	BusinessType       string                `json:"business_type"`
	IntendedBusiness   string                `json:"intended_business"`
	ProductType        string                `json:"product_type"`
	VisitDate          string                `json:"visit_date"`
	VisitTime          string                `json:"visit_time"`
	VisitTimeSlot      string                `json:"visit_time_slot"`
	Note               string                `json:"note"`
	IsQueue            bool                  `json:"is_queue"`
	AppointmentGroupID string                `json:"appointment_group_id"`
	MonthlyRent        float64               `json:"monthly_rent"`
	DepositFee         float64               `json:"deposit_fee"`
	SkipNotification   bool                  `json:"skip_notification"`
}

type BookingStallRequest struct {
	StallID     interface{} `json:"stall_id"`
	StallIdent  string      `json:"stall_ident"`
	StallNumber string      `json:"stall_number"`
	ProjectName string      `json:"project_name"`
	MonthlyRent float64     `json:"monthly_rent"`
	DepositFee  float64     `json:"deposit_fee"`
}

type CreateWaitlistRequest struct {
	StallID     interface{} `json:"stall_id"`
	StallIdent  string      `json:"stall_ident"`
	StallNumber string      `json:"stall_number"`
	MarketID    *int        `json:"market_id"`
	MarketName  string      `json:"market_name"`
	ProjectName string      `json:"project_name"`
	TenantName  string      `json:"tenant_name"`
	TenantPhone string      `json:"tenant_phone"`
}

type UpdateBookingStatusRequest struct {
	RentalType      string `json:"rental_type"`
	Status          string `json:"status" validate:"nonzero"`
	Reason          string `json:"reason"`
	RejectionReason string `json:"rejection_reason"`
	MeetingPoint    string `json:"meeting_point"`
	EscortOfficer   string `json:"escort_officer"`
	OfficerNotes    string `json:"officer_notes"`
	Note            string `json:"note"`
	ByOfficer       bool   `json:"by_officer"`
	StallIdent      string `json:"stall_ident"`
	CancelAll       bool   `json:"cancel_all"`
	ActorRole       string `json:"-"`
}

type UpdateBookingPolicySettingsRequest struct {
	MarketID                     int    `json:"market_id"`
	MarketName                   string `json:"market_name"`
	LeaseApplicationDeadlineDays int    `json:"lease_application_deadline_days"`
	WaitlistOfferHours           int    `json:"waitlist_offer_hours"`
}
