package responses

import (
	"emarket-api/models"
	"time"
)

type ApplicationResponse struct {
	DirectorReviewSupported    bool                         `json:"director_review_supported"`
	MarketDirectorReview       *models.MarketDirectorReview `json:"market_director_review,omitempty"`
	BureauDirectorReview       *models.BureauDirectorReview `json:"bureau_director_review,omitempty"`
	OfficerReviewSupported     bool                         `json:"officer_review_supported"`
	OfficerReview              *models.OfficerReview        `json:"officer_review,omitempty"`
	TradeOptionID              string                       `json:"trade_option_id"`
	TradeStartDate             string                       `json:"trade_start_date"`
	TradeEndDate               string                       `json:"trade_end_date"`
	TradeDuration              int                          `json:"trade_duration"`
	TradeDurationUnit          string                       `json:"trade_duration_unit"`
	RequiresContract           *bool                        `json:"requires_contract"`
	HasSignedDocument          bool                         `json:"has_signed_document"`
	HasOfficerSignature        bool                         `json:"has_officer_signature"`
	HasFinanceSignature        bool                         `json:"has_finance_signature"`
	HasMarketDirectorSignature bool                         `json:"has_market_director_signature"`
	HasBureauDirectorSignature bool                         `json:"has_bureau_director_signature"`
	DocumentNumber             string                       `json:"document_number,omitempty"`
	TenantSignature            string                       `json:"tenant_signature,omitempty"`
	TenantSignedAt             *time.Time                   `json:"tenant_signed_at,omitempty"`
	BookingNumber              string                       `json:"booking_number,omitempty"`
	ID                         int                          `json:"id"`
	ApplicationNo              string                       `json:"application_no"`
	ApplicationNumber          string                       `json:"application_number"`
	ApplicationType            string                       `json:"application_type"`
	Subject                    string                       `json:"subject,omitempty"`
	BookingID                  *int                         `json:"booking_id,omitempty"`
	StallID                    *int                         `json:"stall_id,omitempty"`
	MarketID                   *int                         `json:"market_id,omitempty"`
	MarketName                 string                       `json:"market_name,omitempty"`
	StallIdent                 string                       `json:"stall_ident"`
	StallNumber                string                       `json:"stall_number"`
	ApplicantName              string                       `json:"applicant_name"`
	TenantName                 string                       `json:"tenant_name"`
	ApplicantPhone             string                       `json:"applicant_phone"`
	TenantPhone                string                       `json:"tenant_phone"`
	ApplicantIDCard            string                       `json:"applicant_id_card"`
	ApplicantAge               *int                         `json:"applicant_age,omitempty"`
	HasIDCardCopy              bool                         `json:"has_id_card_copy"`
	ApplicantEmail             string                       `json:"applicant_email,omitempty"`
	ApplicantAddress           string                       `json:"applicant_address,omitempty"`
	Details                    string                       `json:"details,omitempty"`
	ShopName                   string                       `json:"shop_name"`
	ProductCategory            string                       `json:"product_category"`
	ProductType                string                       `json:"product_type"`
	ProposedRentType           string                       `json:"proposed_rent_type,omitempty"`
	RequestedStartDate         string                       `json:"requested_start_date"`
	LeaseTermMonths            int                          `json:"lease_term_months"`
	RequestedPeriodMonths      int                          `json:"requested_period_months"`
	MonthlyRate                float64                      `json:"monthly_rate"`
	StallMonthlyRate           float64                      `json:"stall_monthly_rate"`
	Status                     string                       `json:"status"`
	OfficerComment             string                       `json:"officer_comment,omitempty"`
	MarketDirectorComment      string                       `json:"market_director_comment,omitempty"`
	BureauDirectorComment      string                       `json:"bureau_director_comment,omitempty"`
	ReviewedAt                 *time.Time                   `json:"reviewed_at,omitempty"`
	ReviewNote                 string                       `json:"review_note,omitempty"`
	CreatedAt                  time.Time                    `json:"created_at"`
}
