package usecases

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	"emarket-api/models"
	"emarket-api/responses"
)

func encodeOfficerReview(input *models.OfficerReview) (*string, error) {
	if input == nil || (input.StallStatus != "none" && input.StallStatus != "has" && input.StallStatus != "other") {
		return nil, responses.NewAppErr(422, "กรุณาเลือกข้อมูลแผงค้า")
	}
	review := *input
	review.StallNumbers = strings.TrimSpace(review.StallNumbers)
	review.Project = strings.TrimSpace(review.Project)
	review.TradeType = strings.TrimSpace(review.TradeType)
	review.OtherDetails = strings.TrimSpace(review.OtherDetails)
	review.BookingReference = strings.TrimSpace(review.BookingReference)
	if review.StallStatus == "has" && (review.StallCount < 1 || review.StallCount > 100000) {
		return nil, responses.NewAppErr(422, "กรุณาระบุจำนวนแผงค้า")
	}
	if review.StallStatus == "none" || review.StallStatus == "other" {
		review.StallCount = 0
		review.StallNumbers = ""
		review.Project = ""
		review.TradeType = ""
	}
	if review.StallStatus == "other" {
		review.OtherChecked = true
	}
	if !review.OtherChecked {
		review.OtherDetails = ""
	}
	if !review.DocumentsComplete {
		review.BookingReference = ""
	}
	if review.OtherChecked && review.OtherDetails == "" {
		return nil, responses.NewAppErr(422, "กรุณาระบุรายละเอียดอื่น ๆ")
	}
	for _, value := range []string{review.StallNumbers, review.Project, review.TradeType, review.OtherDetails, review.BookingReference} {
		if utf8.RuneCountInString(value) > 255 {
			return nil, responses.NewAppErr(422, "รายละเอียดแต่ละช่องต้องไม่เกิน 255 ตัวอักษร")
		}
	}
	encoded, err := json.Marshal(review)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	text := string(encoded)
	return &text, nil
}
