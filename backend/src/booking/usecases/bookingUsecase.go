package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"emarket-api/logs"
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/booking/domains"
	"gorm.io/gorm"
)

type useCase struct {
	repo domains.Repository
}

func NewUseCase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func isTerminalBookingStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "rejected", "cancelled", "not_interested", "no_show", "terminated", "expired":
		return true
	}
	return false
}

func computeAggregateActiveStatus(active []responses.BookingResponse) string {
	if len(active) == 0 {
		return "cancelled"
	}
	// Priority across active stalls in the same booking:
	// pending_approval -> approved -> interested -> lease_applied -> completed
	hasPending := false
	hasApproved := false
	hasInterested := false
	hasLeaseApplied := false
	hasCompleted := false

	for _, item := range active {
		switch item.Status {
		case "pending_approval", "pending", "pending_slot":
			hasPending = true
		case "approved":
			hasApproved = true
		case "interested", "visited", "pending_decision":
			hasInterested = true
		case "lease_applied":
			hasLeaseApplied = true
		case "completed":
			hasCompleted = true
		}
	}

	if hasInterested {
		return "interested"
	}
	if hasApproved {
		return "approved"
	}
	if hasPending {
		return "pending_approval"
	}
	if hasLeaseApplied {
		return "lease_applied"
	}
	if hasCompleted {
		return "completed"
	}
	return active[0].Status
}

func ensureStatusMilestone(b *models.Booking, items *[]responses.TimelineItem) {
	if b == nil || items == nil {
		return
	}
	updatedStr := b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	actorName := b.ApplicantName
	if actorName == "" {
		actorName = "ผู้เช่า"
	}

	switch b.BookingStatus {
	case "lease_applied":
		hasApp := false
		for _, it := range *items {
			if strings.Contains(it.Action, "ยื่นคำร้อง") {
				hasApp = true
				break
			}
		}
		if !hasApp {
			*items = append(*items, responses.TimelineItem{
				Action:    "ยื่นคำร้องขอเช่าแผงค้า",
				ActorRole: "tenant",
				ActorName: actorName,
				Timestamp: updatedStr,
				Note:      "ส่งคำร้องขอเช่าแผงค้าเข้าสู่ระบบเรียบร้อยแล้ว อยู่ระหว่างตรวจสอบเอกสาร",
			})
		}
	case "completed":
		hasContract := false
		for _, it := range *items {
			if strings.Contains(it.Action, "สัญญา") {
				hasContract = true
				break
			}
		}
		if !hasContract {
			*items = append(*items, responses.TimelineItem{
				Action:    "อนุมัติสัญญาเช่าแผงค้าสำเร็จ",
				ActorRole: "officer",
				ActorName: "ผู้อำนวยการสำนักตลาด",
				Timestamp: updatedStr,
				Note:      "คำร้องผ่านการอนุมัติและจัดทำสัญญาเช่าเรียบร้อยแล้ว",
			})
		}
	case "terminated":
		hasTerm := false
		for _, it := range *items {
			if strings.Contains(it.Action, "สิ้นสุดสัญญา") {
				hasTerm = true
				break
			}
		}
		if !hasTerm {
			*items = append(*items, responses.TimelineItem{
				Action:    "สิ้นสุดสัญญาเช่า",
				ActorRole: "officer",
				ActorName: "เจ้าหน้าที่ตลาด",
				Timestamp: updatedStr,
				Note:      "สัญญาเช่าสิ้นสุดเรียบร้อยแล้ว แผงค้ากลับมาว่างอีกครั้ง",
			})
		}
	}
}

func parseTimeline(b *models.Booking) []responses.TimelineItem {
	var items []responses.TimelineItem
	if b.Timeline != "" {
		if err := json.Unmarshal([]byte(b.Timeline), &items); err == nil && len(items) > 0 {
			ensureStatusMilestone(b, &items)
			return items
		}
	}

	items = make([]responses.TimelineItem, 0)
	createdStr := b.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	updatedStr := b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")

	actorName := b.ApplicantName
	if actorName == "" {
		actorName = "ผู้เช่า"
	}

	initNote := ""
	if b.VisitDate != "" && b.VisitTime != "" {
		initNote = fmt.Sprintf("เลือกนัดหมายเข้าชมวันที่ %s รอบเวลา %s", b.VisitDate, b.VisitTime)
	}
	items = append(items, responses.TimelineItem{
		Action:    "ลงทะเบียนจองแผงและนัดหมายเข้าชม",
		ActorRole: "tenant",
		ActorName: actorName,
		Timestamp: createdStr,
		Note:      initNote,
	})

	switch b.BookingStatus {
	case "approved":
		items = append(items, responses.TimelineItem{
			Action:    "เจ้าหน้าที่อนุมัติการนัดหมาย",
			ActorRole: "officer",
			ActorName: "เจ้าหน้าที่ตลาด",
			Timestamp: updatedStr,
			Note:      "ตรวจสอบเอกสารเบื้องต้นและยืนยันรอบเวลานัดหมายเข้าชมเรียบร้อย",
		})
	case "visited", "interested":
		items = append(items, responses.TimelineItem{
			Action:    "เจ้าหน้าที่อนุมัติการนัดหมาย",
			ActorRole: "officer",
			ActorName: "เจ้าหน้าที่ตลาด",
			Timestamp: createdStr,
		})
		items = append(items, responses.TimelineItem{
			Action:    "ยืนยันการเข้าชมแผงเรียบร้อย",
			ActorRole: "officer",
			ActorName: "เจ้าหน้าที่ตลาด",
			Timestamp: updatedStr,
			Note:      "ผู้เช่าเดินทางมาชมแผงค้าตามนัดหมายเรียบร้อยแล้ว",
		})
	case "lease_applied":
		items = append(items, responses.TimelineItem{
			Action:    "เจ้าหน้าที่อนุมัติการนัดหมาย",
			ActorRole: "officer",
			ActorName: "เจ้าหน้าที่ตลาด",
			Timestamp: createdStr,
		})
		items = append(items, responses.TimelineItem{
			Action:    "ยืนยันการเข้าชมแผงเรียบร้อย",
			ActorRole: "officer",
			ActorName: "เจ้าหน้าที่ตลาด",
			Timestamp: updatedStr,
		})
		items = append(items, responses.TimelineItem{
			Action:    "ยื่นคำร้องขอเช่าแผงค้า",
			ActorRole: "tenant",
			ActorName: actorName,
			Timestamp: updatedStr,
			Note:      "ส่งคำร้องขอเช่าแผงค้าเข้าสู่ระบบเรียบร้อยแล้ว อยู่ระหว่างตรวจสอบเอกสาร",
		})
	case "completed":
		items = append(items, responses.TimelineItem{
			Action:    "อนุมัติสัญญาเช่าแผงค้าสำเร็จ",
			ActorRole: "officer",
			ActorName: "ผู้อำนวยการสำนักตลาด",
			Timestamp: updatedStr,
			Note:      "คำร้องผ่านการอนุมัติครบทุกขั้นตอน พร้อมทำสัญญาเช่า",
		})
	case "rejected":
		items = append(items, responses.TimelineItem{
			Action:    "ปฏิเสธการนัดหมาย",
			ActorRole: "officer",
			ActorName: "เจ้าหน้าที่ตลาด",
			Timestamp: updatedStr,
			Note:      b.Note,
		})
	case "cancelled":
		items = append(items, responses.TimelineItem{
			Action:    "ยกเลิกการนัดหมาย",
			ActorRole: "tenant",
			ActorName: actorName,
			Timestamp: updatedStr,
			Note:      b.Note,
		})
	case "no_show":
		items = append(items, responses.TimelineItem{
			Action:    "บันทึกไม่มาตามนัดหมาย",
			ActorRole: "officer",
			ActorName: "เจ้าหน้าที่ตลาด",
			Timestamp: updatedStr,
			Note:      "ผู้เช่าไม่เดินทางมาตามวันและเวลานัดหมาย",
		})
	}

	return items
}

