package models

type OfficerReview struct {
	StallStatus       string `json:"stall_status"`
	StallCount        int    `json:"stall_count"`
	StallNumbers      string `json:"stall_numbers"`
	Project           string `json:"project"`
	TradeType         string `json:"trade_type"`
	OtherContract     bool   `json:"other_contract"`
	OtherChecked      bool   `json:"other_checked"`
	OtherDetails      string `json:"other_details"`
	DocumentsComplete bool   `json:"documents_complete"`
	BookingReference  string `json:"booking_reference"`
}
