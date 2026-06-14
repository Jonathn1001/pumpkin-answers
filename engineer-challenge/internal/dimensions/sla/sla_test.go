package sla_test

import (
	"testing"
	"time"

	"claimsplatform/internal/dimensions/sla"
	"claimsplatform/internal/domain"
)

func TestPerTypeOverrideWins(t *testing.T) {
	cfg := domain.ConfigDocument{SLA: domain.SLAConfig{
		DefaultDays: 7, PerClaimType: map[domain.ClaimType]int{domain.Outpatient: 5},
		Escalation: domain.Escalation{WarnBeforeDays: 2, NotifyRole: "supervisor"}}}
	submitted := time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC)
	dec := &domain.ClaimDecision{Accepted: true}
	sla.New().Evaluate(cfg, domain.Claim{Type: domain.Outpatient, SubmittedAt: submitted}, dec)
	if dec.SLADays != 5 {
		t.Fatalf("expected 5-day override, got %d", dec.SLADays)
	}
	want := time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC)
	if dec.SLADeadline == nil || !dec.SLADeadline.Equal(want) {
		t.Fatalf("expected deadline %v, got %v", want, dec.SLADeadline)
	}
}

func TestDefaultWhenNoOverride(t *testing.T) {
	cfg := domain.ConfigDocument{SLA: domain.SLAConfig{DefaultDays: 15}}
	dec := &domain.ClaimDecision{Accepted: true}
	sla.New().Evaluate(cfg, domain.Claim{Type: domain.Inpatient, SubmittedAt: time.Now()}, dec)
	if dec.SLADays != 15 {
		t.Fatalf("expected default 15, got %d", dec.SLADays)
	}
}