func appendTimeline(b *models.Booking, item responses.TimelineItem) {
	current := parseTimeline(b)
	current = append(current, item)
	bytes, err := json.Marshal(current)
	if err == nil {
		b.Timeline = string(bytes)
	}
}

func isWaitingQueueBooking(b *models.Booking) bool {
	if b == nil {
		return false
	}
	switch b.BookingStatus {
	case "rejected", "cancelled", "not_interested", "no_show", "lease_applied", "completed", "terminated":
		return false
	}
	return b.BookingStatus == "waiting_queue" || b.QueueNumber > 0
}

func (u *useCase) toBookingResponse(b *models.Booking) responses.BookingResponse {
	normalizedStatus := b.BookingStatus
	if normalizedStatus == "waiting_queue" || normalizedStatus == "pending_slot" {
		normalizedStatus = "pending_approval"
	}
	isQueue := isWaitingQueueBooking(b)
	queuePos := 0
	if isQueue {
		queuePos = b.QueueNumber
		if queuePos <= 0 {
			queuePos = 1
		}
	}
	holdingStall := !isQueue && !isTerminalBookingStatus(normalizedStatus)

	rejectionReason := b.RejectionReason
	cancelReason := ""
	notInterestedReason := ""
	switch normalizedStatus {
	case "rejected":
		if rejectionReason == "" {
			rejectionReason = b.Note
		}
	case "cancelled":
		cancelReason = b.RejectionReason
		if cancelReason == "" {
			cancelReason = b.Note
		}
		if rejectionReason == "" {
			rejectionReason = cancelReason
		}
	case "not_interested":
		notInterestedReason = b.RejectionReason
		if notInterestedReason == "" {
			notInterestedReason = b.Note
		}
	}

	timelineItems := parseTimeline(b)
	appNumber := ""
	for _, item := range timelineItems {
		if idx := strings.Index(item.Note, "APP-"); idx != -1 {
			rest := item.Note[idx:]
			fields := strings.Fields(rest)
			if len(fields) > 0 {
				appNumber = strings.TrimRight(fields[0], ".,;)")
			}
		}
	}

	monthlyRate := b.MonthlyRent
	depositFee := b.DepositFee
	if monthlyRate <= 0 {
		rent, dep, _, _ := u.repo.LookupStallRate(b.StallIdent, b.MarketID)
		monthlyRate = rent
		depositFee = dep
	}
	if depositFee <= 0 {
		depositFee = monthlyRate
	}

	waitlistCount := 0
	if cnt, err := u.repo.CountWaitlistByStall(b.StallIdent); err == nil {
		waitlistCount = int(cnt)
	}

	activeStallCount := 1
	totalMonthlyRent := monthlyRate
	totalDeposit := depositFee
	if isTerminalBookingStatus(normalizedStatus) {
		activeStallCount = 0
		totalMonthlyRent = 0
		totalDeposit = 0
	}

	mID := 0
	if b.MarketID != nil {
		mID = *b.MarketID
	}
	leaseDays, _, _, _ := u.repo.GetPolicySettings(mID, b.MarketName)
	var expPtr *time.Time
	var remAppSeconds int64
	if normalizedStatus == "interested" {
		if b.ExpiresAt != nil {
			expPtr = b.ExpiresAt
		} else {
			base := b.UpdatedAt
			if base.IsZero() {
				base = time.Now()
			}
			fallbackExp := base.Add(time.Duration(leaseDays) * 24 * time.Hour)
			expPtr = &fallbackExp
		}
		if expPtr != nil {
			diff := int64(time.Until(*expPtr).Seconds())
			if diff > 0 {
				remAppSeconds = diff
			}
		}
	}

	return responses.BookingResponse{
		ID:                          b.ID,
		BookingNo:                   b.BookingNo,
		BookingNumber:               b.BookingNo,
		StallID:                     b.StallID,
		StallIdent:                  b.StallIdent,
		StallNumber:                 b.StallIdent,
		MarketID:                    b.MarketID,
		MarketName:                  b.MarketName,
		ProjectName:                 b.ProjectName,
		ApplicantName:               b.ApplicantName,
		TenantName:                  b.ApplicantName,
		TenantUserID:                b.CreatedBy,
		ApplicantPhone:              b.ApplicantPhone,
		TenantPhone:                 b.ApplicantPhone,
		ApplicantEmail:              b.ApplicantEmail,
		TenantEmail:                 b.ApplicantEmail,
		TenantIdCard:                b.ApplicantIdCard,
		BusinessType:                b.BusinessType,
		IntendedBusiness:            b.BusinessType,
		VisitDate:                   b.VisitDate,
		VisitTime:                   b.VisitTime,
		VisitTimeSlot:               b.VisitTime,
		MeetingPoint:                b.MeetingPoint,
		EscortOfficer:               b.EscortOfficer,
		OfficerNotes:                b.OfficerNotes,
		RejectionReason:             rejectionReason,
		CancelReason:                cancelReason,
		NotInterestedReason:         notInterestedReason,
		Note:                        b.Note,
		Status:                      normalizedStatus,
		ApplicationNumber:           appNumber,
		QueueNumber:                 queuePos,
		QueuePosition:               queuePos,
		HoldingStall:                holdingStall,
		AppointmentGroupID:          b.AppointmentGroupID,
		Timeline:                    timelineItems,
		ExpiresAt:                   expPtr,
		ApplicationDeadlineDays:     leaseDays,
		RemainingApplicationSeconds: remAppSeconds,
		CreatedAt:                   b.CreatedAt,
		StallCount:                  1,
		ActiveStallCount:            activeStallCount,
		TotalStallCount:             1,
		StallMonthlyRate:            monthlyRate,
		DepositFee:                  depositFee,
		TotalMonthlyRent:            totalMonthlyRent,
		TotalDeposit:                totalDeposit,
		WaitlistCount:               waitlistCount,
	}
}

