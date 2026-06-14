// Package usecase holds the application services that orchestrate domain + persistence.
package usecase

import (
	"context"

	"claimsplatform/internal/configvalidation"
	"claimsplatform/internal/domain"
)

type Service struct {
	repo domain.ConfigurationRepository
}

func New(repo domain.ConfigurationRepository) *Service { return &Service{repo: repo} }

// CreateTenant starts from the default config, or clones cloneFrom's active config when set.
func (s *Service) CreateTenant(ctx context.Context, slug, name, cloneFrom string) (domain.Tenant, error) {
	cfg := DefaultDocument()
	if cloneFrom != "" {
		src, err := s.repo.GetActiveConfig(ctx, cloneFrom)
		if err != nil {
			return domain.Tenant{}, err
		}
		cfg = src
	}
	if errs := configvalidation.Validate(cfg); len(errs) > 0 {
		return domain.Tenant{}, domain.ValidationError{Fields: errs}
	}
	return s.repo.CreateTenant(ctx, domain.Tenant{Slug: slug, Name: name, Status: domain.TenantActive}, cfg)
}

func (s *Service) GetTenant(ctx context.Context, slug string) (domain.Tenant, error) {
	return s.repo.GetTenantBySlug(ctx, slug)
}

func (s *Service) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	return s.repo.ListTenants(ctx)
}

func (s *Service) UpdateTenant(ctx context.Context, slug, name, status string) (domain.Tenant, error) {
	if status == "" {
		status = domain.TenantActive
	}
	if status != domain.TenantActive && status != domain.TenantArchived {
		return domain.Tenant{}, domain.ValidationError{Fields: []domain.FieldError{
			{Field: "status", Message: "must be 'active' or 'archived'"},
		}}
	}
	return s.repo.UpdateTenantMeta(ctx, slug, name, status)
}
