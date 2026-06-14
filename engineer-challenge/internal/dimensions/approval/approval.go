package approval

import (
	"fmt"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

type dimension struct{}

func New() registry.Dimension { return dimension{} }

func init() { registry.Register(New()) }

func (dimension) Key() string { return registry.KeyApproval }

func (dimension) Evaluate(cfg domain.ConfigDocument, claim domain.Claim, dec *domain.ClaimDecision) {
	if !dec.Accepted {
		return
	}
	a := cfg.Approval
	if a.AutoApproveThreshold > 0 && claim.Amount <= a.AutoApproveThreshold {
		dec.Approval = &domain.ApprovalDecision{Outcome: domain.AutoApproved}
		dec.Trace = append(dec.Trace, domain.TraceEntry{Dimension: "approval",
			Explanation: fmt.Sprintf("amount %d <= threshold %d -> auto-approved", claim.Amount, a.AutoApproveThreshold)})
		return
	}
	route := &domain.ApprovalRoute{Model: a.Model}
	if a.Model == domain.ApprovalModelCommittee && a.Committee != nil {
		route.CommitteeName = a.Committee.Name
		route.RequiredApprovals = a.Committee.RequiredApprovals
		dec.Trace = append(dec.Trace, domain.TraceEntry{Dimension: "approval",
			Explanation: fmt.Sprintf("routed to committee %q (%d approvals)", route.CommitteeName, route.RequiredApprovals)})
	} else {
		matched := false
		for _, tier := range a.Tiers {
			if tier.MaxAmount == nil || claim.Amount <= *tier.MaxAmount {
				route.TierLabel = tier.Label
				route.ApproverRole = tier.ApproverRole
				matched = true
				break
			}
		}
		if matched {
			dec.Trace = append(dec.Trace, domain.TraceEntry{Dimension: "approval",
				Explanation: fmt.Sprintf("amount %d -> tier %q (%s)", claim.Amount, route.TierLabel, route.ApproverRole)})
		} else {
			dec.Trace = append(dec.Trace, domain.TraceEntry{Dimension: "approval",
				Explanation: fmt.Sprintf("amount %d matched no tier (check approval.tiers configuration)", claim.Amount)})
		}
	}
	dec.Approval = &domain.ApprovalDecision{Outcome: domain.Routed, Route: route}
}

func (dimension) Validate(cfg domain.ConfigDocument) []domain.FieldError {
	var errs []domain.FieldError
	a := cfg.Approval
	if a.AutoApproveThreshold < 0 {
		errs = append(errs, domain.FieldError{Field: "approval.autoApproveThreshold", Message: "must be >= 0"})
	}
	switch a.Model {
	case domain.ApprovalModelTiered:
		if len(a.Tiers) == 0 {
			errs = append(errs, domain.FieldError{Field: "approval.tiers", Message: "tiered model requires at least one tier"})
		}
		var prev int64 = -1
		openSeen := false
		for i, tier := range a.Tiers {
			if tier.MaxAmount == nil {
				openSeen = true
				continue
			}
			if openSeen {
				errs = append(errs, domain.FieldError{Field: fmt.Sprintf("approval.tiers[%d]", i), Message: "no tier may follow the open-ended tier"})
			}
			if *tier.MaxAmount <= prev {
				errs = append(errs, domain.FieldError{Field: fmt.Sprintf("approval.tiers[%d].maxAmount", i), Message: "tier ceilings must be strictly ascending"})
			}
			prev = *tier.MaxAmount
		}
		if !openSeen {
			errs = append(errs, domain.FieldError{Field: "approval.tiers", Message: "final tier must be open-ended (maxAmount null)"})
		}
	case domain.ApprovalModelCommittee:
		if a.Committee == nil || a.Committee.RequiredApprovals < 1 {
			errs = append(errs, domain.FieldError{Field: "approval.committee", Message: "committee model requires committee with requiredApprovals >= 1"})
		}
	default:
		errs = append(errs, domain.FieldError{Field: "approval.model", Message: "model must be 'tiered' or 'committee'"})
	}
	return errs
}

func (dimension) DefaultConfig() any {
	return domain.ApprovalConfig{AutoApproveThreshold: 0, Model: domain.ApprovalModelTiered,
		Tiers: []domain.ApprovalTier{{Label: "Default", MaxAmount: nil, ApproverRole: "manager"}}}
}

func (dimension) UISchema() []registry.FieldDescriptor {
	return []registry.FieldDescriptor{
		{Key: "approval.autoApproveThreshold", Label: "Auto-approve threshold", Type: "number", Widget: "number", Required: true},
		{Key: "approval.model", Label: "Approval model", Type: "select", Widget: "select", Required: true, Options: []string{"tiered", "committee"}},
		{Key: "approval.tiers", Label: "Tiers", Type: "array", Widget: "tier-list"},
		{Key: "approval.committee", Label: "Committee", Type: "object", Widget: "committee"},
	}
}
