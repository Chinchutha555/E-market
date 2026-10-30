package models

import "time"

type Booking struct {
	ID                 int        `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingNo          string     `gorm:"size:50;unique;not null;column:booking_no" json:"booking_no"`
	StallID            *int       `gorm:"column:stall_id" json:"stall_id"`
	StallIdent         string     `gorm:"size:50;not null;column:stall_ident" json:"stall_ident"`
	MarketID           *int       `gorm:"column:market_id" json:"market_id"`
	MarketName         string     `gorm:"size:255;column:market_name" json:"market_name"`
	ProjectName        string     `gorm:"size:255;column:project_name" json:"project_name"`
	MonthlyRent        float64    `gorm:"column:monthly_rent;default:1800" json:"monthly_rent"`
	DepositFee         float64    `gorm:"column:deposit_fee;default:1800" json:"deposit_fee"`
	ApplicantName      string     `gorm:"size:255;not null;column:applicant_name" json:"applicant_name"`
	ApplicantPhone     string     `gorm:"size:50;not null;column:applicant_phone" json:"applicant_phone"`
	ApplicantEmail     string     `gorm:"size:255;column:applicant_email" json:"applicant_email"`
	ApplicantIdCard    string     `gorm:"size:50;column:applicant_id_card" json:"applicant_id_card"`
	BusinessType       string     `gorm:"size:255;column:business_type" json:"business_type"`
	VisitDate          string     `gorm:"size:20;not null;column:visit_date" json:"visit_date"`
	VisitTime          string     `gorm:"size:20;not null;column:visit_time" json:"visit_time"`
	MeetingPoint       string     `gorm:"size:255;column:meeting_point" json:"meeting_point"`
	EscortOfficer      string     `gorm:"size:255;column:escort_officer" json:"escort_officer"`
	RentalType         string     `gorm:"size:30;column:rental_type" json:"rental_type"`
	OfficerNotes       string     `gorm:"type:text;column:officer_notes" json:"officer_notes"`
	RejectionReason    string     `gorm:"type:text;column:rejection_reason" json:"rejection_reason"`
	Note               string     `gorm:"type:text;column:note" json:"note"`
	BookingStatus      string     `gorm:"size:50;default:'pending';column:booking_status" json:"booking_status"`
	QueueNumber        int        `gorm:"default:0;column:queue_number" json:"queue_number"`
	AppointmentGroupID string     `gorm:"size:100;column:appointment_group_id" json:"appointment_group_id"`
	SkipNotification   bool       `gorm:"-" json:"-"`
	Timeline           string     `gorm:"type:text;column:timeline" json:"timeline"`
	ExpiresAt          *time.Time `gorm:"column:expires_at" json:"expires_at"`
	Status             string     `gorm:"type:char(1);default:'A';column:status" json:"status"`
	CreatedBy          *int       `gorm:"column:created_by" json:"created_by"`
	UpdatedBy          *int       `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Booking) TableName() string {
	return "booking"
}

type BookingStall struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	BookingID          uint       `gorm:"column:booking_id;index" json:"booking_id"`
	SubBookingID       uint       `gorm:"column:sub_booking_id;index" json:"sub_booking_id"`
	BookingNo          string     `gorm:"column:booking_no;index" json:"booking_no"`
	AppointmentGroupID string     `gorm:"column:appointment_group_id;index" json:"appointment_group_id"`
	MarketID           uint       `gorm:"column:market_id;default:1" json:"market_id"`
	MarketName         string     `gorm:"column:market_name" json:"market_name"`
	ProjectName        string     `gorm:"column:project_name" json:"project_name"`
	StallID            uint       `gorm:"column:stall_id" json:"stall_id"`
	StallIdent         string     `gorm:"column:stall_ident;index" json:"stall_ident"`
	MonthlyRent        float64    `gorm:"column:monthly_rent;default:1800" json:"monthly_rent"`
	DepositFee         float64    `gorm:"column:deposit_fee;default:1800" json:"deposit_fee"`
	Status             string     `gorm:"column:status;default:pending_approval" json:"status"`
	CancelReason       string     `gorm:"column:cancel_reason" json:"cancel_reason"`
	CancelledBy        string     `gorm:"column:cancelled_by" json:"cancelled_by"`
	CancelledAt        *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (BookingStall) TableName() string {
	return "booking_stall"
}
