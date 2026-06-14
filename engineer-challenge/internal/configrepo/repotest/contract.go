// Package repotest holds the shared ConfigurationRepository contract: one behavior
// suite both the in-memory fake and the GORM implementation must satisfy.
package repotest

import (
	"context"
	"errors"
	"testing"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/seed"
)

// Contract runs the full repository behavior suite against a fresh repo from newRepo.
func Contract(t *testing.T, newRepo func(t *testing.T) domain.ConfigurationRepository) {
	mk := func(t *testing.T, slug string, cfg domain.ConfigDocument) domain.ConfigurationRepository {
		repo := newRepo(t)
		_, err := repo.CreateTenant(context.Background(),
			domain.Tenant{Slug: slug, Name: slug, Status: domain.TenantActive}, cfg)
		if err != nil {
			t.Fatalf("seed CreateTenant: %v", err)
		}
		return repo
	}

	t.Run("CreateTenant publishes version 1 and sets active pointer", func(t *testing.T) {
		repo := newRepo(t)
		ctx := context.Background()
		tn, err := repo.CreateTenant(ctx, domain.Tenant{Slug: "acme", Name: "Acme", Status: domain.TenantActive}, seed.SafeGuard())
		if err != nil {
			t.Fatal(err)
		}
		if tn.ActiveVersionNumber == nil || *tn.ActiveVersionNumber != 1 {
			t.Fatalf("expected active version 1, got %v", tn.ActiveVersionNumber)
		}
		cfg, err := repo.GetActiveConfig(ctx, "acme")
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Approval.AutoApproveThreshold != seed.SafeGuard().Approval.AutoApproveThreshold {
			t.Fatal("active config should equal the initial config")
		}
		vs, err := repo.ListVersions(ctx, "acme")
		if err != nil {
			t.Fatal(err)
		}
		if len(vs) != 1 || vs[0].Status != domain.VersionPublished {
			t.Fatalf("expected one published version, got %+v", vs)
		}
	})

	t.Run("duplicate slug returns ErrSlugTaken", func(t *testing.T) {
		repo := mk(t, "dup", seed.SafeGuard())
		_, err := repo.CreateTenant(context.Background(),
			domain.Tenant{Slug: "dup", Name: "B", Status: domain.TenantActive}, seed.SafeGuard())
		if !errors.Is(err, domain.ErrSlugTaken) {
			t.Fatalf("expected ErrSlugTaken, got %v", err)
		}
	})

	t.Run("unknown tenant returns ErrTenantNotFound", func(t *testing.T) {
		repo := newRepo(t)
		_, err := repo.GetTenantBySlug(context.Background(), "ghost")
		if !errors.Is(err, domain.ErrTenantNotFound) {
			t.Fatalf("expected ErrTenantNotFound, got %v", err)
		}
	})

	t.Run("CreateDraft does not change active; Publish switches the pointer", func(t *testing.T) {
		repo := mk(t, "hf", seed.SafeGuard())
		ctx := context.Background()
		draft, err := repo.CreateDraft(ctx, "hf", seed.HealthFirst(), "switch to HF rules", "tester")
		if err != nil {
			t.Fatal(err)
		}
		if draft.VersionNumber != 2 || draft.Status != domain.VersionDraft {
			t.Fatalf("expected draft v2, got %+v", draft)
		}
		cfg, _ := repo.GetActiveConfig(ctx, "hf")
		if cfg.Approval.AutoApproveThreshold != seed.SafeGuard().Approval.AutoApproveThreshold {
			t.Fatal("active should still be v1 before publish")
		}
		if err := repo.Publish(ctx, "hf", 2); err != nil {
			t.Fatal(err)
		}
		cfg, _ = repo.GetActiveConfig(ctx, "hf")
		if cfg.Approval.AutoApproveThreshold != seed.HealthFirst().Approval.AutoApproveThreshold {
			t.Fatal("active should be v2 after publish")
		}
		tn, _ := repo.GetTenantBySlug(ctx, "hf")
		if tn.ActiveVersionNumber == nil || *tn.ActiveVersionNumber != 2 {
			t.Fatalf("expected active pointer 2, got %v", tn.ActiveVersionNumber)
		}
	})

	t.Run("Rollback clones target into a new published version, preserving history", func(t *testing.T) {
		repo := mk(t, "gv", seed.SafeGuard()) // v1
		ctx := context.Background()
		_, _ = repo.CreateDraft(ctx, "gv", seed.GovHealth(), "gov", "tester") // v2
		_ = repo.Publish(ctx, "gv", 2)
		newV, err := repo.Rollback(ctx, "gv", 1) // back to SafeGuard rules
		if err != nil {
			t.Fatal(err)
		}
		if newV.VersionNumber != 3 {
			t.Fatalf("rollback should create v3, got %d", newV.VersionNumber)
		}
		cfg, _ := repo.GetActiveConfig(ctx, "gv")
		if cfg.Approval.Model != seed.SafeGuard().Approval.Model {
			t.Fatal("active after rollback should match v1 config")
		}
		vs, _ := repo.ListVersions(ctx, "gv")
		if len(vs) != 3 {
			t.Fatalf("history must be preserved (3 versions), got %d", len(vs))
		}
	})

	t.Run("GetVersion unknown number returns ErrVersionNotFound", func(t *testing.T) {
		repo := mk(t, "v", seed.SafeGuard())
		_, err := repo.GetVersion(context.Background(), "v", 99)
		if !errors.Is(err, domain.ErrVersionNotFound) {
			t.Fatalf("expected ErrVersionNotFound, got %v", err)
		}
	})

	t.Run("UpdateTenantMeta changes name and status", func(t *testing.T) {
		repo := mk(t, "m", seed.SafeGuard())
		tn, err := repo.UpdateTenantMeta(context.Background(), "m", "New Name", domain.TenantArchived)
		if err != nil {
			t.Fatal(err)
		}
		if tn.Name != "New Name" || tn.Status != domain.TenantArchived {
			t.Fatalf("meta not updated: %+v", tn)
		}
	})

	t.Run("ListTenants returns all created tenants", func(t *testing.T) {
		repo := mk(t, "a", seed.SafeGuard())
		_, _ = repo.CreateTenant(context.Background(),
			domain.Tenant{Slug: "b", Name: "B", Status: domain.TenantActive}, seed.HealthFirst())
		ts, err := repo.ListTenants(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if len(ts) < 2 {
			t.Fatalf("expected >=2 tenants, got %d", len(ts))
		}
	})
}
