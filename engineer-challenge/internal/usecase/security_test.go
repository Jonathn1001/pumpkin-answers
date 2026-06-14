package usecase_test

import (
	"context"
	"regexp"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/usecase"
)

// L2: tenant slugs must be URL- and ref-parser-safe. Slugs are derived
// server-side (never client-supplied), so even a hostile name must still yield a
// slug matching ^[a-z0-9][a-z0-9-]{0,62}$.
func TestDerivedSlugIsAlwaysRefSafe(t *testing.T) {
	svc := usecase.New(memory.New())
	safe := regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)
	for _, name := range []string{
		"Bad Slug!", "has@at", "../etc/passwd", "with/slash", "a b@c", "DROP TABLE tenants;",
	} {
		tn, err := svc.CreateTenant(context.Background(), name, usecase.DefaultDocument())
		if err != nil {
			t.Fatalf("name %q: unexpected error %v", name, err)
		}
		if !safe.MatchString(tn.Slug) {
			t.Fatalf("name %q produced unsafe slug %q", name, tn.Slug)
		}
	}
}
