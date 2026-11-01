package requests

import "encoding/json"

type UpdateMarketUtilitySettingsRequest struct {
	Settings json.RawMessage `json:"settings" binding:"required"`
}
