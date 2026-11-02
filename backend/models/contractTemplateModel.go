package models

import "time"

type MarketContractTemplate struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	MarketID     int       `gorm:"column:market_id;not null;uniqueIndex:idx_market_contract_template_code" json:"market_id"`
	Code         string    `gorm:"size:40;not null;uniqueIndex:idx_market_contract_template_code" json:"code"`
	Title        string    `gorm:"size:255;not null" json:"title"`
	ContractKind string    `gorm:"size:30;not null" json:"contract_kind"`
	Description  string    `gorm:"type:text" json:"description"`
	IsActive     bool      `gorm:"column:is_active;default:true;not null" json:"is_active"`
	SortOrder    int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (MarketContractTemplate) TableName() string { return "market_contract_templates" }
