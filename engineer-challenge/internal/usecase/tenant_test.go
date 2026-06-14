package usecase_test

import (
	"context"
	"errors"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/configvalidation"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/usecase"
)

func TestDefaultDocumentIsValid(t *testing.T) {
	if errs := configvalidation.Validate(usecase.DefaultDocument()); len(errs) != 0 {
		t.Fatalf("default document must be valid, got %v", errs)
	}
}

func TestCreateTenantStartsFromValidDefault(t *testing.T) {
	svc := usecase.New(memory.New())
	tn, err := svc.CreateTenant(context.Background(), "New Co", usecase.DefaultDocument())
	if err != nil {
		t.Fatal(err)
	}
	if tn.ActiveVersionNumber == nil || *tn.ActiveVersionNumber != 1 {
		t.Fatalf("expected active version 1, got %v", tn.ActiveVersionNumber)
	}
}

func TestUpdateTenantRejectsInvalidStatus(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	if _, err := svc.CreateTenant(ctx, "Co", usecase.DefaultDocument()); err != nil {
		t.Fatal(err)
	}
	_, err := svc.UpdateTenant(ctx, "co", "New", "banana")
	var ve domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for invalid status, got %v", err)
	}
}

func TestUpdateTenantAcceptsArchivedStatus(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	if _, err := svc.CreateTenant(ctx, "Co", usecase.DefaultDocument()); err != nil {
		t.Fatal(err)
	}
	tn, err := svc.UpdateTenant(ctx, "co", "Co", domain.TenantArchived)
	if err != nil {
		t.Fatalf("expected no error for valid 'archived' status, got %v", err)
	}
	if tn.Status != domain.TenantArchived {
		t.Fatalf("expected status 'archived', got %q", tn.Status)
	}
}

func TestCreateTenantDerivesSlugFromName(t *testing.T) {
	svc := usecase.New(memory.New())
	tn, err := svc.CreateTenant(context.Background(), "SafeGuard Insurance", usecase.DefaultDocument())
	if err != nil {
		t.Fatal(err)
	}
	if tn.Slug != "safeguard-insurance" {
		t.Fatalf("derived slug = %q, want %q", tn.Slug, "safeguard-insurance")
	}
}

func TestCreateTenantAutoSuffixesDuplicateDerivedSlug(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	if _, err := svc.CreateTenant(ctx, "Acme", usecase.DefaultDocument()); err != nil {
		t.Fatal(err)
	}
	tn, err := svc.CreateTenant(ctx, "Acme", usecase.DefaultDocument())
	if err != nil {
		t.Fatalf("second create should auto-suffix, got %v", err)
	}
	if tn.Slug != "acme-2" {
		t.Fatalf("second slug = %q, want %q", tn.Slug, "acme-2")
	}
}

func TestCreateTenantRejectsNameWithNoSlugChars(t *testing.T) {
	svc := usecase.New(memory.New())
	_, err := svc.CreateTenant(context.Background(), "@@@", usecase.DefaultDocument())
	var ve domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if len(ve.Fields) == 0 || ve.Fields[0].Field != "name" {
		t.Fatalf("expected validation on 'name', got %+v", ve.Fields)
	}
}

// The client sends the full config; create persists it verbatim as v1.
func TestCreateTenantUsesProvidedConfig(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	cfg := usecase.DefaultDocument()
	cfg.Branding.DisplayName = "Provided Brand"
	if _, err := svc.CreateTenant(ctx, "Provided Co", cfg); err != nil {
		t.Fatal(err)
	}
	got, _ := svc.GetActiveConfig(ctx, "provided-co")
	if got.Branding.DisplayName != "Provided Brand" {
		t.Fatalf("create should persist the provided config, got displayName %q", got.Branding.DisplayName)
	}
}

// An invalid client config is rejected before any tenant is written.
func TestCreateTenantRejectsInvalidConfig(t *testing.T) {
	svc := usecase.New(memory.New())
	bad := usecase.DefaultDocument()
	bad.Branding.DisplayName = "" // branding.displayName is required
	_, err := svc.CreateTenant(context.Background(), "Bad Co", bad)
	var ve domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for invalid config, got %v", err)
	}
}
