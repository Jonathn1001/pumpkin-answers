package sla

import (
	"fmt"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

type dimension struct{}

func New() registry.Dimension { return dimension{} }

func init() { registry.Register(New()) }

func (dimension) Key() string { return registry.KeySLA }

func (dimension) Evaluate(cfg domain.ConfigDocument, claim domain.Claim, dec *domain.ClaimDecision) {
	if !dec.Accepted {
		return
	}
	days := cfg.SLA.DefaultDays
	if override, ok := cfg.SLA.PerClaimType[claim.Type]; ok {
		days = override
	}
	deadline := claim.SubmittedAt.AddDate(0, 0, days)
	dec.SLADays = days
	dec.SLADeadline = &deadline
	esc := cfg.SLA.Escalation
	dec.Escalation = &esc
	dec.Trace = append(dec.Trace, domain.TraceEntry{Dimension: "sla",
		Explanation: fmt.Sprintf("SLA %d days -> deadline %s", days, deadline.Format("2006-01-02"))})
}

func (dimension) Validate(cfg domain.ConfigDocument) []domain.FieldError {
	var errs []domain.FieldError
	if cfg.SLA.DefaultDays < 1 {
		errs = append(errs, domain.FieldError{Field: "sla.defaultDays", Message: "defaultDays must be >= 1"})
	}
	return errs
}

func (dimension) DefaultConfig() any {
	return domain.SLAConfig{DefaultDays: 7, PerClaimType: map[domain.ClaimType]int{},
		Escalation: domain.Escalation{WarnBeforeDays: 2, NotifyRole: "supervisor"}}
}

func (dimension) UISchema() []registry.FieldDescriptor {
	return []registry.FieldDescriptor{
		{Key: "sla.defaultDays", Label: "Default SLA (days)", Type: "number", Widget: "number", Required: true},
		{Key: "sla.perClaimType", Label: "Per-type overrides", Type: "object", Widget: "claimtype-number-map"},
		{Key: "sla.escalation", Label: "Escalation", Type: "object", Widget: "escalation"},
	}
}
