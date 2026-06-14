package notifications_test

import (
	"testing"

	"claimsplatform/internal/dimensions/notifications"
	"claimsplatform/internal/domain"
)

func TestFiresSubmittedAndOutcomeEvents(t *testing.T) {
	cfg := domain.ConfigDocument{Notifications: domain.NotificationsConfig{
		Channels: []string{"email", "sms"},
		Events: map[string][]string{
			"claim_submitted": {"email", "sms"},
			"claim_routed":    {"email"},
		}}}
	dec := &domain.ClaimDecision{Accepted: true, Approval: &domain.ApprovalDecision{Outcome: domain.Routed}}
	notifications.New().Evaluate(cfg, domain.Claim{}, dec)

	got := map[string][]string{}
	for _, n := range dec.Notifications {
		got[n.Event] = n.Channels
	}
	if len(got["claim_submitted"]) != 2 {
		t.Fatalf("expected submitted on email+sms, got %v", got["claim_submitted"])
	}
	if len(got["claim_routed"]) != 1 || got["claim_routed"][0] != "email" {
		t.Fatalf("expected routed on email only, got %v", got["claim_routed"])
	}
	if _, ok := got["claim_auto_approved"]; ok {
		t.Fatal("should not fire auto_approved when routed")
	}
}

func TestValidateWebhookRequiresURL(t *testing.T) {
	c := domain.ConfigDocument{Notifications: domain.NotificationsConfig{Channels: []string{"email", "webhook"}}}
	if errs := notifications.New().Validate(c); len(errs) == 0 {
		t.Fatal("expected error: webhook channel without webhookUrl")
	}
}

func TestFallbackToDefaultChannelsWhenEventUnmapped(t *testing.T) {
	c := domain.ConfigDocument{Notifications: domain.NotificationsConfig{
		Channels: []string{"email", "sms"},
		Events:   map[string][]string{},
	}}
	dec := &domain.ClaimDecision{Accepted: true, Approval: &domain.ApprovalDecision{Outcome: domain.AutoApproved}}
	notifications.New().Evaluate(c, domain.Claim{}, dec)
	if len(dec.Notifications) == 0 {
		t.Fatal("expected notifications to fire")
	}
	for _, n := range dec.Notifications {
		if len(n.Channels) != 2 {
			t.Fatalf("event %s: expected fallback to 2 default channels, got %v", n.Event, n.Channels)
		}
	}
}
