package usecase_test

import (
	"context"
	"errors"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/seed"
	"claimsplatform/internal/usecase"
)

func TestSaveDraftRejectsInvalidConfig(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	_, _ = svc.CreateTenant(ctx, "co", "Co", "")
	bad := seed.SafeGuard()
	bad.SLA.DefaultDays = 0 // invalid: must be >= 1
	_, err := svc.SaveDraftConfig(ctx, "co", bad, "oops", "tester")
	var ve domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestSaveDraftThenPublishSwitchesActive(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	_, _ = svc.CreateTenant(ctx, "co", "Co", "")
	v, err := svc.SaveDraftConfig(ctx, "co", seed.HealthFirst(), "to HF", "tester")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.PublishVersion(ctx, "co", v.VersionNumber); err != nil {
		t.Fatal(err)
	}
	cfg, _ := svc.GetActiveConfig(ctx, "co")
	if cfg.Approval.AutoApproveThreshold != seed.HealthFirst().Approval.AutoApproveThreshold {
		t.Fatal("active config should be the published draft")
	}
}

func TestRollbackRestoresOldConfig(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	_, _ = svc.CreateTenant(ctx, "co", "Co", "") // v1 default
	v2, _ := svc.SaveDraftConfig(ctx, "co", seed.GovHealth(), "gov", "tester")
	_ = svc.PublishVersion(ctx, "co", v2.VersionNumber)
	rolled, err := svc.RollbackVersion(ctx, "co", 1, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if rolled.CreatedBy != "tester" {
		t.Fatalf("rollback version must record the actor, got %q", rolled.CreatedBy)
	}
	cfg, _ := svc.GetActiveConfig(ctx, "co")
	if cfg.Approval.Model == seed.GovHealth().Approval.Model {
		t.Fatal("rollback to v1 should restore the default (tiered) model")
	}
}
