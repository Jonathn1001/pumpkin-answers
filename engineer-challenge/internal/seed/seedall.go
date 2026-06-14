package seed

import (
	"context"
	"errors"

	"claimsplatform/internal/domain"
)

// SeedAll inserts the three sample tenants if they are not already present (idempotent).
func SeedAll(ctx context.Context, repo domain.ConfigurationRepository) error {
	tenants := []struct {
		slug, name string
		cfg        domain.ConfigDocument
	}{
		{"safeguard", "SafeGuard Insurance", SafeGuard()},
		{"healthfirst", "HealthFirst", HealthFirst()},
		{"govhealth", "GovHealth", GovHealth()},
	}
	for _, t := range tenants {
		_, err := repo.GetTenantBySlug(ctx, t.slug)
		if err == nil {
			continue // already seeded
		}
		if !errors.Is(err, domain.ErrTenantNotFound) {
			return err
		}
		if _, err := repo.CreateTenant(ctx,
			domain.Tenant{Slug: t.slug, Name: t.name, Status: domain.TenantActive}, t.cfg); err != nil {
			return err
		}
	}
	return nil
}