var stallPrefixProjectMap = map[string]string{
	"BKA-01": "โครงการ 01 (แผงร้านค้า)",
	"BKA-02": "โครงการ 02 (อาคารพาณิชย์)",
	"BKA-03": "โครงการ 03 (พื้นที่ด้านซอย 127)",
	"BKA-04": "โครงการ 04 (พื้นที่ด้านซอย 125)",
	"BKH-01": "โครงการ 01 (5 เบ็ดเตล็ด)",
	"BKH-02": "โครงการ 02 (4 เบ็ดเตล็ด)",
	"BKH-03": "โครงการ 03 (3 อาหารสำเร็จรูป)",
	"BKH-04": "โครงการ 04 (6 เนื้อสัตว์ อาหารทะเล)",
	"BKH-05": "โครงการ 05 (1 ผัก-ผลไม้)",
	"BKH-06": "โครงการ 06 (2 อาหารปรุงและอาหารสำเร็จรูป)",
	"JJ-1":   "โครงการที่ 1",
	"JJ-3":   "โครงการที่ 3",
	"JJ-4":   "โครงการที่ 4",
	"JJ-5":   "โครงการที่ 5",
	"JJ-6":   "โครงการที่ 6",
	"JJ-7":   "โครงการที่ 7",
	"JJ-8":   "โครงการที่ 8",
	"JJ-9":   "โครงการที่ 9",
	"JJ-10":  "โครงการที่ 10",
	"JJ-11":  "โครงการที่ 11",
	"JJ-12":  "โครงการที่ 12",
	"JJ-13":  "โครงการที่ 13",
	"JJ-14":  "โครงการที่ 14",
	"JJ-15":  "โครงการที่ 15",
	"JJ-16":  "โครงการที่ 16",
	"JJ-17":  "โครงการที่ 17",
	"JJ-18":  "โครงการที่ 18",
	"JJ-19":  "โครงการที่ 19",
	"JJ-20":  "โครงการที่ 20/2",
	"JJ-21":  "โครงการที่ 21",
	"JJ-22":  "โครงการที่ 22",
	"JJ-23":  "โครงการที่ 23",
	"JJ-24":  "โครงการที่ 24",
	"JJ-25":  "โครงการที่ 25",
	"JJ-26":  "โครงการที่ 26",
	"JJ-27":  "โครงการที่ 27",
	"JJ-28":  "โครงการที่ 28",
	"JJ-29":  "โครงการที่ 29/2",
	"MIN-01": "โครงการ 01 (โซนหลัก (ชั้นล่าง))",
	"MIN-02": "โครงการ 02 (แผงค้าเบ็ดเตล็ด (ชั้นบน))",
	"MIN-03": "โครงการ 03 (โครงการที่ 3 (ชั้นล่าง))",
	"MIN-04": "โครงการ 04 (โครงการที่ 2 (ชั้นล่าง))",
	"MIN-05": "โครงการ 05 (โครงการที่ 4 (ชั้นล่าง))",
	"MIN-06": "โครงการ 06 (โครงการที่ 1 (ชั้นล่าง))",
	"MIN-07": "โครงการ 07 (โครงการที่ 7 (ชั้นล่าง))",
	"MIN-08": "โครงการ 08 (โครงการที่ 6 (ชั้นล่าง))",
	"MIN-09": "โครงการ 09 (โครงการที่ 8 (ชั้นล่าง))",
	"MIN-10": "โครงการ 10 (โครงการที่ 5 (ชั้นล่าง))",
	"MIN-11": "โครงการ 11 (แผงค้าต้นไม้ (ชั้นบน))",
	"MIN-12": "โครงการ 12 (โซนหลัก (ชั้นบน))",
	"NON-01": "โครงการ 01 (ด้านใน)",
	"NON-02": "โครงการ 02 (โซนหลัก)",
	"NON-03": "โครงการ 03 (ด้านนอก)",
	"NON-04": "โครงการ 04 (นอกตลาด)",
	"PRA-01": "โครงการ 01 (D)",
	"PRA-02": "โครงการ 02 (F)",
	"PRA-03": "โครงการ 03 (C)",
	"PRA-04": "โครงการ 04 (A)",
	"PRA-05": "โครงการ 05 (8)",
	"PRA-06": "โครงการ 06 (6)",
	"PRA-07": "โครงการ 07 (4)",
	"PRA-08": "โครงการ 08 (2)",
	"PRA-09": "โครงการ 09 (E)",
	"PRA-10": "โครงการ 10 (9)",
	"PRA-11": "โครงการ 11 (7)",
	"PRA-12": "โครงการ 12 (5)",
	"PRA-13": "โครงการ 13 (3)",
	"PRA-14": "โครงการ 14 (1)",
	"PRA-15": "โครงการ 15 (โซนหลัก)",
	"PRA-16": "โครงการ 16 (G)",
	"RBU-01": "โครงการ 01 (อาคารตลาดราษฎร์บูรณะ)",
	"RDP-01": "โครงการ 01 (M)",
	"RDP-02": "โครงการ 02 (B)",
	"RDP-03": "โครงการ 03 (D)",
	"RDP-04": "โครงการ 04 (F)",
	"RDP-05": "โครงการ 05 (E)",
	"RDP-06": "โครงการ 06 (C)",
	"RDP-07": "โครงการ 07 (A)",
	"RDP-08": "โครงการ 08 (G)",
	"RDP-09": "โครงการ 09 (H)",
	"RDP-10": "โครงการ 10 (โซนหลัก)",
	"SIN-01": "โครงการ 01",
	"THO-01": "โครงการ 01 (4)",
	"THO-02": "โครงการ 02 (3)",
	"THO-03": "โครงการ 03 (1)",
	"THO-04": "โครงการ 04 (11)",
	"THO-05": "โครงการ 05 (2)",
	"THO-06": "โครงการ 06 (6)",
	"THO-07": "โครงการ 07 (7)",
	"THO-08": "โครงการ 08 (8)",
	"THO-09": "โครงการ 09 (โซนหลัก)",
	"THO-10": "โครงการ 10 (5)",
	"THO-11": "โครงการ 11 (9)",
	"THO-12": "โครงการ 12 (1-2)",
	"THO-13": "โครงการ 13 (ด้านนอก)",
	"TWR-01": "โซน 1 (อาคารตลาดสดด้านใน)",
	"TWR-02": "โซน 2 (อาคารตลาดสดรอบนอก)",
	"TWR-03": "โซน 3 (ศูนย์อาหารเทเวศร์)",
	"TWR-04": "โซน 4 (ลานต้นโพธิ์)",
	"WON-01": "โครงการ 01 (ฉ)",
	"WON-02": "โครงการ 02 (ง)",
	"WON-03": "โครงการ 03 (จ)",
	"WON-04": "โครงการ 04 (ก)",
	"WON-05": "โครงการ 05 (ข)",
	"WON-06": "โครงการ 06 (ค)",
	"WON-07": "โครงการ 07 (โซนหลัก)",
}

func inferProjectNameFromStallIdent(stallIdent string, fallback string) string {
	code := strings.ToUpper(strings.TrimSpace(strings.Split(stallIdent, ",")[0]))
	parts := strings.Split(code, "-")
	if len(parts) >= 2 {
		prefix := parts[0] + "-" + parts[1]
		if mapped, ok := stallPrefixProjectMap[prefix]; ok && mapped != "" {
			return mapped
		}
	}
	return strings.TrimSpace(fallback)
}

