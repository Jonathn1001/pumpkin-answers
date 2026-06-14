package usecase

import (
	"context"

	"claimsplatform/internal/configvalidation"
	"claimsplatform/internal/domain"
)

// SaveDraftConfig validates first; an invalid config yields ValidationError and no draft.
func (s *Service) SaveDraftConfig(ctx context.Context, slug string, cfg domain.ConfigDocument, note, by string) (domain.ConfigVersion, error) {
	if errs := configvalidation.Validate(cfg); len(errs) > 0 {
		return domain.ConfigVersion{}, domain.ValidationError{Fields: errs}
	}
	return s.repo.CreateDraft(ctx, slug, cfg, note, by)
}

func (s *Service) PublishVersion(ctx context.Context, slug string, version int) error {
	return s.repo.Publish(ctx, slug, version)
}

func (s *Service) RollbackVersion(ctx context.Context, slug string, targetVersion int, by string) (domain.ConfigVersion, error) {
	return s.repo.Rollback(ctx, slug, targetVersion, by)
}

func (s *Service) ListVersions(ctx context.Context, slug string) ([]domain.ConfigVersion, error) {
	return s.repo.ListVersions(ctx, slug)
}

func (s *Service) GetVersion(ctx context.Context, slug string, version int) (domain.ConfigVersion, error) {
	return s.repo.GetVersion(ctx, slug, version)
}

func (s *Service) GetActiveConfig(ctx context.Context, slug string) (domain.ConfigDocument, error) {
	return s.repo.GetActiveConfig(ctx, slug)
}
