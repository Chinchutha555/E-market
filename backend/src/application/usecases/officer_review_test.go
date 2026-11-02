package usecases

import (
	"emarket-api/models"
	"emarket-api/requests"
	"encoding/json"
	"testing"
)

type officerTestRepo struct {
	app         models.Application
	stallStatus string
}

func (r *officerTestRepo) FindAll(requests.ApplicationFilter) ([]models.Application, error) {
	return nil, nil
}
func (r *officerTestRepo) FindByID(int) (*models.Application, error) {
	copy := r.app
	return &copy, nil
}
func (r *officerTestRepo) Create(*models.Application) error   { return nil }
func (r *officerTestRepo) Update(a *models.Application) error { r.app = *a; return nil }
func (r *officerTestRepo) UpdateStallStatus(_ string, status string) error {
	r.stallStatus = status
	return nil
}
func TestOfficerReviewActions(t *testing.T) {
	tests := []struct {
		action, note, wantState, wantStall string
		wantError                          bool
	}{
		{"approve", "ครบถ้วน", "pending_market_director", "", false},
		{"return_revision", "เอกสารไม่ครบ", "revision_required", "", false},
		{"reject", "ไม่ผ่านเกณฑ์", "rejected", "available", false},
		{"return_revision", "  ", "pending_officer", "", true},
		{"reject", "", "pending_officer", "", true},
		{"unknown", "", "pending_officer", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.action+tt.note, func(t *testing.T) {
			repo := &officerTestRepo{app: models.Application{ID: 1, StallIdent: "A-1", ApplicationState: "pending_officer"}}
			_, err := NewUseCase(repo).Action(1, 1, requests.ApplicationActionRequest{Action: tt.action, Comment: tt.note, OfficerReview: &models.OfficerReview{StallStatus: "none"}, SignatureBase64: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="})
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v", err)
			}
			if repo.app.ApplicationState != tt.wantState || repo.stallStatus != tt.wantStall {
				t.Fatalf("state=%s stall=%s", repo.app.ApplicationState, repo.stallStatus)
			}
			if !tt.wantError && repo.app.OfficerComment != tt.note {
				t.Fatalf("officer comment = %q", repo.app.OfficerComment)
			}
		})
	}
}

func TestOfficerApprovalRequiresSignature(t *testing.T) {
	repo := &officerTestRepo{app: models.Application{ID: 1, ApplicationType: "stall_lease", StallIdent: "A-1", ApplicationState: "pending_officer"}}
	uc := NewUseCase(repo)
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "approve"}); err == nil {
		t.Fatal("approval without signature succeeded")
	}
	if repo.app.ApplicationState != "pending_officer" {
		t.Fatal("missing signature changed status")
	}
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "approve", OfficerReview: &models.OfficerReview{StallStatus: "none"}, SignatureBase64: "dGVzdA=="}); err == nil {
		t.Fatal("approval with invalid image succeeded")
	}

	signature := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "approve", SignatureBase64: signature}); err == nil {
		t.Fatal("approval without officer review succeeded")
	}
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "approve", OfficerReview: &models.OfficerReview{StallStatus: "none"}, SignatureBase64: signature}); err != nil {
		t.Fatal(err)
	}
	if repo.app.OfficerReviewJSON == "" {
		t.Fatal("officer review was not saved")
	}
	if repo.app.ApplicationState != "pending_market_director" || len(repo.app.OfficerSignature) == 0 || repo.app.OfficerSignatureType != "image/png" {
		t.Fatal("signature or approval was not saved")
	}
}

func TestOfficerReviewRoutesByStallResult(t *testing.T) {
	signature := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	tests := []struct {
		name   string
		review *models.OfficerReview
		state  string
	}{
		{"no stall skips finance", &models.OfficerReview{StallStatus: "none"}, "pending_market_director"},
		{"has stall goes finance", &models.OfficerReview{StallStatus: "has", StallCount: 1}, "pending_debt_check"},
		{"other closes application", &models.OfficerReview{StallStatus: "other", OtherDetails: "complaint"}, "cancelled"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &officerTestRepo{app: models.Application{ID: 1, ApplicationType: "stall_lease", StallIdent: "A-1", ApplicationState: "pending_officer"}}
			_, err := NewUseCase(repo).Action(1, 1, requests.ApplicationActionRequest{Action: "approve", OfficerReview: tt.review, SignatureBase64: signature})
			if err != nil {
				t.Fatal(err)
			}
			if repo.app.ApplicationState != tt.state {
				t.Fatalf("state=%s want=%s", repo.app.ApplicationState, tt.state)
			}
		})
	}
}