func (u *useCase) groupBookingModels(bookings []models.Booking) []responses.BookingResponse {
	type groupBucket struct {
		key    string
		models []models.Booking
	}
	buckets := make(map[string]*groupBucket)
	order := make([]string, 0, len(bookings))

	for _, b := range bookings {
		key := strings.TrimSpace(b.AppointmentGroupID)
		if key == "" {
			key = fmt.Sprintf("single-%d", b.ID)
		}
		if _, exists := buckets[key]; !exists {
			buckets[key] = &groupBucket{key: key, models: make([]models.Booking, 0, 2)}
			order = append(order, key)
		}
		buckets[key].models = append(buckets[key].models, b)
	}

	result := make([]responses.BookingResponse, 0, len(order))
	for _, key := range order {
		bucket := buckets[key]
		// Sort children by ID ascending so the original booking number and stall order remain deterministic
		sort.Slice(bucket.models, func(i, j int) bool {
			return bucket.models[i].ID < bucket.models[j].ID
		})

		for i := range bucket.models {
			bucket.models[i].ProjectName = inferProjectNameFromStallIdent(bucket.models[i].StallIdent, bucket.models[i].ProjectName)
		}

		firstModel := bucket.models[0]
		u.repo.SyncBookingStalls(firstModel.ID, firstModel.BookingNo, firstModel.AppointmentGroupID, bucket.models)

		children := make([]responses.BookingResponse, len(bucket.models))
		for i := range bucket.models {
			children[i] = u.toBookingResponse(&bucket.models[i])
		}

		parent := children[0]
		// Keep canonical booking number from first stall in the group
		canonicalBookingNo := children[0].BookingNo

		activeChildren := make([]responses.BookingResponse, 0, len(children))
		activeCodes := make([]string, 0, len(children))
		allCodes := make([]string, 0, len(children))
		activeProjects := make([]string, 0, len(children))
		activeProjSeen := make(map[string]bool)
		allProjects := make([]string, 0, len(children))
		allProjSeen := make(map[string]bool)
		var activeRentSum float64
		var activeDepositSum float64
		var allRentSum float64
		var allDepositSum float64
		totalWaitlist := 0

		for _, child := range children {
			allCodes = append(allCodes, child.StallNumber)
			allRentSum += child.StallMonthlyRate
			allDepositSum += child.DepositFee
			if pName := strings.TrimSpace(child.ProjectName); pName != "" && !allProjSeen[pName] {
				allProjSeen[pName] = true
				allProjects = append(allProjects, pName)
			}
			totalWaitlist += child.WaitlistCount
			if !isTerminalBookingStatus(child.Status) {
				activeChildren = append(activeChildren, child)
				activeCodes = append(activeCodes, child.StallNumber)
				if pName := strings.TrimSpace(child.ProjectName); pName != "" && !activeProjSeen[pName] {
					activeProjSeen[pName] = true
					activeProjects = append(activeProjects, pName)
				}
				activeRentSum += child.StallMonthlyRate
				activeDepositSum += child.DepositFee
			}
		}

		if len(activeChildren) > 0 {
			// Use first active child as the primary reference for status/timeline/actions
			parent = activeChildren[0]
			parent.BookingNo = canonicalBookingNo
			parent.BookingNumber = canonicalBookingNo
			parent.Status = computeAggregateActiveStatus(activeChildren)
			parent.HoldingStall = true
			parent.StallIdent = strings.Join(activeCodes, ", ")
			parent.StallNumber = strings.Join(activeCodes, ", ")
			if len(activeProjects) > 0 {
				parent.ProjectName = strings.Join(activeProjects, ", ")
			}
			parent.StallCount = len(activeChildren)
			parent.ActiveStallCount = len(activeChildren)
			parent.TotalStallCount = len(children)
			parent.StallMonthlyRate = activeRentSum
			parent.DepositFee = activeDepositSum
			parent.TotalMonthlyRent = activeRentSum
			parent.TotalDeposit = activeDepositSum
		} else {
			// Only when ALL stalls in the booking are cancelled/rejected does the booking itself become cancelled/terminal
			lastChild := children[len(children)-1]
			parent = children[0]
			parent.BookingNo = canonicalBookingNo
			parent.BookingNumber = canonicalBookingNo
			if lastChild.Status == "rejected" || lastChild.Status == "no_show" || lastChild.Status == "not_interested" {
				parent.Status = lastChild.Status
			} else {
				parent.Status = "cancelled"
			}
			parent.HoldingStall = false
			parent.StallIdent = strings.Join(allCodes, ", ")
			parent.StallNumber = strings.Join(allCodes, ", ")
			if len(allProjects) > 0 {
				parent.ProjectName = strings.Join(allProjects, ", ")
			}
			parent.StallCount = len(children)
			parent.ActiveStallCount = 0
			parent.TotalStallCount = len(children)
			parent.StallMonthlyRate = allRentSum
			parent.DepositFee = allDepositSum
			parent.TotalMonthlyRent = allRentSum
			parent.TotalDeposit = allDepositSum
		}

		parent.WaitlistCount = totalWaitlist
		parent.SubBookings = children
		result = append(result, parent)
	}
	return result
}

func canManageBookings(role string) bool {
	return role == "admin" || role == "staff" || role == "market_director"
}

func (u *useCase) GetBookings(filter requests.BookingFilter, actorID int, actorRole string) ([]responses.BookingResponse, error) {
	if actorID <= 0 {
		return nil, responses.NewAppErr(http.StatusUnauthorized, "กรุณาเข้าสู่ระบบก่อนดูรายการจอง")
	}

	var ownerID *int
	if !canManageBookings(actorRole) {
		ownerID = &actorID
	}

	bookings, err := u.repo.FindAll(filter, ownerID)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	return u.groupBookingModels(bookings), nil
}

func (u *useCase) GetBookingByID(id int, actorID int, actorRole string) (*responses.BookingResponse, error) {
	if actorID <= 0 {
		return nil, responses.NewAppErr(http.StatusUnauthorized, "กรุณาเข้าสู่ระบบก่อนดูรายละเอียดการจอง")
	}

	booking, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลการจอง")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	if !canManageBookings(actorRole) && (booking.CreatedBy == nil || *booking.CreatedBy != actorID) {
		return nil, responses.NewAppErr(http.StatusForbidden, "คุณไม่มีสิทธิ์ดูรายละเอียดการจองนี้")
	}

	if strings.TrimSpace(booking.AppointmentGroupID) != "" {
		if siblings, err := u.repo.FindByGroupID(booking.AppointmentGroupID); err == nil && len(siblings) > 0 {
			grouped := u.groupBookingModels(siblings)
			if len(grouped) > 0 {
				return &grouped[0], nil
			}
		}
	}

	grouped := u.groupBookingModels([]models.Booking{*booking})
	return &grouped[0], nil
}

func (u *useCase) GetPolicySettings(marketID int, marketName string) (*responses.BookingPolicySettingsResponse, error) {
	leaseDays, offerHours, resolvedID, resolvedName := u.repo.GetPolicySettings(marketID, marketName)
	return &responses.BookingPolicySettingsResponse{
		MarketID:                     resolvedID,
		MarketName:                   resolvedName,
		LeaseApplicationDeadlineDays: leaseDays,
		WaitlistOfferHours:           offerHours,
	}, nil
}

func (u *useCase) UpdatePolicySettings(actorRole string, req requests.UpdateBookingPolicySettingsRequest) (*responses.BookingPolicySettingsResponse, error) {
	if !canManageBookings(actorRole) {
		return nil, responses.NewAppErr(http.StatusForbidden, "เฉพาะเจ้าหน้าที่เท่านั้นที่แก้ไขนโยบายการจองได้")
	}
	if req.LeaseApplicationDeadlineDays < 1 || req.LeaseApplicationDeadlineDays > 30 {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กำหนดยื่นคำร้องต้องอยู่ระหว่าง 1–30 วัน")
	}
	if req.WaitlistOfferHours < 1 || req.WaitlistOfferHours > 168 {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "เวลาตอบรับ waiting list ต้องอยู่ระหว่าง 1–168 ชั่วโมง")
	}
	result, err := u.repo.UpdatePolicySettings(req.MarketID, req.MarketName, req.LeaseApplicationDeadlineDays, req.WaitlistOfferHours)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return result, nil
}

