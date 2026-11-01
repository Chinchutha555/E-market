package responses

import (
	"encoding/json"
	"time"
)

type MarketUtilitySettingsResponse struct {
	MarketID  int             `json:"market_id"`
	Settings  json.RawMessage `json:"settings"`
	UpdatedBy *int            `json:"updated_by"`
	UpdatedAt *time.Time      `json:"updated_at"`
}
