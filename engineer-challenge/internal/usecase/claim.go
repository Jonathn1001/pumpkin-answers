package usecase

import (
	"context"

	"claimsplatform/internal/comparison"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/engine"
)

// ConfigRef points at a tenant's active config (Version nil) or a specific version.
type ConfigRef struct {
	Slug    string
	Version *int
}

func (s *Service) ProcessClaim(ctx context.Context, slug string, claim domain.Claim) (domain.ClaimDecision, error) {
	if errs := claim.Validate(); len(errs) > 0 {
		return domain.ClaimDecision{}, domain.ValidationError{Fields: errs}
	}
	cfg, err := s.repo.GetActiveConfig(ctx, slug)
	if err != nil {
		return domain.ClaimDecision{}, err
	}
	return engine.ProcessClaim(cfg, claim), nil
}

// PreviewClaim resolves config with precedence inline > version > active, runs the
// engine, and persists nothing.
func (s *Service) PreviewClaim(ctx context.Context, slug string, claim domain.Claim, version *int, inline *domain.ConfigDocument) (domain.ClaimDecision, error) {
	if errs := claim.Validate(); len(errs) > 0 {
		return domain.ClaimDecision{}, domain.ValidationError{Fields: errs}
	}
	cfg, err := s.resolve(ctx, ConfigRef{Slug: slug, Version: version}, inline)
	if err != nil {
		return domain.ClaimDecision{}, err
	}
	return engine.ProcessClaim(cfg, claim), nil
}

func (s *Service) CompareConfigs(ctx context.Context, left, right ConfigRef) ([]comparison.Change, error) {
	l, err := s.resolve(ctx, left, nil)
	if err != nil {
		return nil, err
	}
	r, err := s.resolve(ctx, right, nil)
	if err != nil {
		return nil, err
	}
	return comparison.Diff(l, r)
}

func (s *Service) resolve(ctx context.Context, ref ConfigRef, inline *domain.ConfigDocument) (domain.ConfigDocument, error) {
	if inline != nil {
		return *inline, nil
	}
	if ref.Version != nil {
		v, err := s.repo.GetVersion(ctx, ref.Slug, *ref.Version)
		if err != nil {
			return domain.ConfigDocument{}, err
		}
		return v.Config, nil
	}
	return s.repo.GetActiveConfig(ctx, ref.Slug)
}
