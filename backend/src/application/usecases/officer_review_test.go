package usecases

import (
	"emarket-api/models"
	"emarket-api/requests"
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
			_, err := NewUseCase(repo).Action(1, 1, requests.ApplicationActionRequest{Action: tt.action, Comment: tt.note, SignatureBase64: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="})
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
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "approve", SignatureBase64: "dGVzdA=="}); err == nil {
		t.Fatal("approval with invalid image succeeded")
	}

	signature := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	if _, err := uc.Action(1, 1, requests.ApplicationActionRequest{Action: "approve", SignatureBase64: signature}); err != nil {
		t.Fatal(err)
	}
	if repo.app.ApplicationState != "pending_market_director" || len(repo.app.OfficerSignature) == 0 || repo.app.OfficerSignatureType != "image/png" {
		t.Fatal("signature or approval was not saved")
	}
}
