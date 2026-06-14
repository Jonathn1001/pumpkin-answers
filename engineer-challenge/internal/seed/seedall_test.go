package seed_test

import (
	"context"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/seed"
)

func TestSeedAllIsIdempotent(t *testing.T) {
	repo := memory.New()
	ctx := context.Background()
	if err := seed.SeedAll(ctx, repo); err != nil {
		t.Fatal(err)
	}
	if err := seed.SeedAll(ctx, repo); err != nil { // second run must not error or duplicate
		t.Fatal(err)
	}
	ts, _ := repo.ListTenants(ctx)
	if len(ts) != 3 {
		t.Fatalf("expected 3 seeded tenants, got %d", len(ts))
	}
}
