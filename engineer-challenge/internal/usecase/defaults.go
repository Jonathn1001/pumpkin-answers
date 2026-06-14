package usecase

import (
	_ "claimsplatform/internal/dimensions" // register dimensions (engine + defaults need them)
	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

// DefaultDocument assembles a valid starter config from each dimension's own default,
// so a brand-new tenant is valid and processable immediately. Unknown future
// dimensions are ignored here (forward-compatible).
func DefaultDocument() domain.ConfigDocument {
	var doc domain.ConfigDocument
	for _, d := range registry.All() {
		switch v := d.DefaultConfig().(type) {
		case domain.BrandingConfig:
			doc.Branding = v
		case map[domain.ClaimType]domain.ClaimTypeConfig:
			doc.ClaimTypes = v
		case domain.ApprovalConfig:
			doc.Approval = v
		case domain.NotificationsConfig:
			doc.Notifications = v
		case domain.SLAConfig:
			doc.SLA = v
		case []domain.CustomFieldConfig:
			doc.CustomFields = v
		}
	}
	return doc
}
