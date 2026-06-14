package branding_test

import (
	"testing"

	"claimsplatform/internal/dimensions/branding"
	"claimsplatform/internal/domain"
)

func TestInvalidColorFails(t *testing.T) {
	cfg := domain.ConfigDocument{Branding: domain.BrandingConfig{DisplayName: "X", PrimaryColor: "red"}}
	if errs := branding.New().Validate(cfg); len(errs) == 0 {
		t.Fatal("expected validation error for non-hex color")
	}
}

func TestValidBrandingPasses(t *testing.T) {
	cfg := domain.ConfigDocument{Branding: domain.BrandingConfig{DisplayName: "X", PrimaryColor: "#0A4D2C", SecondaryColor: "#082B19"}}
	if errs := branding.New().Validate(cfg); len(errs) != 0 {
		t.Fatalf("expected valid branding, got %v", errs)
	}
}
