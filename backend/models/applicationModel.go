package models

import "time"

type Application struct {
	ID                    int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationNo         string     `gorm:"size:50;unique;not null;column:application_no" json:"application_no"`
	ApplicationType       string     `gorm:"size:50;default:'stall_lease';column:application_type" json:"application_type"`
	Subject               string     `gorm:"size:255;column:subject" json:"subject"`
	BookingID             *int       `gorm:"column:booking_id;index" json:"booking_id"`
	StallID               *int       `gorm:"column:stall_id" json:"stall_id"`
	StallIdent            string     `gorm:"size:50;not null;column:stall_ident" json:"stall_ident"`
	MarketID              *int       `gorm:"column:market_id;index" json:"market_id"`
	MarketName            string     `gorm:"size:255;column:market_name" json:"market_name"`
	ApplicantName         string     `gorm:"size:255;not null;column:applicant_name" json:"applicant_name"`
	ApplicantPhone        string     `gorm:"size:50;not null;column:applicant_phone" json:"applicant_phone"`
	ApplicantIDCard       string     `gorm:"size:50;column:applicant_id_card" json:"applicant_id_card"`
	ApplicantAge          *int       `gorm:"column:applicant_age" json:"applicant_age,omitempty"`
	IDCardCopy            []byte     `gorm:"column:id_card_copy" json:"-"`
	IDCardCopyType        string     `gorm:"size:50;column:id_card_copy_type" json:"-"`
	ShopName              string     `gorm:"size:255;not null;column:shop_name" json:"shop_name"`
	ProductCategory       string     `gorm:"size:255;column:product_category" json:"product_category"`
	LeaseTermMonths       int        `gorm:"default:12;column:lease_term_months" json:"lease_term_months"`
	MonthlyRate           float64    `gorm:"column:monthly_rate" json:"monthly_rate"`
	ProposedRentType      string     `gorm:"size:50;default:'monthly';column:proposed_rent_type" json:"proposed_rent_type"`
	RequestedStartDate    string     `gorm:"size:50;column:requested_start_date" json:"requested_start_date"`
	ApplicantEmail        string     `gorm:"size:255;column:applicant_email" json:"applicant_email"`
	ApplicantAddress      string     `gorm:"type:text;column:applicant_address" json:"applicant_address"`
	Details               string     `gorm:"type:text;column:details" json:"details"`
	OfficerComment        string     `gorm:"type:text;column:officer_comment" json:"officer_comment"`
	MarketDirectorComment string     `gorm:"type:text;column:market_director_comment" json:"market_director_comment"`
	BureauDirectorComment string     `gorm:"type:text;column:bureau_director_comment" json:"bureau_director_comment"`
	ApplicationState      string     `gorm:"size:50;default:'submitted';column:application_state" json:"application_state"`
	ReviewedAt            *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
	ReviewNote            string     `gorm:"type:text;column:review_note" json:"review_note"`
	Status                string     `gorm:"type:char(1);default:'A';column:status" json:"status"`
	CreatedBy             *int       `gorm:"column:created_by" json:"created_by"`
	UpdatedBy             *int       `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt             time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Application) TableName() string {
	return "application"
}
