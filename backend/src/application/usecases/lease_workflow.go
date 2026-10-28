package usecases

import (
	"errors"
	"strings"
)

var errInvalidLeaseTransition = errors.New("invalid lease application transition")
var errLeaseReasonRequired = errors.New("a reason is required")
var errLeaseDebtNotVerified = errors.New("debt must be verified as clear")

// leaseTransition describes only the requested transition. The caller must
// check identity and market access and persist it atomically with notifications.
func leaseTransition(state, action, role, comment string, debtClear bool) (string, error) {
	if (action == "return_revision" || action == "reject") && strings.TrimSpace(comment) == "" {
		return "", errLeaseReasonRequired
	}
	if role == "tenant" && action == "resubmit" && (state == "revision_required" || state == "rejected") {
		return "pending_officer", nil
	}
	required := map[string]string{"pending_officer": "staff", "pending_market_director": "market_director", "pending_bureau_director": "bureau_director"}
	if required[state] == "" || role != required[state] {
		return "", errInvalidLeaseTransition
	}
	switch action {
	case "return_revision":
		return "revision_required", nil
	case "reject":
		return "rejected", nil
	case "approve":
		switch state {
		case "pending_officer":
			return "pending_market_director", nil
		case "pending_market_director":
			if !debtClear {
				return "", errLeaseDebtNotVerified
			}
			return "pending_bureau_director", nil
		case "pending_bureau_director":
			return "approved", nil
		}
	}
	return "", errInvalidLeaseTransition
}