func (u *useCase) CreateBooking(userID int, req requests.CreateBookingRequest) (*responses.BookingResponse, error) {
	if userID <= 0 {
		return nil, responses.NewAppErr(http.StatusUnauthorized, "กรุณาเข้าสู่ระบบก่อนจองแผง")
	}

	// Enrich applicant contact details from authenticated user if omitted
	if strings.TrimSpace(req.ApplicantName) == "" || strings.TrimSpace(req.ApplicantPhone) == "" {
		if user, err := u.repo.LookupUserByID(userID); err == nil && user != nil {
			if strings.TrimSpace(req.ApplicantName) == "" {
				req.ApplicantName = strings.TrimSpace(user.FullName)
			}
			if strings.TrimSpace(req.ApplicantPhone) == "" {
				req.ApplicantPhone = strings.TrimSpace(user.Phone)
			}
			if strings.TrimSpace(req.ApplicantEmail) == "" {
				req.ApplicantEmail = strings.TrimSpace(user.Email)
			}
		}
	}

	if len(req.Stalls) > 0 {
		groupID := strings.TrimSpace(req.AppointmentGroupID)
		if groupID == "" {
			groupID = fmt.Sprintf("BKG-%d", time.Now().UnixMicro())
		}
		seen := make(map[string]bool)
		// Pre-validate all stalls in the booking so all stalls share the same valid bookable status
		for _, selected := range req.Stalls {
			ident := strings.TrimSpace(selected.StallIdent)
			if ident == "" {
				ident = strings.TrimSpace(selected.StallNumber)
			}
			if ident == "" {
				return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "พบแผงที่ไม่มีรหัส")
			}
			if seen[ident] {
				return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "ไม่สามารถเลือกแผงเดิมซ้ำในใบจองเดียวกัน")
			}
			seen[ident] = true

			if offer, err := u.repo.FindActiveWaitlistOffer(ident); err == nil && offer != nil && offer.TenantID != userID {
				return nil, responses.NewAppErr(http.StatusConflict, fmt.Sprintf("แผง %s กำลังเสนอสิทธิ์ให้ผู้รอคิวลำดับแรก ไม่สามารถเพิ่มเข้าใบจองได้", ident))
			}
			hasHolder, err := u.repo.HasActiveHolder(ident)
			if err != nil {
				return nil, responses.NewInternalServerError()
			}
			if hasHolder || req.IsQueue {
				return nil, responses.NewAppErr(http.StatusConflict, fmt.Sprintf("แผง %s ติดจองอยู่ ไม่สามารถเพิ่มเข้าใบจองได้ กรุณาเข้าคิวรอแผงนี้แทน", ident))
			}
		}

		skipChildNotifs := len(req.Stalls) > 1
		for _, selected := range req.Stalls {
			ident := strings.TrimSpace(selected.StallIdent)
			if ident == "" {
				ident = strings.TrimSpace(selected.StallNumber)
			}
			childReq := req
			childReq.Stalls = nil
			childReq.StallID = selected.StallID
			childReq.StallIdent = ident
			childReq.StallNumber = ident
			childReq.AppointmentGroupID = groupID
			childReq.SkipNotification = skipChildNotifs
			if strings.TrimSpace(selected.ProjectName) != "" {
				childReq.ProjectName = selected.ProjectName
			}
			if selected.MonthlyRent > 0 {
				childReq.MonthlyRent = selected.MonthlyRent
			}
			if selected.DepositFee > 0 {
				childReq.DepositFee = selected.DepositFee
			}
			if _, err := u.CreateBooking(userID, childReq); err != nil {
				return nil, err
			}
		}

		siblings, err := u.repo.FindByGroupID(groupID)
		if err != nil || len(siblings) == 0 {
			return nil, responses.NewInternalServerError()
		}

		if skipChildNotifs && len(siblings) > 0 {
			primary := &siblings[0]
			var stallIdents []string
			var projNames []string
			seenProj := make(map[string]bool)
			for _, sib := range siblings {
				stallIdents = append(stallIdents, sib.StallIdent)
				pName := inferProjectNameFromStallIdent(sib.StallIdent, sib.ProjectName)
				if pName != "" && !seenProj[pName] {
					seenProj[pName] = true
					projNames = append(projNames, pName)
				}
			}
			combinedStalls := strings.Join(stallIdents, ", ")
			combinedProjects := strings.Join(projNames, ", ")
			statusLabel := "รออนุมัติวันนัด"
			bookingID := primary.ID

			officerNotif := models.Notification{
				RecipientRole: "officer", RecipientName: "เจ้าหน้าที่ดูแลตลาด",
				Type: primary.BookingStatus, StatusCode: primary.BookingStatus, StatusLabel: statusLabel,
				Title:     fmt.Sprintf("มีคำขอนัดชมแผงใหม่ แผง %s (%s)", combinedStalls, primary.BookingNo),
				Message:   fmt.Sprintf("%s ส่งคำขอนัดชมแผง %s (%d แผง)", primary.ApplicantName, combinedStalls, len(siblings)),
				BookingID: &bookingID, BookingNumber: primary.BookingNo, StallNumber: combinedStalls,
				MarketName: primary.MarketName, ProjectName: combinedProjects,
				TargetURL: "/bookings/" + fmt.Sprint(primary.ID) + "?role=officer",
			}
			_ = u.repo.CreateNotification(&officerNotif)
		}

		grouped := u.groupBookingModels(siblings)
		return &grouped[0], nil
	}

	// Single-stall validation (or individual stall within a multi-stall group)
	if offer, err := u.repo.FindActiveWaitlistOffer(req.StallIdent); err == nil && offer != nil && offer.TenantID != userID {
		return nil, responses.NewAppErr(http.StatusConflict, "แผงนี้กำลังเปิดสิทธิ์ให้ผู้รอคิวลำดับแรก")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	if existing, err := u.repo.FindActiveDuplicate(req.StallIdent, userID, req.ApplicantName, req.ApplicantPhone); err == nil && existing != nil {
		return nil, responses.NewAppErr(http.StatusConflict, fmt.Sprintf("คุณได้ทำการจองแผง %s ไว้แล้ว (รหัส %s) ไม่สามารถจองซ้ำได้", req.StallIdent, existing.BookingNo))
	}

	hasActiveHolder, err := u.repo.HasActiveHolder(req.StallIdent)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	if hasActiveHolder || req.IsQueue {
		return nil, responses.NewAppErr(http.StatusConflict, fmt.Sprintf("แผง %s ติดจองอยู่ ไม่สามารถเพิ่มเข้าใบจองได้ กรุณากดปุ่ม 'รอคิวแผงนี้' เพื่อเข้าคิวรอแทน", req.StallIdent))
	}

	visitDate := strings.TrimSpace(req.VisitDate)
	visitTime := strings.TrimSpace(req.VisitTime)
	if visitDate == "" || visitTime == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาเลือกวันและเวลานัดหมายเข้าชมแผง")
	}

	bookingNo := fmt.Sprintf("BK-%s-%06d", time.Now().Format("20060102"), time.Now().UnixMicro()%1000000)

	email := strings.TrimSpace(req.ApplicantEmail)
	if email == "" {
		email = strings.TrimSpace(req.TenantEmail)
	}
	idCard := strings.TrimSpace(req.ApplicantIdCard)
	if idCard == "" {
		idCard = strings.TrimSpace(req.TenantIdCard)
	}

	var parsedStallID *int
	switch v := req.StallID.(type) {
	case float64:
		idVal := int(v)
		parsedStallID = &idVal
	case int:
		idVal := v
		parsedStallID = &idVal
	}

	groupID := strings.TrimSpace(req.AppointmentGroupID)
	if groupID == "" {
		groupID = fmt.Sprintf("BKG-%d", time.Now().UnixMicro())
	}

	booking := models.Booking{
		BookingNo:          bookingNo,
		StallID:            parsedStallID,
		StallIdent:         req.StallIdent,
		MarketID:           req.MarketID,
		MarketName:         req.MarketName,
		ProjectName:        req.ProjectName,
		MonthlyRent:        req.MonthlyRent,
		DepositFee:         req.DepositFee,
		ApplicantName:      req.ApplicantName,
		ApplicantPhone:     req.ApplicantPhone,
		ApplicantEmail:     email,
		ApplicantIdCard:    idCard,
		BusinessType:       req.BusinessType,
		VisitDate:          visitDate,
		VisitTime:          visitTime,
		Note:               req.Note,
		BookingStatus:      "pending_approval",
		QueueNumber:        0,
		AppointmentGroupID: groupID,
		SkipNotification:   req.SkipNotification,
	}
	if userID > 0 {
		booking.CreatedBy = &userID
		booking.UpdatedBy = &userID
	}

	initTimeline := []responses.TimelineItem{
		{
			Action:    "ลงทะเบียนจองแผงและนัดหมายเข้าชม",
			ActorRole: "tenant",
			ActorName: func() string {
				if req.ApplicantName != "" {
					return req.ApplicantName
				}
				return "ผู้เช่า"
			}(),
			Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
			Note:      fmt.Sprintf("เลือกนัดหมายเข้าชมวันที่ %s รอบเวลา %s", visitDate, visitTime),
		},
	}
	if tBytes, err := json.Marshal(initTimeline); err == nil {
		booking.Timeline = string(tBytes)
	}

	if err := u.repo.Create(&booking); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	_ = u.repo.AcceptWaitlistOffer(req.StallIdent, userID)

	if siblings, err := u.repo.FindByGroupID(groupID); err == nil && len(siblings) > 0 {
		grouped := u.groupBookingModels(siblings)
		return &grouped[0], nil
	}
	grouped := u.groupBookingModels([]models.Booking{booking})
	return &grouped[0], nil
}