func TestFinanceDebtReviewRoutesToMarketDirectorOnlyWhenClear(t *testing.T) {
	signature := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	for _, tt := range []struct {
		debt  *models.MarketDirectorReview
		state string
	}{
		{&models.MarketDirectorReview{DebtStatus: "none"}, "pending_market_director"},
		{&models.MarketDirectorReview{DebtStatus: "has", DebtAmount: 500}, "debt_pending"},
	} {
		repo := &officerTestRepo{app: models.Application{ID: 1, ApplicationType: "stall_lease", StallIdent: "A-1", ApplicationState: "pending_debt_check"}}
		_, err := NewUseCase(repo).Action(1, 2, requests.ApplicationActionRequest{Action: "approve", MarketDirectorReview: tt.debt, SignatureBase64: signature})
		if err != nil {
			t.Fatal(err)
		}
		if repo.app.ApplicationState != tt.state || len(repo.app.FinanceSignature) == 0 {
			t.Fatalf("state=%s finance_signature=%d", repo.app.ApplicationState, len(repo.app.FinanceSignature))
		}
	}
}

func TestOfficerReviewValidationAndStorage(t *testing.T) {
	if _, err := encodeOfficerReview(nil); err == nil {
		t.Fatal("accepted missing review")
	}
	if _, err := encodeOfficerReview(&models.OfficerReview{StallStatus: "has"}); err == nil {
		t.Fatal("accepted missing stall count")
	}
	if _, err := encodeOfficerReview(&models.OfficerReview{StallStatus: "none", OtherChecked: true}); err == nil {
		t.Fatal("accepted empty other details")
	}
	input := &models.OfficerReview{StallStatus: "has", StallCount: 2, StallNumbers: " A-1, A-2 ", Project: "Market", TradeType: "Food", OtherContract: true, OtherChecked: true, OtherDetails: "Other", DocumentsComplete: true, BookingReference: "BK-1"}
	result, err := encodeOfficerReview(input)
	if err != nil {
		t.Fatal(err)
	}
	var decoded models.OfficerReview
	if err := json.Unmarshal([]byte(*result), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.StallCount != 2 || decoded.StallNumbers != "A-1, A-2" || !decoded.OtherContract || decoded.BookingReference != "BK-1" {
		t.Fatalf("wrong review: %#v", decoded)
	}
}

func TestDirectorReviewsAndBureauDecisions(t *testing.T) {
	if _, err := encodeMarketDirectorReview(nil); err == nil {
		t.Fatal("accepted missing debt review")
	}
	if _, err := encodeMarketDirectorReview(&models.MarketDirectorReview{DebtStatus: "has", DebtAmount: 0}); err == nil {
		t.Fatal("accepted zero debt")
	}
	market, err := encodeMarketDirectorReview(&models.MarketDirectorReview{DebtStatus: "has", DebtAmount: 123.45})
	if err != nil || market == nil {
		t.Fatalf("market review: %v", err)
	}
	repo := &officerTestRepo{app: models.Application{ID: 1, ApplicationType: "stall_lease", StallIdent: "A-1", ApplicationState: "pending_bureau_director"}}
	uc := NewUseCase(repo)
	review := &models.BureauDirectorReview{Decision: "reject"}
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "reject", BureauDirectorReview: review}); err == nil {
		t.Fatal("rejection without OTP succeeded")
	}
	if repo.app.ApplicationState != "pending_bureau_director" {
		t.Fatal("missing OTP changed state")
	}
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "return_revision", Comment: "return", OTP: "123456"}); err == nil {
		t.Fatal("bureau returned a closed decision for revision")
	}
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "reject", BureauDirectorReview: review, OTP: "123456"}); err != nil {
		t.Fatal(err)
	}
	if repo.app.ApplicationState != "rejected" || len(repo.app.BureauDirectorSignature) != 0 {
		t.Fatal("rejection did not close without signature")
	}
	var saved models.BureauDirectorReview
	if err := json.Unmarshal([]byte(repo.app.BureauDirectorReviewJSON), &saved); err != nil || saved.Decision != "reject" {
		t.Fatal("rejection review not saved")
	}
}
