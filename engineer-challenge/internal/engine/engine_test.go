package engine_test

import (
	"testing"
	"time"

	_ "claimsplatform/internal/dimensions" // register the six dimensions
	"claimsplatform/internal/domain"
	"claimsplatform/internal/engine"
	"claimsplatform/internal/seed"
)

func canonicalClaim() domain.Claim {
	return domain.Claim{Type: domain.Outpatient, Amount: 10000,
		SubmittedAt:  time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC),
		CustomFields: map[string]any{"employeeId": "EMP1234", "policyNumber": "HF-12345678", "memberTier": "Gold", "nationalId": "123456789012", "citizenCategory": "General"}}
}

func TestThreeTenantsDifferOnIdenticalClaim(t *testing.T) {
	cases := []struct {
		name                  string
		cfg                   domain.ConfigDocument
		wantOutcome           domain.ApprovalOutcome
		wantSLADays           int
		wantSubmittedChannels int
	}{
		{"SafeGuard", seed.SafeGuard(), domain.AutoApproved, 5, 1},
		{"HealthFirst", seed.HealthFirst(), domain.Routed, 7, 2},
		{"GovHealth", seed.GovHealth(), domain.Routed, 15, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dec := engine.ProcessClaim(tc.cfg, canonicalClaim())
			if !dec.Accepted {
				t.Fatalf("%s: expected accepted", tc.name)
			}
			if dec.Approval.Outcome != tc.wantOutcome {
				t.Fatalf("%s: outcome = %s, want %s", tc.name, dec.Approval.Outcome, tc.wantOutcome)
			}
			if dec.SLADays != tc.wantSLADays {
				t.Fatalf("%s: slaDays = %d, want %d", tc.name, dec.SLADays, tc.wantSLADays)
			}
			var submitted []string
			for _, n := range dec.Notifications {
				if n.Event == "claim_submitted" {
					submitted = n.Channels
				}
			}
			if len(submitted) != tc.wantSubmittedChannels {
				t.Fatalf("%s: submitted channels = %v, want %d", tc.name, submitted, tc.wantSubmittedChannels)
			}
		})
	}
}

func TestGovHealthRoutesToCommittee(t *testing.T) {
	dec := engine.ProcessClaim(seed.GovHealth(), canonicalClaim())
	if dec.Approval.Route == nil || dec.Approval.Route.CommitteeName != "Government Claims Committee" {
		t.Fatalf("expected committee route, got %+v", dec.Approval.Route)
	}
}

func TestDisabledTypeRejectedForGovHealth(t *testing.T) {
	claim := canonicalClaim()
	claim.Type = domain.Dental // GovHealth disables DENTAL
	dec := engine.ProcessClaim(seed.GovHealth(), claim)
	if dec.Accepted {
		t.Fatal("expected DENTAL rejected for GovHealth")
	}
}
