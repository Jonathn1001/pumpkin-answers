// Package usecase holds the application services that orchestrate domain + persistence.
package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"claimsplatform/internal/configvalidation"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/slug"
)

type Service struct {
	repo domain.ConfigurationRepository
}

func New(repo domain.ConfigurationRepository) *Service { return &Service{repo: repo} }

// CreateTenant derives a server-authoritative slug from name (slugs are never
// client-supplied). It starts from the default config, or clones cloneFrom's
// active config when set.
func (s *Service) CreateTenant(ctx context.Context, name, cloneFrom string) (domain.Tenant, error) {
	base := slug.Make(name)
	if base == "" {
		return domain.Tenant{}, domain.ValidationError{Fields: []domain.FieldError{
			{Field: "name", Message: "must contain at least one letter or digit to derive a slug"},
		}}
	}
	tenantSlug, err := s.uniqueSlug(ctx, base)
	if err != nil {
		return domain.Tenant{}, err
	}
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
	return s.repo.CreateTenant(ctx, domain.Tenant{Slug: tenantSlug, Name: name, Status: domain.TenantActive}, cfg)
}

// uniqueSlug appends -2, -3, … until the slug is free, keeping within 63 chars.
func (s *Service) uniqueSlug(ctx context.Context, base string) (string, error) {
	candidate := base
	for i := 2; ; i++ {
		_, err := s.repo.GetTenantBySlug(ctx, candidate)
		if errors.Is(err, domain.ErrTenantNotFound) {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}
		suffix := "-" + strconv.Itoa(i)
		trimmed := base
		if len(trimmed)+len(suffix) > 63 {
			trimmed = strings.TrimRight(base[:63-len(suffix)], "-")
		}
		candidate = trimmed + suffix
	}
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
