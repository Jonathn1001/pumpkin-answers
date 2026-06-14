package domain

import "time"

// Tenant status + ConfigVersion status vocabularies.
const (
	TenantActive   = "active"
	TenantArchived = "archived"

	VersionDraft     = "draft"
	VersionPublished = "published"
)

// Tenant is the aggregate root. ActiveVersionNumber is the published pointer (nil before version #1).
type Tenant struct {
	Slug                string    `json:"slug"`
	Name                string    `json:"name"`
	Status              string    `json:"status"`
	ActiveVersionNumber *int      `json:"activeVersionNumber,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// ConfigVersion is an immutable entity inside the Tenant aggregate.
type ConfigVersion struct {
	TenantSlug    string         `json:"tenantSlug"`
	VersionNumber int            `json:"versionNumber"`
	Status        string         `json:"status"`
	Note          string         `json:"note,omitempty"`
	CreatedBy     string         `json:"createdBy,omitempty"`
	Config        ConfigDocument `json:"config"`
	CreatedAt     time.Time      `json:"createdAt"`
}
