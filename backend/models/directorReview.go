package models

type MarketDirectorReview struct {
	DebtStatus string  `json:"debt_status"`
	DebtAmount float64 `json:"debt_amount"`
}

type BureauDirectorReview struct {
	Decision      string `json:"decision"`
	OtherDetails  string `json:"other_details"`
	OfficeChecked bool   `json:"office_checked"`
	SystemNumber  string `json:"system_number"`
	DocumentDate  string `json:"document_date"`
}
