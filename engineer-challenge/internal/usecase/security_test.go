package usecase_test

import (
	"context"
	"errors"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/usecase"
)

// L2: tenant slugs must be URL- and ref-parser-safe; invalid slugs are rejected.
func TestCreateTenantRejectsInvalidSlug(t *testing.T) {
	svc := usecase.New(memory.New())
	for _, bad := range []string{"Bad Slug!", "has@at", "UPPER", "-leading", "with/slash", ""} {
		_, err := svc.CreateTenant(context.Background(), bad, "X", "")
		var ve domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("slug %q should be rejected with ValidationError, got %v", bad, err)
		}
	}
}

func TestCreateTenantAcceptsValidSlug(t *testing.T) {
	svc := usecase.New(memory.New())
	if _, err := svc.CreateTenant(context.Background(), "valid-slug-1", "X", ""); err != nil {
		t.Fatalf("valid slug was rejected: %v", err)
	}
}
