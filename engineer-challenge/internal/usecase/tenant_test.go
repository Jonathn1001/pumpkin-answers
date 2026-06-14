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
	tn, err := svc.CreateTenant(context.Background(), "newco", "New Co", "")
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
	if _, err := svc.CreateTenant(ctx, "co", "Co", ""); err != nil {
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
	if _, err := svc.CreateTenant(ctx, "co", "Co", ""); err != nil {
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

func TestCreateTenantCloneCopiesSourceConfig(t *testing.T) {
	svc := usecase.New(memory.New())
	ctx := context.Background()
	if _, err := svc.CreateTenant(ctx, "src", "Source", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateTenant(ctx, "clone", "Clone", "src"); err != nil {
		t.Fatal(err)
	}
	a, _ := svc.GetActiveConfig(ctx, "src")
	b, _ := svc.GetActiveConfig(ctx, "clone")
	if a.Approval.Model != b.Approval.Model {
		t.Fatal("clone should copy source config")
	}
}