func toWaitlistResponse(item *models.Waitlist) responses.WaitlistResponse {
	var remSec int64
	expPtr := item.ExpiresAt
	if item.Status == "offered" && expPtr == nil {
		base := item.UpdatedAt
		if item.OfferedAt != nil && !item.OfferedAt.IsZero() {
			base = *item.OfferedAt
		} else if base.IsZero() {
			base = time.Now()
		}
		fallback := base.Add(48 * time.Hour)
		expPtr = &fallback
	}
	if item.Status == "offered" && expPtr != nil {
		diff := int64(time.Until(*expPtr).Seconds())
		if diff > 0 {
			remSec = diff
		}
	}
	dismissReason := strings.TrimSpace(item.DismissReason)
	if dismissReason == "" {
		switch item.Status {
		case "cancelled":
			dismissReason = "ตัดสิทธิ์เนื่องจากผู้จองหลักทำสัญญาเช่าแผงสำเร็จ"
		case "expired":
			dismissReason = "หมดเวลาตอบรับสิทธิ์ 48 ชั่วโมง (สละสิทธิ์)"
		case "withdrawn":
			dismissReason = "ยกเลิกการต่อคิว"
		}
	}
	return responses.WaitlistResponse{
		ID:                    item.ID,
		StallID:               item.StallID,
		StallIdent:            item.StallIdent,
		StallNumber:           item.StallIdent,
		MarketID:              item.MarketID,
		MarketName:            item.MarketName,
		ProjectName:           inferProjectNameFromStallIdent(item.StallIdent, item.ProjectName),
		TenantID:              item.TenantID,
		TenantName:            item.TenantName,
		TenantPhone:           item.TenantPhone,
		Position:              item.Position,
		Status:                item.Status,
		DismissReason:         dismissReason,
		TenantDismissed:       item.TenantDismissed,
		OfferedAt:             item.OfferedAt,
		ExpiresAt:             expPtr,
		OfferExpiresAt:        expPtr,
		RemainingSeconds:      remSec,
		OfferRemainingSeconds: remSec,
		CreatedAt:             item.CreatedAt,
		UpdatedAt:             item.UpdatedAt,
	}
}

