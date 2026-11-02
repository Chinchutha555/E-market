package models

import "time"

type Contract struct {
	ContractTemplateID          *int    `gorm:"column:contract_template_id" json:"contract_template_id"`
	ContractTemplateCode        string  `gorm:"size:40;column:contract_template_code" json:"contract_template_code"`
	ContractTemplateTitle       string  `gorm:"size:255;column:contract_template_title" json:"contract_template_title"`
	ContractTemplateKind        string  `gorm:"size:30;column:contract_template_kind" json:"contract_template_kind"`
	ContractTemplateDescription string  `gorm:"type:text;column:contract_template_description" json:"contract_template_description"`
	TenantEmail                 string  `gorm:"column:tenant_email" json:"-"`
	TenantAddress               string  `gorm:"column:tenant_address" json:"-"`
	BusinessType                string  `gorm:"column:business_type" json:"-"`
	BillingCycle                string  `gorm:"column:billing_cycle" json:"-"`
	TermsAndConditions          string  `gorm:"column:terms_and_conditions" json:"-"`
	AdvancePayment              float64 `gorm:"column:advance_payment" json:"-"`
	TenantSignature             []byte  `gorm:"column:tenant_signature" json:"-"`
	TenantSignatureType         string  `gorm:"column:tenant_signature_type" json:"-"`
	OfficerSignature            []byte  `gorm:"column:officer_signature" json:"-"`
	OfficerSignatureType        string  `gorm:"column:officer_signature_type" json:"-"`
	MarketDirectorSignature     []byte  `gorm:"column:market_director_signature" json:"-"`
	MarketDirectorSignatureType string  `gorm:"column:market_director_signature_type" json:"-"`
	BureauDirectorSignature     []byte  `gorm:"column:bureau_director_signature" json:"-"`
	BureauDirectorSignatureType string  `gorm:"column:bureau_director_signature_type" json:"-"`
	SupportingDocument          []byte  `gorm:"column:supporting_document" json:"-"`
	SupportingDocumentType      string  `gorm:"column:supporting_document_type" json:"-"`
	ReviewReason                string  `gorm:"column:review_reason" json:"-"`
	ID                          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	ContractNumber              string  `gorm:"size:50;unique;not null;column:contract_number" json:"contract_number"`
	ApplicationID               *int    `gorm:"column:application_id" json:"application_id"`
	StallID                     *int    `gorm:"column:stall_id" json:"stall_id"`
	StallIdent                  string  `gorm:"size:50;not null;column:stall_ident" json:"stall_ident"`
	MarketID                    *int    `gorm:"column:market_id;index" json:"market_id"`
	TenantUserID                *int    `gorm:"column:tenant_user_id;index" json:"tenant_user_id"`
	TenantName                  string  `gorm:"size:255;not null;column:tenant_name" json:"tenant_name"`
	TenantPhone                 string  `gorm:"size:50;not null;column:tenant_phone" json:"tenant_phone"`
	TenantIDCard                string  `gorm:"size:50;column:tenant_id_card" json:"tenant_id_card"`
	ShopName                    string  `gorm:"size:255;not null;column:shop_name" json:"shop_name"`
	StartDate                   string  `gorm:"size:20;not null;column:start_date" json:"start_date"`
	EndDate                     string  `gorm:"size:20;not null;column:end_date" json:"end_date"`
	MonthlyRental               float64 `gorm:"column:monthly_rental" json:"monthly_rental"`
	SecurityDeposit             float64 `gorm:"column:security_deposit" json:"security_deposit"`
	ContractState               string  `gorm:"size:50;default:'active';column:contract_state" json:"contract_state"`
	// สถานะติดตามสัญญาที่เจ้าหน้าที่บันทึก (ว่าง = ยังไม่มีเรื่องที่ต้องติดตาม)
	RenewalStatus string    `gorm:"size:30;default:'';column:renewal_status" json:"renewal_status"`
	ReturnStatus  string    `gorm:"size:30;default:'';column:return_status" json:"return_status"`
	LegalStatus   string    `gorm:"size:30;default:'';column:legal_status" json:"legal_status"`
	Status        string    `gorm:"type:char(1);default:'A';column:status" json:"status"`
	CreatedBy     *int      `gorm:"column:created_by" json:"created_by"`
	UpdatedBy     *int      `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Contract) TableName() string {
	return "contract"
}
