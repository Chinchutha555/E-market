package responses

import "time"

type TimelineItem struct {
	Action    string `json:"action"`
	ActorRole string `json:"actor_role"`
	ActorName string `json:"actor_name"`
	Timestamp string `json:"timestamp"`
	Note      string `json:"note,omitempty"`
}

type BookingResponse struct {
	RentalType                  string            `json:"rental_type,omitempty"`
	ID                          int               `json:"id"`
	BookingNo                   string            `json:"booking_no"`
	BookingNumber               string            `json:"booking_number"`
	StallID                     *int              `json:"stall_id,omitempty"`
	StallIdent                  string            `json:"stall_ident"`
	StallNumber                 string            `json:"stall_number"`
	MarketID                    *int              `json:"market_id,omitempty"`
	MarketName                  string            `json:"market_name"`
	ProjectName                 string            `json:"project_name"`
	ApplicantName               string            `json:"applicant_name"`
	TenantName                  string            `json:"tenant_name"`
	TenantUserID                *int              `json:"tenant_user_id,omitempty"`
	ApplicantPhone              string            `json:"applicant_phone"`
	TenantPhone                 string            `json:"tenant_phone"`
	ApplicantEmail              string            `json:"applicant_email"`
	TenantEmail                 string            `json:"tenant_email"`
	TenantIdCard                string            `json:"tenant_id_card,omitempty"`
	BusinessType                string            `json:"business_type"`
	IntendedBusiness            string            `json:"intended_business"`
	VisitDate                   string            `json:"visit_date"`
	VisitTime                   string            `json:"visit_time"`
	VisitTimeSlot               string            `json:"visit_time_slot"`
	MeetingPoint                string            `json:"meeting_point,omitempty"`
	EscortOfficer               string            `json:"escort_officer,omitempty"`
	OfficerNotes                string            `json:"officer_notes,omitempty"`
	RejectionReason             string            `json:"rejection_reason,omitempty"`
	CancelReason                string            `json:"cancel_reason,omitempty"`
	NotInterestedReason         string            `json:"not_interested_reason,omitempty"`
	Note                        string            `json:"note"`
	Status                      string            `json:"status"`
	ApplicationNumber           string            `json:"application_number,omitempty"`
	QueueNumber                 int               `json:"queue_number"`
	QueuePosition               int               `json:"queue_position"`
	HoldingStall                bool              `json:"holding_stall"`
	AppointmentGroupID          string            `json:"appointment_group_id,omitempty"`
	Timeline                    []TimelineItem    `json:"timeline"`
	ExpiresAt                   *time.Time        `json:"expires_at,omitempty"`
	ApplicationDeadlineDays     int               `json:"application_deadline_days,omitempty"`
	RemainingApplicationSeconds int64             `json:"remaining_application_seconds,omitempty"`
	CreatedAt                   time.Time         `json:"created_at"`
	SubBookings                 []BookingResponse `json:"sub_bookings,omitempty"`
	StallCount                  int               `json:"stall_count,omitempty"`
	ActiveStallCount            int               `json:"active_stall_count"`
	TotalStallCount             int               `json:"total_stall_count"`
	StallMonthlyRate            float64           `json:"stall_monthly_rate"`
	DepositFee                  float64           `json:"deposit_fee"`
	TotalMonthlyRent            float64           `json:"total_monthly_rent"`
	TotalDeposit                float64           `json:"total_deposit"`
	WaitlistCount               int               `json:"waitlist_count"`
}

type BookingPolicySettingsResponse struct {
	MarketID                     int    `json:"market_id"`
	MarketName                   string `json:"market_name"`
	LeaseApplicationDeadlineDays int    `json:"lease_application_deadline_days"`
	WaitlistOfferHours           int    `json:"waitlist_offer_hours"`
}

type WaitlistResponse struct {
	ID                    int        `json:"id"`
	StallID               *int       `json:"stall_id,omitempty"`
	StallIdent            string     `json:"stall_ident"`
	StallNumber           string     `json:"stall_number"`
	MarketID              *int       `json:"market_id,omitempty"`
	MarketName            string     `json:"market_name"`
	ProjectName           string     `json:"project_name"`
	TenantID              int        `json:"tenant_id"`
	TenantName            string     `json:"tenant_name"`
	TenantPhone           string     `json:"tenant_phone"`
	Position              int        `json:"position"`
	Status                string     `json:"status"`
	DismissReason         string     `json:"dismiss_reason,omitempty"`
	TenantDismissed       bool       `json:"tenant_dismissed"`
	OfferedAt             *time.Time `json:"offered_at,omitempty"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"`
	OfferExpiresAt        *time.Time `json:"offer_expires_at,omitempty"`
	RemainingSeconds      int64      `json:"remaining_seconds"`
	OfferRemainingSeconds int64      `json:"offer_remaining_seconds"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