func (u *useCase) CreateWaitlist(userID int, req requests.CreateWaitlistRequest) (*responses.WaitlistResponse, error) {
	if userID <= 0 {
		return nil, responses.NewAppErr(http.StatusUnauthorized, "กรุณาเข้าสู่ระบบก่อนเข้าคิวแผง")
	}
	ident := strings.TrimSpace(req.StallIdent)
	if ident == "" {
		ident = strings.TrimSpace(req.StallNumber)
	}
	if ident == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุรหัสแผง")
	}
	// Check if user already holds this stall as primary holder
	if existing, err := u.repo.FindActiveDuplicate(ident, userID, "", ""); err == nil && existing != nil {
		return nil, responses.NewAppErr(http.StatusConflict, fmt.Sprintf("คุณเป็นผู้จองหลักของแผง %s อยู่แล้ว ไม่สามารถเข้าคิวซ้ำได้", ident))
	}
	hasHolder, err := u.repo.HasActiveHolder(ident)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	activeOffer, offerErr := u.repo.FindActiveWaitlistOffer(ident)
	if offerErr != nil && !errors.Is(offerErr, gorm.ErrRecordNotFound) {
		logs.Error(offerErr)
		return nil, responses.NewInternalServerError()
	}
	if !hasHolder && activeOffer == nil {
		return nil, responses.NewAppErr(http.StatusConflict, fmt.Sprintf("แผง %s ยังว่างอยู่ สามารถจองได้ทันทีโดยไม่ต้องเข้าคิว", ident))
	}

	var stallID *int
	switch v := req.StallID.(type) {
	case float64:
		n := int(v)
		stallID = &n
	case int:
		n := v
		stallID = &n
	}

	tenantName := strings.TrimSpace(req.TenantName)
	tenantPhone := strings.TrimSpace(req.TenantPhone)
	if user, err := u.repo.LookupUserByID(userID); err == nil && user != nil {
		if tenantName == "" {
			tenantName = strings.TrimSpace(user.FullName)
		}
		if tenantPhone == "" {
			tenantPhone = strings.TrimSpace(user.Phone)
		}
	}
	if tenantName == "" {
		tenantName = fmt.Sprintf("ผู้เช่า #%d", userID)
	}

	mktName := strings.TrimSpace(req.MarketName)
	projName := strings.TrimSpace(req.ProjectName)
	if mktName == "" || projName == "" {
		_, _, dbProj, dbMkt := u.repo.LookupStallRate(ident, req.MarketID)
		if mktName == "" {
			mktName = dbMkt
		}
		if projName == "" {
			projName = dbProj
		}
	}

	item := &models.Waitlist{
		StallID:     stallID,
		StallIdent:  ident,
		MarketID:    req.MarketID,
		MarketName:  mktName,
		ProjectName: projName,
		TenantID:    userID,
		TenantName:  tenantName,
		TenantPhone: tenantPhone,
	}
	if err := u.repo.CreateWaitlist(item); err != nil {
		if strings.Contains(err.Error(), "WAITLIST_DUPLICATE") {
			return nil, responses.NewAppErr(http.StatusConflict, "คุณอยู่ในคิวรอแผงนี้แล้ว")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	returnPtr := toWaitlistResponse(item)
	return &returnPtr, nil
}

func (u *useCase) GetMyWaitlist(userID int) ([]responses.WaitlistResponse, error) {
	items, err := u.repo.FindWaitlistByTenant(userID)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	result := make([]responses.WaitlistResponse, len(items))
	for i := range items {
		result[i] = toWaitlistResponse(&items[i])
	}
	return result, nil
}

func (u *useCase) GetWaitlistByStall(stallIdent string, actorID int, actorRole string) ([]responses.WaitlistResponse, error) {
	items, err := u.repo.FindWaitlistByStall(stallIdent)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	result := make([]responses.WaitlistResponse, len(items))
	for i := range items {
		resp := toWaitlistResponse(&items[i])
		if !canManageBookings(actorRole) && items[i].TenantID != actorID {
			resp.TenantPhone = ""
		}
		result[i] = resp
	}
	return result, nil
}

func (u *useCase) WithdrawWaitlist(id int, userID int, actorRole string) error {
	item, err := u.repo.FindWaitlistByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.NewAppErr(http.StatusNotFound, "ไม่พบรายการรอคิว")
		}
		return responses.NewInternalServerError()
	}
	if !canManageBookings(actorRole) && item.TenantID != userID {
		return responses.NewAppErr(http.StatusForbidden, "คุณไม่มีสิทธิ์ออกจากคิวนี้")
	}
	if err := u.repo.WithdrawWaitlist(item); err != nil {
		logs.Error(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (u *useCase) AdvanceWaitlistByStall(stallIdent string, actorID int, actorRole string) ([]responses.WaitlistResponse, error) {
	if !canManageBookings(actorRole) {
		return nil, responses.NewAppErr(http.StatusForbidden, "เฉพาะเจ้าหน้าที่เท่านั้นที่สามารถเลื่อนคิวได้")
	}
	if err := u.repo.AdvanceWaitlistByStall(stallIdent); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return u.GetWaitlistByStall(stallIdent, actorID, actorRole)
}

func (u *useCase) CancelBookingStall(id int, stallIdent string, userID int, actorRole string, reason string) (*responses.BookingResponse, error) {
	cleanIdent := strings.TrimSpace(stallIdent)
	if cleanIdent == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุรหัสแผงที่ต้องการยกเลิก")
	}
	return u.UpdateStatus(id, userID, requests.UpdateBookingStatusRequest{
		Status:          "cancelled",
		Reason:          reason,
		RejectionReason: reason,
		StallIdent:      cleanIdent,
		ByOfficer:       canManageBookings(actorRole),
		ActorRole:       actorRole,
	})
}

func (u *useCase) applyStatusToBookingRow(booking *models.Booking, userID int, req requests.UpdateBookingStatusRequest) error {
	if (req.Status == "interested" || req.Status == "visited") && isWaitingQueueBooking(booking) {
		return responses.NewAppErr(http.StatusBadRequest, "ยังไม่สามารถกดยืนยันการเช่าได้ เนื่องจากผู้จองคิวหลักยังไม่สละสิทธิ์")
	}

	wasQueue := isWaitingQueueBooking(booking)
	targetStatus := req.Status
	if targetStatus == "visited" {
		targetStatus = "interested"
	}
	booking.BookingStatus = targetStatus
	booking.UpdatedBy = &userID
	if strings.TrimSpace(req.MeetingPoint) != "" {
		booking.MeetingPoint = strings.TrimSpace(req.MeetingPoint)
	}
	if strings.TrimSpace(req.EscortOfficer) != "" {
		booking.EscortOfficer = strings.TrimSpace(req.EscortOfficer)
	}
	if strings.TrimSpace(req.OfficerNotes) != "" {
		booking.OfficerNotes = strings.TrimSpace(req.OfficerNotes)
	}

	actionName := "อัปเดตสถานะการจอง"
	actorRole := "officer"
	actorName := "เจ้าหน้าที่ตลาด"

	switch targetStatus {
	case "approved":
		actionName = fmt.Sprintf("เจ้าหน้าที่อนุมัติการนัดหมาย (แผง %s)", booking.StallIdent)
	case "interested":
		actionName = fmt.Sprintf("ยืนยันการเช่าแผง %s เรียบร้อย (เปิดให้ผู้เช่ายื่นคำร้อง)", booking.StallIdent)
	case "rejected":
		actionName = fmt.Sprintf("ปฏิเสธการนัดหมาย (แผง %s)", booking.StallIdent)
	case "cancelled":
		if req.ByOfficer || canManageBookings(req.ActorRole) {
			actionName = fmt.Sprintf("เจ้าหน้าที่ยกเลิกรายการจองแผง %s", booking.StallIdent)
			actorRole = "officer"
			actorName = "เจ้าหน้าที่ตลาด"
		} else {
			if wasQueue {
				actionName = fmt.Sprintf("ผู้เช่าสละสิทธิ์คิวแผง %s", booking.StallIdent)
			} else {
				actionName = fmt.Sprintf("ผู้เช่ายกเลิกการจองแผง %s", booking.StallIdent)
			}
			actorRole = "tenant"
			actorName = booking.ApplicantName
		}
	case "not_interested":
		actionName = fmt.Sprintf("ผู้เช่าสละสิทธิ์ไม่ยื่นคำร้องขอเช่าแผง %s", booking.StallIdent)
		if !req.ByOfficer && !canManageBookings(req.ActorRole) {
			actorRole = "tenant"
			actorName = booking.ApplicantName
		}
	case "no_show":
		actionName = fmt.Sprintf("บันทึกไม่มาตามนัดหมาย (แผง %s)", booking.StallIdent)
	case "lease_applied":
		actionName = fmt.Sprintf("ยื่นคำร้องขอเช่าแผง %s", booking.StallIdent)
		actorRole = "tenant"
		actorName = booking.ApplicantName
	case "completed":
		actionName = fmt.Sprintf("อนุมัติสัญญาเช่าแผง %s สำเร็จ", booking.StallIdent)
		actorName = "ผู้อำนวยการสำนักตลาด"
	}

	note := strings.TrimSpace(req.Reason)
	if note == "" {
		note = strings.TrimSpace(req.RejectionReason)
	}
	if note == "" {
		note = strings.TrimSpace(req.Note)
	}
	if note == "" {
		note = strings.TrimSpace(req.OfficerNotes)
	}
	if targetStatus == "rejected" || targetStatus == "cancelled" || targetStatus == "not_interested" || targetStatus == "no_show" {
		if strings.TrimSpace(req.RejectionReason) != "" && strings.TrimSpace(req.OfficerNotes) != "" && req.RejectionReason != req.OfficerNotes {
			booking.RejectionReason = fmt.Sprintf("%s (%s)", strings.TrimSpace(req.RejectionReason), strings.TrimSpace(req.OfficerNotes))
		} else if note != "" {
			booking.RejectionReason = note
		}
		if booking.RejectionReason != "" {
			booking.Note = booking.RejectionReason
		}
	}

	appendTimeline(booking, responses.TimelineItem{
		Action:    actionName,
		ActorRole: actorRole,
		ActorName: actorName,
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Note:      note,
	})

	return u.repo.Update(booking)
}

func (u *useCase) UpdateStatus(id int, userID int, req requests.UpdateBookingStatusRequest) (*responses.BookingResponse, error) {
	booking, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลการจอง")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	if !canManageBookings(req.ActorRole) && (booking.CreatedBy == nil || *booking.CreatedBy != userID) {
		return nil, responses.NewAppErr(http.StatusForbidden, "คุณไม่มีสิทธิ์แก้ไขรายการจองนี้")
	}

	groupID := strings.TrimSpace(booking.AppointmentGroupID)
	targetStallIdent := strings.TrimSpace(req.StallIdent)

	if groupID != "" {
		siblings, err := u.repo.FindByGroupID(groupID)
		if err == nil && len(siblings) > 1 {
			// Case 1: Specific stall requested via req.StallIdent
			if targetStallIdent != "" {
				found := false
				for i := range siblings {
					if strings.EqualFold(strings.TrimSpace(siblings[i].StallIdent), targetStallIdent) {
						found = true
						if err := u.applyStatusToBookingRow(&siblings[i], userID, req); err != nil {
							return nil, err
						}
						break
					}
				}
				if !found {
					return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบแผงที่ระบุในใบจองนี้")
				}
				updatedSiblings, _ := u.repo.FindByGroupID(groupID)
				grouped := u.groupBookingModels(updatedSiblings)
				return &grouped[0], nil
			}

			// Case 2: Cancel without req.StallIdent and without req.CancelAll -> only cancel the targeted row `id`
			// so cancelling 1 stall (`POST /bookings/:subId/cancel`) never cancels the entire multi-stall booking automatically!
			if req.Status == "cancelled" && !req.CancelAll {
				if err := u.applyStatusToBookingRow(booking, userID, req); err != nil {
					return nil, err
				}
				updatedSiblings, _ := u.repo.FindByGroupID(groupID)
				grouped := u.groupBookingModels(updatedSiblings)
				return &grouped[0], nil
			}

			// Case 3: Whole-booking action (Approve visit, Confirm visit, Reject all, No-show all, or CancelAll)
			var affectedSiblings []*models.Booking
			for i := range siblings {
				if isTerminalBookingStatus(siblings[i].BookingStatus) {
					continue
				}
				siblings[i].SkipNotification = true
				if err := u.applyStatusToBookingRow(&siblings[i], userID, req); err != nil {
					return nil, err
				}
				affectedSiblings = append(affectedSiblings, &siblings[i])
			}
			if len(affectedSiblings) > 0 {
				primary := affectedSiblings[0]
				var stallIdents []string
				var projNames []string
				seenProj := make(map[string]bool)
				for _, sib := range affectedSiblings {
					stallIdents = append(stallIdents, sib.StallIdent)
					pName := inferProjectNameFromStallIdent(sib.StallIdent, sib.ProjectName)
					if pName != "" && !seenProj[pName] {
						seenProj[pName] = true
						projNames = append(projNames, pName)
					}
				}
				combinedStalls := strings.Join(stallIdents, ", ")
				combinedProjects := strings.Join(projNames, ", ")

				statusLabels := map[string]string{
					"approved":         "นัดหมายสำเร็จ",
					"rejected":         "ไม่อนุมัตินัดหมาย",
					"interested":       "พร้อมยื่นคำร้องขอเช่า",
					"cancelled":        "ยกเลิกรายการแล้ว",
					"pending_approval": "รออนุมัติวันนัด",
					"no_show":          "ไม่มาตามนัดหมาย",
					"not_interested":   "สละสิทธิ์การเช่า",
				}
				label := statusLabels[primary.BookingStatus]
				if label == "" {
					label = primary.BookingStatus
				}
				notifTitle := label + " แผง " + combinedStalls + " (" + primary.BookingNo + ")"
				notifMsg := fmt.Sprintf("สถานะรายการจองแผง %s (%d แผง) เปลี่ยนเป็น %s", combinedStalls, len(affectedSiblings), label)
				if primary.BookingStatus == "interested" {
					mID := 0
					if primary.MarketID != nil {
						mID = *primary.MarketID
					}
					leaseDays, _, _, _ := u.repo.GetPolicySettings(mID, primary.MarketName)
					notifTitle = "เปิดสิทธิ์ยื่นคำร้องขอเช่าแผง " + combinedStalls + " (" + primary.BookingNo + ")"
					notifMsg = fmt.Sprintf("เจ้าหน้าที่ยืนยันการเช่าแผง %s (%d แผง) เรียบร้อยแล้ว กรุณากดปุ่ม 'ยื่นคำร้องขอเช่าแผงนี้' ภายใน %d วัน", combinedStalls, len(affectedSiblings), leaseDays)
				}
				bID := primary.ID
				groupNotif := models.Notification{
					RecipientRole: "tenant", RecipientUserID: primary.CreatedBy, RecipientName: primary.ApplicantName,
					Type: primary.BookingStatus, StatusCode: primary.BookingStatus, StatusLabel: label,
					Title:     notifTitle,
					Message:   notifMsg,
					BookingID: &bID, BookingNumber: primary.BookingNo, StallNumber: combinedStalls,
					MarketName: primary.MarketName, ProjectName: combinedProjects,
					TargetURL: "/bookings/" + fmt.Sprint(primary.ID) + "?role=tenant",
				}
				_ = u.repo.CreateNotification(&groupNotif)
			}
			updatedSiblings, _ := u.repo.FindByGroupID(groupID)
			grouped := u.groupBookingModels(updatedSiblings)
			return &grouped[0], nil
		}
	}

	if err := u.applyStatusToBookingRow(booking, userID, req); err != nil {
		return nil, err
	}

	if groupID != "" {
		if updatedSiblings, err := u.repo.FindByGroupID(groupID); err == nil && len(updatedSiblings) > 0 {
			grouped := u.groupBookingModels(updatedSiblings)
			return &grouped[0], nil
		}
	}
	grouped := u.groupBookingModels([]models.Booking{*booking})
	return &grouped[0], nil
}

