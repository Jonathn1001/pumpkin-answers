package usecase_test

import (
	"context"
	"testing"
	"time"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/seed"
	"claimsplatform/internal/usecase"
)

func canonicalClaim() domain.Claim {
	return domain.Claim{
		Type: domain.Outpatient, Amount: 10000,
		SubmittedAt: time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC),
		CustomFields: map[string]any{
			"employeeId": "EMP1234", "policyNumber": "HF-12345678", "memberTier": "Gold",
			"nationalId": "123456789012", "citizenCategory": "General",
		},
	}
}

func seeded(t *testing.T) *usecase.Service {
	t.Helper()
	repo := memory.New()
	_, _ = repo.CreateTenant(context.Background(),
		domain.Tenant{Slug: "sg", Name: "SafeGuard", Status: domain.TenantActive}, seed.SafeGuard())
	return usecase.New(repo)
}

func TestProcessClaimUsesActiveConfig(t *testing.T) {
	dec, err := seeded(t).ProcessClaim(context.Background(), "sg", canonicalClaim())
	if err != nil {
		t.Fatal(err)
	}
	if dec.Approval == nil || dec.Approval.Outcome != domain.AutoApproved {
		t.Fatalf("SafeGuard should auto-approve, got %+v", dec.Approval)
	}
}

func TestPreviewWithInlineConfigDoesNotPersist(t *testing.T) {
	svc := seeded(t)
	ctx := context.Background()
	gov := seed.GovHealth()
	dec, err := svc.PreviewClaim(ctx, "sg", canonicalClaim(), nil, &gov)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Approval.Route == nil || dec.Approval.Route.CommitteeName == "" {
		t.Fatalf("inline GovHealth config should route to committee, got %+v", dec.Approval)
	}
	vs, _ := svc.ListVersions(ctx, "sg")
	if len(vs) != 1 {
		t.Fatalf("preview must not persist a version, got %d", len(vs))
	}
}

func TestCompareConfigsDetectsDifference(t *testing.T) {
	svc := seeded(t)
	ctx := context.Background()
	_, _ = svc.SaveDraftConfig(ctx, "sg", seed.HealthFirst(), "hf", "t") // v2 draft
	two := 2
	changes, err := svc.CompareConfigs(ctx,
		usecase.ConfigRef{Slug: "sg"},                // active (v1, SafeGuard)
		usecase.ConfigRef{Slug: "sg", Version: &two}, // v2 (HealthFirst)
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected differences between SafeGuard and HealthFirst")
	}
}
