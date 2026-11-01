package models

import "time"

// MarketUtilitySetting เก็บการตั้งค่าอัตราค่าน้ำ-ไฟและรหัสมิเตอร์แผงของตลาด
// SettingsJSON เป็นเอกสาร JSON เดียวต่อตลาด: scope, defaults, projects, stalls, meters
type MarketUtilitySetting struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	MarketID     int       `gorm:"column:market_id;uniqueIndex;not null" json:"market_id"`
	SettingsJSON string    `gorm:"type:text;column:settings_json" json:"-"`
	UpdatedBy    *int      `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (MarketUtilitySetting) TableName() string { return "market_utility_settings" }
