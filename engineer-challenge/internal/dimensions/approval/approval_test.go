package approval_test

import (
	"testing"

	"claimsplatform/internal/dimensions/approval"
	"claimsplatform/internal/domain"
)

func ptr(v int64) *int64 { return &v }

func tiered() domain.ConfigDocument {
	return domain.ConfigDocument{Approval: domain.ApprovalConfig{
		AutoApproveThreshold: 20000, Model: domain.ApprovalModelTiered,
		Tiers: []domain.ApprovalTier{
			{Label: "Manager", MaxAmount: ptr(50000), ApproverRole: "claims_manager"},
			{Label: "Board", MaxAmount: nil, ApproverRole: "board"},
		}}}
}

func TestAutoApproveBelowThreshold(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	approval.New().Evaluate(tiered(), domain.Claim{Amount: 10000}, dec)
	if dec.Approval == nil || dec.Approval.Outcome != domain.AutoApproved {
		t.Fatalf("expected auto_approved, got %+v", dec.Approval)
	}
}

func TestRouteToFirstFittingTier(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	approval.New().Evaluate(tiered(), domain.Claim{Amount: 30000}, dec)
	if dec.Approval.Outcome != domain.Routed || dec.Approval.Route.TierLabel != "Manager" {
		t.Fatalf("expected routed->Manager, got %+v", dec.Approval)
	}
}

func TestZeroThresholdNeverAutoApproves(t *testing.T) {
	cfg := domain.ConfigDocument{Approval: domain.ApprovalConfig{
		AutoApproveThreshold: 0, Model: domain.ApprovalModelCommittee,
		Committee: &domain.Committee{Name: "Gov Committee", RequiredApprovals: 3}}}
	dec := &domain.ClaimDecision{Accepted: true}
	approval.New().Evaluate(cfg, domain.Claim{Amount: 1}, dec)
	if dec.Approval.Outcome != domain.Routed || dec.Approval.Route.CommitteeName != "Gov Committee" {
		t.Fatalf("expected routed->committee, got %+v", dec.Approval)
	}
}
