package usecases

import (
	"encoding/json"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"emarket-api/models"
	"emarket-api/responses"
)

func encodeMarketDirectorReview(input *models.MarketDirectorReview) (*string, error) {
	if input == nil || (input.DebtStatus != "none" && input.DebtStatus != "has") {
		return nil, responses.NewAppErr(422, "กรุณาเลือกข้อมูลหนี้ค้าง")
	}
	review := *input
	if review.DebtStatus == "none" {
		review.DebtAmount = 0
	}
	if review.DebtStatus == "has" && (math.IsNaN(review.DebtAmount) || math.IsInf(review.DebtAmount, 0) || review.DebtAmount <= 0 || review.DebtAmount > 1000000000000) {
		return nil, responses.NewAppErr(422, "กรุณาระบุจำนวนหนี้ค้าง")
	}
	encoded, err := json.Marshal(review)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	value := string(encoded)
	return &value, nil
}

func encodeBureauDirectorReview(input *models.BureauDirectorReview) (*string, error) {
	if input == nil || (input.Decision != "approve" && input.Decision != "reject") {
		return nil, responses.NewAppErr(422, "กรุณาเลือกผลการพิจารณา")
	}
	review := *input
	review.OtherDetails = strings.TrimSpace(review.OtherDetails)
	review.SystemNumber = strings.TrimSpace(review.SystemNumber)
	if review.Decision == "other" && review.OtherDetails == "" {
		return nil, responses.NewAppErr(422, "กรุณาระบุรายละเอียดอื่น ๆ")
	}
	if review.Decision != "other" {
		review.OtherDetails = ""
	}
	if !review.OfficeChecked {
		review.SystemNumber = ""
		review.DocumentDate = ""
	}
	if utf8.RuneCountInString(review.OtherDetails) > 255 || utf8.RuneCountInString(review.SystemNumber) > 255 {
		return nil, responses.NewAppErr(422, "รายละเอียดแต่ละช่องต้องไม่เกิน 255 ตัวอักษร")
	}
	if review.DocumentDate != "" {
		if _, err := time.Parse("2006-01-02", review.DocumentDate); err != nil {
			return nil, responses.NewAppErr(422, "วันที่ไม่ถูกต้อง")
		}
	}
	encoded, err := json.Marshal(review)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	value := string(encoded)
	return &value, nil
}
