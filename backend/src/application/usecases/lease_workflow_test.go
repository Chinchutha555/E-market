package usecases

import (
	"errors"
	"testing"
)

func TestLeaseTransitions(t *testing.T) {
	tests := []struct {
		state, action, role, note string
		clear                     bool
		want                      string
		wantErr                   error
	}{
		{"pending_officer", "approve", "staff", "", false, "pending_market_director", nil},
		{"pending_officer", "return_revision", "staff", "เอกสารไม่ครบ", false, "revision_required", nil},
		{"pending_officer", "reject", "staff", "ไม่ผ่านเกณฑ์", false, "rejected", nil},
		{"pending_officer", "reject", "staff", " ", false, "", errLeaseReasonRequired},
		{"revision_required", "resubmit", "tenant", "", false, "pending_officer", nil},
		{"pending_market_director", "approve", "market_director", "", false, "", errLeaseDebtNotVerified},
		{"pending_market_director", "approve", "market_director", "", true, "pending_bureau_director", nil},
		{"pending_bureau_director", "approve", "bureau_director", "", true, "approved", nil},
		{"pending_bureau_director", "approve", "staff", "", true, "", errInvalidLeaseTransition},
		{"approved", "approve", "bureau_director", "", true, "", errInvalidLeaseTransition},
	}
	for _, tt := range tests {
		got, err := leaseTransition(tt.state, tt.action, tt.role, tt.note, tt.clear)
		if got != tt.want || !errors.Is(err, tt.wantErr) {
			t.Errorf("%s/%s/%s = %q, %v; want %q, %v", tt.state, tt.action, tt.role, got, err, tt.want, tt.wantErr)
		}
	}
}
