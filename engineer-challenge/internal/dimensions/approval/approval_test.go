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

func TestRouteToOpenEndedTier(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	approval.New().Evaluate(tiered(), domain.Claim{Amount: 60000}, dec)
	if dec.Approval.Outcome != domain.Routed || dec.Approval.Route.TierLabel != "Board" {
		t.Fatalf("expected routed->Board (open-ended tier), got %+v", dec.Approval)
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

func TestValidateAcceptsWellFormedTiers(t *testing.T) {
	if errs := approval.New().Validate(tiered()); len(errs) != 0 {
		t.Fatalf("expected valid tiered config, got %v", errs)
	}
}

func TestValidateRejectsNonAscendingTiers(t *testing.T) {
	c := domain.ConfigDocument{Approval: domain.ApprovalConfig{Model: domain.ApprovalModelTiered, Tiers: []domain.ApprovalTier{
		{Label: "A", MaxAmount: ptr(50000)},
		{Label: "B", MaxAmount: ptr(50000)},
		{Label: "C", MaxAmount: nil},
	}}}
	if errs := approval.New().Validate(c); len(errs) == 0 {
		t.Fatal("expected error: tier ceilings must be strictly ascending")
	}
}

func TestValidateRequiresOpenEndedFinalTier(t *testing.T) {
	c := domain.ConfigDocument{Approval: domain.ApprovalConfig{Model: domain.ApprovalModelTiered, Tiers: []domain.ApprovalTier{
		{Label: "A", MaxAmount: ptr(50000)},
	}}}
	if errs := approval.New().Validate(c); len(errs) == 0 {
		t.Fatal("expected error: final tier must be open-ended")
	}
}

func TestValidateCommitteeRequiresApprovals(t *testing.T) {
	c := domain.ConfigDocument{Approval: domain.ApprovalConfig{Model: domain.ApprovalModelCommittee, Committee: &domain.Committee{Name: "X", RequiredApprovals: 0}}}
	if errs := approval.New().Validate(c); len(errs) == 0 {
		t.Fatal("expected error: committee requiredApprovals >= 1")
	}
}

func TestValidateRejectsUnknownModel(t *testing.T) {
	c := domain.ConfigDocument{Approval: domain.ApprovalConfig{Model: "bogus"}}
	if errs := approval.New().Validate(c); len(errs) == 0 {
		t.Fatal("expected error: unknown approval model")
	}
}

func TestNegativeThresholdFails(t *testing.T) {
	cfg := domain.ConfigDocument{Approval: domain.ApprovalConfig{
		AutoApproveThreshold: -1, Model: "tiered",
		Tiers: []domain.ApprovalTier{{Label: "Open", MaxAmount: nil, ApproverRole: "manager"}}}}
	errs := approval.New().Validate(cfg)
	found := false
	for _, e := range errs {
		if e.Field == "approval.autoApproveThreshold" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected autoApproveThreshold error for negative value, got %v", errs)
	}
}
