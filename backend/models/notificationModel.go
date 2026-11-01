package models

import "time"

// Notification is persisted because read/unread state and ownership must be
// consistent across API instances and deployments.
type Notification struct {
	SourceNotificationID *int64    `gorm:"column:source_notification_id;uniqueIndex:idx_notification_source_recipient" json:"-"`
	ID                   int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipientRole        string    `gorm:"size:30;not null;index;column:recipient_role" json:"recipient_role"`
	RecipientUserID      *int      `gorm:"index;column:recipient_user_id;uniqueIndex:idx_notification_source_recipient" json:"recipient_user_id,omitempty"`
	RecipientName        string    `gorm:"size:255;column:recipient_name" json:"recipient_name"`
	Type                 string    `gorm:"size:50;not null;column:type" json:"type"`
	StatusCode           string    `gorm:"size:50;column:status_code" json:"status_code"`
	StatusLabel          string    `gorm:"size:255;column:status_label" json:"status_label"`
	Title                string    `gorm:"size:500;not null;column:title" json:"title"`
	Message              string    `gorm:"type:text;column:message" json:"message"`
	BookingID            *int      `gorm:"index;column:booking_id" json:"booking_id,omitempty"`
	BookingNumber        string    `gorm:"size:50;column:booking_number" json:"booking_number"`
	StallNumber          string    `gorm:"size:50;column:stall_number" json:"stall_number"`
	MarketName           string    `gorm:"size:255;column:market_name" json:"market_name"`
	ProjectName          string    `gorm:"size:255;column:project_name" json:"project_name"`
	TargetURL            string    `gorm:"size:500;column:target_url" json:"target_url"`
	IsRead               bool      `gorm:"default:false;index;column:is_read" json:"is_read"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
	UpdatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Notification) TableName() string {
	return "notification"
}
