package configvalidation_test

import (
	"testing"

	"claimsplatform/internal/configvalidation"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/seed"
)

func TestValidSeedConfigPasses(t *testing.T) {
	if errs := configvalidation.Validate(seed.SafeGuard()); len(errs) != 0 {
		t.Fatalf("expected valid seed config, got %v", errs)
	}
}

func TestBrokenConfigProducesFieldErrors(t *testing.T) {
	cfg := seed.SafeGuard()
	cfg.Notifications.Channels = []string{"webhook"} // webhook channel but no URL
	cfg.Notifications.WebhookURL = ""
	cfg.ClaimTypes = map[domain.ClaimType]domain.ClaimTypeConfig{domain.Outpatient: {Enabled: false}} // zero enabled
	errs := configvalidation.Validate(cfg)
	want := map[string]bool{"notifications.webhookUrl": false, "claimTypes": false}
	for _, e := range errs {
		if _, ok := want[e.Field]; ok {
			want[e.Field] = true
		}
	}
	for field, seen := range want {
		if !seen {
			t.Fatalf("expected field error %q, got %v", field, errs)
		}
	}
}
