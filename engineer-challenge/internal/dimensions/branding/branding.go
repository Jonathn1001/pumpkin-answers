package branding

import (
	"regexp"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

var hexColor = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type dimension struct{}

func New() registry.Dimension { return dimension{} }

func init() { registry.Register(New()) }

func (dimension) Key() string { return registry.KeyBranding }

// Branding does not affect the claim decision.
func (dimension) Evaluate(_ domain.ConfigDocument, _ domain.Claim, _ *domain.ClaimDecision) {}

func (dimension) Validate(cfg domain.ConfigDocument) []domain.FieldError {
	var errs []domain.FieldError
	b := cfg.Branding
	if b.DisplayName == "" {
		errs = append(errs, domain.FieldError{Field: "branding.displayName", Message: "display name is required"})
	}
	if b.PrimaryColor != "" && !hexColor.MatchString(b.PrimaryColor) {
		errs = append(errs, domain.FieldError{Field: "branding.primaryColor", Message: "must be a #RRGGBB hex color"})
	}
	if b.SecondaryColor != "" && !hexColor.MatchString(b.SecondaryColor) {
		errs = append(errs, domain.FieldError{Field: "branding.secondaryColor", Message: "must be a #RRGGBB hex color"})
	}
	return errs
}

func (dimension) DefaultConfig() any {
	return domain.BrandingConfig{DisplayName: "New Tenant", PrimaryColor: "#1F2937", SecondaryColor: "#374151"}
}

func (dimension) UISchema() []registry.FieldDescriptor {
	return []registry.FieldDescriptor{
		{Key: "branding.displayName", Label: "Display name", Type: "string", Widget: "text", Required: true},
		{Key: "branding.logoUrl", Label: "Logo URL", Type: "string", Widget: "text"},
		{Key: "branding.primaryColor", Label: "Primary color", Type: "string", Widget: "color"},
		{Key: "branding.secondaryColor", Label: "Secondary color", Type: "string", Widget: "color"},
		{Key: "branding.supportEmail", Label: "Support email", Type: "string", Widget: "text"},
	}
}
