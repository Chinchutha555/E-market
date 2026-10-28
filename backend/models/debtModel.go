package models

import "time"

type Debt struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	TenantID     int       `gorm:"not null;index" json:"tenant_id"`
	MarketID     int       `gorm:"not null;index" json:"market_id"`
	Description  string    `gorm:"size:500;not null" json:"description"`
	AmountSatang int64     `gorm:"not null" json:"amount_satang"`
	Status       string    `gorm:"size:20;not null;default:outstanding;index" json:"status"`
	CreatedBy    int       `gorm:"not null" json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Debt) TableName() string { return "tenant_debt" }
