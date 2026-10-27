package usecases_test

import (
	"fmt"
	"testing"
	"time"

	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	bookingrepo "emarket-api/src/booking/repositories"
	bookinguc "emarket-api/src/booking/usecases"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newBookingFlow(t *testing.T) (*gorm.DB, *bookingHarness) {
	t.Helper()
	dsn := fmt.Sprintf("file:booking-flow-%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Market{}, &models.Stall{}, &models.Booking{}, &models.BookingStall{}, &models.Waitlist{}, &models.Notification{}); err != nil {
		t.Fatal(err)
	}
	users := []models.User{{ID: 101, Username: "tenant-a", Password: "x", FullName: "ผู้เช่า A", Phone: "0800000001", Role: "tenant", Status: "A"}, {ID: 102, Username: "tenant-b", Password: "x", FullName: "ผู้เช่า B", Phone: "0800000002", Role: "tenant", Status: "A"}, {ID: 103, Username: "tenant-c", Password: "x", FullName: "ผู้เช่า C", Phone: "0800000003", Role: "tenant", Status: "A"}}
	if err := db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	return db, &bookingHarness{uc: bookinguc.NewUseCase(bookingrepo.NewRepository(db))}
}

type bookingHarness struct {
	uc interface {
		CreateBooking(int, requests.CreateBookingRequest) (*responses.BookingResponse, error)
		CreateWaitlist(int, requests.CreateWaitlistRequest) (*responses.WaitlistResponse, error)
		CancelBookingStall(int, string, int, string, string) (*responses.BookingResponse, error)
		WithdrawWaitlist(int, int, string) error
	}
}

func baseBooking(stalls ...string) requests.CreateBookingRequest {
	items := make([]requests.BookingStallRequest, 0, len(stalls))
	for _, ident := range stalls {
		items = append(items, requests.BookingStallRequest{StallIdent: ident, MonthlyRent: 2000, DepositFee: 2000})
	}
	return requests.CreateBookingRequest{Stalls: items, MarketName: "ตลาดทดสอบ", ApplicantName: "ผู้เช่าทดสอบ", ApplicantPhone: "0800000000", VisitDate: "2026-11-01", VisitTime: "09:00 - 10:00"}
}

func TestNormalMultiStallBookingAndSingleStallCancellation(t *testing.T) {
	db, h := newBookingFlow(t)
	for _, ident := range []string{"A-01", "A-02"} {
		if err := db.Create(&models.Stall{Ident: ident, MarketID: 1, MarketName: "ตลาดทดสอบ", MonthlyRent: 2000, StallStatus: "available", Status: "A"}).Error; err != nil {
			t.Fatal(err)
		}
	}
	created, err := h.uc.CreateBooking(101, baseBooking("A-01", "A-02"))
	if err != nil {
		t.Fatal(err)
	}
	if created.StallCount != 2 || len(created.SubBookings) != 2 {
		t.Fatalf("expected one card with 2 stalls, got count=%d children=%d", created.StallCount, len(created.SubBookings))
	}
	var rows int64
	db.Model(&models.Booking{}).Count(&rows)
	if rows != 2 {
		t.Fatalf("expected 2 booking rows, got %d", rows)
	}
	var links int64
	db.Model(&models.BookingStall{}).Count(&links)
	if links != 2 {
		t.Fatalf("expected 2 booking_stall rows, got %d", links)
	}
	_, err = h.uc.CancelBookingStall(created.ID, "A-01", 101, "tenant", "เปลี่ยนใจ")
	if err != nil {
		t.Fatal(err)
	}
	var a1, a2 models.Booking
	db.Where("stall_ident = ?", "A-01").First(&a1)
	db.Where("stall_ident = ?", "A-02").First(&a2)
	if a1.BookingStatus != "cancelled" || a2.BookingStatus == "cancelled" {
		t.Fatalf("single-stall cancellation leaked: A-01=%s A-02=%s", a1.BookingStatus, a2.BookingStatus)
	}
}

func TestWaitlistOfferProtectionAndAcceptance(t *testing.T) {
	db, h := newBookingFlow(t)
	db.Create(&models.Stall{Ident: "W-01", MarketID: 1, MarketName: "ตลาดทดสอบ", MonthlyRent: 2000, StallStatus: "available", Status: "A"})
	holder, err := h.uc.CreateBooking(101, baseBooking("W-01"))
	if err != nil {
		t.Fatal(err)
	}
	w2, err := h.uc.CreateWaitlist(102, requests.CreateWaitlistRequest{StallIdent: "W-01", MarketName: "ตลาดทดสอบ"})
	if err != nil {
		t.Fatal(err)
	}
	w3, err := h.uc.CreateWaitlist(103, requests.CreateWaitlistRequest{StallIdent: "W-01", MarketName: "ตลาดทดสอบ"})
	if err != nil {
		t.Fatal(err)
	}

	// The original booking remains the sole primary holder. Waiting list
	// registrations must never rewrite it into a legacy queue booking.
	var persistedHolder models.Booking
	if err := db.First(&persistedHolder, holder.ID).Error; err != nil {
		t.Fatalf("reload primary holder: %v", err)
	}
	if persistedHolder.QueueNumber != 0 || persistedHolder.BookingStatus != "pending_approval" {
		t.Fatalf("primary holder was changed into a queue: status=%s queue=%d", persistedHolder.BookingStatus, persistedHolder.QueueNumber)
	}
	var activeBookingCount int64
	if err := db.Model(&models.Booking{}).
		Where("stall_ident = ? AND status = 'A'", "W-01").
		Count(&activeBookingCount).Error; err != nil {
		t.Fatalf("count primary bookings: %v", err)
	}
	if activeBookingCount != 1 {
		t.Fatalf("active booking count = %d, want exactly one primary holder", activeBookingCount)
	}
	if w2.Position != 1 || w3.Position != 2 {
		t.Fatalf("wrong positions: %d,%d", w2.Position, w3.Position)
	}
	if _, err := h.uc.CreateWaitlist(102, requests.CreateWaitlistRequest{StallIdent: "W-01"}); err == nil {
		t.Fatal("duplicate waitlist should fail")
	}
	if _, err := h.uc.CancelBookingStall(holder.ID, "W-01", 101, "tenant", "ยกเลิก"); err != nil {
		t.Fatal(err)
	}
	var offer models.Waitlist
	if err := db.Where("stall_ident = ? AND status = ?", "W-01", "offered").First(&offer).Error; err != nil {
		t.Fatal(err)
	}
	if offer.TenantID != 102 || offer.ExpiresAt == nil {
		t.Fatalf("first waiter was not offered: %+v", offer)
	}
	if _, err := h.uc.CreateBooking(103, baseBooking("W-01")); err == nil {
		t.Fatal("non-offered tenant bypassed queue")
	}
	if _, err := h.uc.CreateBooking(102, baseBooking("W-01")); err != nil {
		t.Fatalf("offered tenant could not book: %v", err)
	}
	db.First(&offer, offer.ID)
	if offer.Status != "accepted" {
		t.Fatalf("expected accepted, got %s", offer.Status)
	}
}

func TestOfferedTenantWithdrawPromotesNext(t *testing.T) {
	db, h := newBookingFlow(t)
	db.Create(&models.Stall{Ident: "W-02", MarketID: 1, MarketName: "ตลาดทดสอบ", MonthlyRent: 2000, StallStatus: "available", Status: "A"})
	holder, _ := h.uc.CreateBooking(101, baseBooking("W-02"))
	first, _ := h.uc.CreateWaitlist(102, requests.CreateWaitlistRequest{StallIdent: "W-02"})
	h.uc.CreateWaitlist(103, requests.CreateWaitlistRequest{StallIdent: "W-02"})
	h.uc.CancelBookingStall(holder.ID, "W-02", 101, "tenant", "ยกเลิก")
	if err := h.uc.WithdrawWaitlist(first.ID, 102, "tenant"); err != nil {
		t.Fatal(err)
	}
	var next models.Waitlist
	if err := db.Where("stall_ident = ? AND tenant_id = ?", "W-02", 103).First(&next).Error; err != nil {
		t.Fatal(err)
	}
	if next.Status != "offered" || next.Position != 1 {
		t.Fatalf("next waiter not promoted: status=%s position=%d", next.Status, next.Position)
	}
}

func TestWaitlistRejectsAvailableStall(t *testing.T) {
	db, h := newBookingFlow(t)
	db.Create(&models.Stall{Ident: "FREE-01", MarketID: 1, MarketName: "ตลาดทดสอบ", StallStatus: "available", Status: "A"})
	if _, err := h.uc.CreateWaitlist(102, requests.CreateWaitlistRequest{StallIdent: "FREE-01"}); err == nil {
		t.Fatal("available stall must not accept a waitlist entry")
	}
}

func TestTenantCannotCancelAnotherTenantsBooking(t *testing.T) {
	db, h := newBookingFlow(t)
	db.Create(&models.Stall{Ident: "SEC-01", MarketID: 1, MarketName: "ตลาดทดสอบ", StallStatus: "available", Status: "A"})
	created, err := h.uc.CreateBooking(101, baseBooking("SEC-01"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h.uc.CancelBookingStall(created.ID, "SEC-01", 103, "tenant", "พยายามยกเลิกแทน"); err == nil {
		t.Fatal("tenant must not cancel another tenant's booking")
	}
}

func TestMultiStallPrevalidationPreventsPartialBooking(t *testing.T) {
	db, h := newBookingFlow(t)
	for _, ident := range []string{"SAFE-01", "HELD-01"} {
		db.Create(&models.Stall{Ident: ident, MarketID: 1, MarketName: "ตลาดทดสอบ", StallStatus: "available", Status: "A"})
	}
	if _, err := h.uc.CreateBooking(101, baseBooking("HELD-01")); err != nil {
		t.Fatal(err)
	}
	if _, err := h.uc.CreateBooking(102, baseBooking("SAFE-01", "HELD-01")); err == nil {
		t.Fatal("mixed available/held request must fail")
	}
	var count int64
	db.Model(&models.Booking{}).Where("stall_ident = ? AND created_by = ?", "SAFE-01", 102).Count(&count)
	if count != 0 {
		t.Fatalf("partial booking was written for SAFE-01: %d", count)
	}
}

func TestExpiredOfferAdvancesToNextWaiter(t *testing.T) {
	db, _ := newBookingFlow(t)
	repo := bookingrepo.NewRepository(db)
	past := time.Now().Add(-time.Minute)
	db.Create(&models.Stall{Ident: "EXP-01", MarketID: 1, MarketName: "ตลาดทดสอบ", StallStatus: "offered", Status: "A"})
	db.Create(&models.Waitlist{StallIdent: "EXP-01", TenantID: 102, TenantName: "ผู้เช่า B", Position: 1, Status: "offered", ExpiresAt: &past})
	db.Create(&models.Waitlist{StallIdent: "EXP-01", TenantID: 103, TenantName: "ผู้เช่า C", Position: 2, Status: "waiting"})
	if err := repo.ProcessWaitlistLifecycle(); err != nil {
		t.Fatal(err)
	}
	var expired, next models.Waitlist
	db.Where("stall_ident = ? AND tenant_id = ?", "EXP-01", 102).First(&expired)
	db.Where("stall_ident = ? AND tenant_id = ?", "EXP-01", 103).First(&next)
	if expired.Status != "expired" {
		t.Fatalf("expected expired, got %s", expired.Status)
	}
	if next.Status != "offered" || next.Position != 1 || next.ExpiresAt == nil {
		t.Fatalf("next waiter was not offered: %+v", next)
	}
}

func TestLifecycleClosesStaleWaitlistAfterPrimaryLeaseApplication(t *testing.T) {
	db, _ := newBookingFlow(t)
	repo := bookingrepo.NewRepository(db)
	tenantID := 101
	holder := models.Booking{
		BookingNo:          "BK-STALE-WAITLIST",
		StallIdent:         "TWR-03-026",
		ApplicantName:      "ผู้เช่าหลัก",
		ApplicantPhone:     "0800000000",
		VisitDate:          "2026-10-26",
		VisitTime:          "09:00",
		BookingStatus:      "lease_applied",
		AppointmentGroupID: "BKG-STALE-WAITLIST",
		Status:             "A",
		CreatedBy:          &tenantID,
	}
	if err := db.Create(&holder).Error; err != nil {
		t.Fatalf("create lease-applied holder: %v", err)
	}
	waiter := models.Waitlist{
		StallIdent: "TWR-03-026",
		TenantID:   102,
		TenantName: "ผู้รอคิวเดิม",
		Position:   1,
		Status:     "waiting",
	}
	if err := db.Create(&waiter).Error; err != nil {
		t.Fatalf("create stale waitlist: %v", err)
	}

	if err := repo.ProcessWaitlistLifecycle(); err != nil {
		t.Fatalf("process lifecycle: %v", err)
	}
	if err := db.First(&waiter, waiter.ID).Error; err != nil {
		t.Fatalf("reload waitlist: %v", err)
	}
	if waiter.Status != "cancelled" {
		t.Fatalf("waitlist status = %q, want cancelled", waiter.Status)
	}

	var notificationCount int64
	if err := db.Model(&models.Notification{}).
		Where("recipient_user_id = ? AND stall_number = ? AND status_code = ?", 102, "TWR-03-026", "waitlist_cancelled").
		Count(&notificationCount).Error; err != nil {
		t.Fatalf("count reconciliation notification: %v", err)
	}
	if notificationCount != 1 {
		t.Fatalf("notification count = %d, want 1", notificationCount)
	}
}

func TestThreePeopleBookingCancelPromotesSecondAndMovesThirdUp(t *testing.T) {
	db, h := newBookingFlow(t)
	db.Create(&models.Stall{Ident: "THREE-01", MarketID: 1, MarketName: "ตลาดทดสอบ", MonthlyRent: 2500, StallStatus: "available", Status: "A"})

	// 1. Person 1 books the stall
	holder, err := h.uc.CreateBooking(101, baseBooking("THREE-01"))
	if err != nil {
		t.Fatalf("Person 1 create booking: %v", err)
	}

	// 2. Person 2 joins waitlist -> Position 1
	w2, err := h.uc.CreateWaitlist(102, requests.CreateWaitlistRequest{StallIdent: "THREE-01"})
	if err != nil {
		t.Fatalf("Person 2 join waitlist: %v", err)
	}
	if w2.Position != 1 || w2.Status != "waiting" {
		t.Fatalf("Person 2 initial waitlist: status=%s pos=%d, want waiting/1", w2.Status, w2.Position)
	}

	// 3. Person 3 joins waitlist -> Position 2
	w3, err := h.uc.CreateWaitlist(103, requests.CreateWaitlistRequest{StallIdent: "THREE-01"})
	if err != nil {
		t.Fatalf("Person 3 join waitlist: %v", err)
	}
	if w3.Position != 2 || w3.Status != "waiting" {
		t.Fatalf("Person 3 initial waitlist: status=%s pos=%d, want waiting/2", w3.Status, w3.Position)
	}

	// 4. Person 1 cancels booking
	_, err = h.uc.CancelBookingStall(holder.ID, "THREE-01", 101, "tenant", "ขอยกเลิกการจอง")
	if err != nil {
		t.Fatalf("Person 1 cancel booking: %v", err)
	}

	// 5. Verify Person 2 is now offered (ถึงสิทธิ์)
	var reloadedW2 models.Waitlist
	if err := db.First(&reloadedW2, w2.ID).Error; err != nil {
		t.Fatalf("reload Person 2: %v", err)
	}
	if reloadedW2.Status != "offered" || reloadedW2.Position != 1 || reloadedW2.ExpiresAt == nil {
		t.Fatalf("Person 2 after cancel: status=%s pos=%d expires=%v, want offered/1", reloadedW2.Status, reloadedW2.Position, reloadedW2.ExpiresAt)
	}

	// 6. Verify Person 3 queue has moved up to position 1 (คิวเลื่อนขึ้นมาเป็นคิวที่ 1)
	var reloadedW3 models.Waitlist
	if err := db.First(&reloadedW3, w3.ID).Error; err != nil {
		t.Fatalf("reload Person 3: %v", err)
	}
	if reloadedW3.Status != "waiting" || reloadedW3.Position != 1 {
		t.Fatalf("Person 3 after cancel: status=%s pos=%d, want waiting/1", reloadedW3.Status, reloadedW3.Position)
	}
}