func (u *useCase) UpdateVisit(id int, actorID int, actorRole string, visitDate string, visitTime string) (*responses.BookingResponse, error) {
	booking, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลการจอง")
		}
		return nil, responses.NewInternalServerError()
	}
	if !canManageBookings(actorRole) && (booking.CreatedBy == nil || *booking.CreatedBy != actorID) {
		return nil, responses.NewAppErr(http.StatusForbidden, "คุณไม่มีสิทธิ์เปลี่ยนวันนัดหมายของรายการนี้")
	}
	if visitDate == "" || visitTime == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุวันและเวลานัดหมาย")
	}

	groupID := strings.TrimSpace(booking.AppointmentGroupID)
	if groupID != "" {
		if siblings, err := u.repo.FindByGroupID(groupID); err == nil && len(siblings) > 0 {
			for i := range siblings {
				if isTerminalBookingStatus(siblings[i].BookingStatus) {
					continue
				}
				siblings[i].VisitDate = visitDate
				siblings[i].VisitTime = visitTime
				siblings[i].BookingStatus = "pending_approval"
				siblings[i].UpdatedBy = &actorID
				appendTimeline(&siblings[i], responses.TimelineItem{
					Action:    "เลือกวัน-เวลานัดหมายเข้าชมแผง",
					ActorRole: actorRole,
					ActorName: siblings[i].ApplicantName,
					Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
					Note:      fmt.Sprintf("เปลี่ยนรอบนัดหมายเป็นวันที่ %s รอบเวลา %s", visitDate, visitTime),
				})
				_ = u.repo.Update(&siblings[i])
			}
			updatedSiblings, _ := u.repo.FindByGroupID(groupID)
			grouped := u.groupBookingModels(updatedSiblings)
			return &grouped[0], nil
		}
	}

	booking.VisitDate = visitDate
	booking.VisitTime = visitTime
	booking.BookingStatus = "pending_approval"
	booking.UpdatedBy = &actorID

	appendTimeline(booking, responses.TimelineItem{
		Action:    "เลือกวัน-เวลานัดหมายเข้าชมแผง",
		ActorRole: actorRole,
		ActorName: booking.ApplicantName,
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Note:      fmt.Sprintf("เปลี่ยนรอบนัดหมายเป็นวันที่ %s รอบเวลา %s", visitDate, visitTime),
	})
	if err := u.repo.Update(booking); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	grouped := u.groupBookingModels([]models.Booking{*booking})
	return &grouped[0], nil
}

func (u *useCase) ClearAll(actorRole string) error {
	if err := u.repo.ClearAll(); err != nil {
		logs.Error(err)
		return responses.NewInternalServerError()
	}
	return nil
}
