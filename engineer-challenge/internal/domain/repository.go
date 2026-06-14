package domain

import "context"

// ConfigurationRepository is the persistence port for the Tenant aggregate.
// Implementations live in infrastructure (configrepo). GORM never leaks past this line.
type ConfigurationRepository interface {
	CreateTenant(ctx context.Context, t Tenant, initial ConfigDocument) (Tenant, error)
	GetTenantBySlug(ctx context.Context, slug string) (Tenant, error)
	ListTenants(ctx context.Context) ([]Tenant, error)
	UpdateTenantMeta(ctx context.Context, slug, name, status string) (Tenant, error)

	CreateDraft(ctx context.Context, slug string, cfg ConfigDocument, note, by string) (ConfigVersion, error)
	Publish(ctx context.Context, slug string, version int) error
	Rollback(ctx context.Context, slug string, targetVersion int) (ConfigVersion, error)

	ListVersions(ctx context.Context, slug string) ([]ConfigVersion, error)
	GetVersion(ctx context.Context, slug string, version int) (ConfigVersion, error)
	GetActiveConfig(ctx context.Context, slug string) (ConfigDocument, error)
}
