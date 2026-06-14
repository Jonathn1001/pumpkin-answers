package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/seed"
)

// ProcessClaim/PreviewClaim must reject structurally-invalid claims with a
// ValidationError (HTTP 422) before touching config or the engine.
func TestProcessClaimRejectsInvalidClaim(t *testing.T) {
	svc := seeded(t)
	good := canonicalClaim()

	cases := map[string]domain.Claim{
		"negative amount":  withClaim(good, func(c *domain.Claim) { c.Amount = -1 }),
		"empty type":       withClaim(good, func(c *domain.Claim) { c.Type = "" }),
		"unknown type":     withClaim(good, func(c *domain.Claim) { c.Type = "BOGUS" }),
		"zero submittedAt": withClaim(good, func(c *domain.Claim) { c.SubmittedAt = time.Time{} }),
	}
	for name, claim := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := svc.ProcessClaim(context.Background(), "sg", claim)
			var ve domain.ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("expected ValidationError, got %v", err)
			}
		})
	}
}

func TestProcessClaimAcceptsValidClaim(t *testing.T) {
	if _, err := seeded(t).ProcessClaim(context.Background(), "sg", canonicalClaim()); err != nil {
		t.Fatalf("valid claim was rejected: %v", err)
	}
}

func TestPreviewClaimRejectsInvalidClaim(t *testing.T) {
	svc := seeded(t)
	gov := seed.GovHealth()
	bad := withClaim(canonicalClaim(), func(c *domain.Claim) { c.Amount = -1 })
	_, err := svc.PreviewClaim(context.Background(), "sg", bad, nil, &gov)
	var ve domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func withClaim(base domain.Claim, mutate func(*domain.Claim)) domain.Claim {
	mutate(&base)
	return base
}
