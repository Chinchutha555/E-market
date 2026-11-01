package models

import "time"

// MaintenanceRecord stores the full maintenance/request payload in the backend.
// Frequently used scope fields remain indexed columns; the JSON payload keeps the
// UI workflow extensible without requiring a migration for every optional field.
type MaintenanceRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MarketID    int       `gorm:"not null;index:idx_maintenance_market_stall,priority:1" json:"market_id"`
	StallIdent  string    `gorm:"size:100;not null;index:idx_maintenance_market_stall,priority:2" json:"stall_ident"`
	Status      string    `gorm:"size:30;not null;default:pending;index" json:"status"`
	PayloadJSON string    `gorm:"type:text;not null" json:"-"`
	CreatedBy   int       `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MaintenanceRecord) TableName() string { return "maintenance_record" }
