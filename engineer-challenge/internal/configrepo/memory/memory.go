// Package memory is an in-memory ConfigurationRepository for fast use-case and HTTP
// tests. Its observable behavior mirrors the GORM implementation.
package memory

import (
	"context"
	"strconv"
	"sync"
	"time"

	"claimsplatform/internal/domain"
)

type record struct {
	tenant   domain.Tenant
	versions []domain.ConfigVersion // versions[i] == version i+1
}

type Repo struct {
	mu  sync.Mutex
	rec map[string]*record
}

func New() *Repo { return &Repo{rec: map[string]*record{}} }

var _ domain.ConfigurationRepository = (*Repo)(nil)

func (r *Repo) get(slug string) (*record, error) {
	rec, ok := r.rec[slug]
	if !ok {
		return nil, domain.ErrTenantNotFound
	}
	return rec, nil
}

func (r *Repo) CreateTenant(_ context.Context, t domain.Tenant, initial domain.ConfigDocument) (domain.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rec[t.Slug]; ok {
		return domain.Tenant{}, domain.ErrSlugTaken
	}
	now := time.Now().UTC()
	t.CreatedAt, t.UpdatedAt = now, now
	one := 1
	t.ActiveVersionNumber = &one
	v1 := domain.ConfigVersion{TenantSlug: t.Slug, VersionNumber: 1, Status: domain.VersionPublished, Config: initial, CreatedAt: now}
	r.rec[t.Slug] = &record{tenant: t, versions: []domain.ConfigVersion{v1}}
	return t, nil
}

func (r *Repo) GetTenantBySlug(_ context.Context, slug string) (domain.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return domain.Tenant{}, err
	}
	return rec.tenant, nil
}

func (r *Repo) ListTenants(_ context.Context) ([]domain.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]domain.Tenant, 0, len(r.rec))
	for _, rec := range r.rec {
		out = append(out, rec.tenant)
	}
	return out, nil
}

func (r *Repo) UpdateTenantMeta(_ context.Context, slug, name, status string) (domain.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return domain.Tenant{}, err
	}
	rec.tenant.Name = name
	rec.tenant.Status = status
	rec.tenant.UpdatedAt = time.Now().UTC()
	return rec.tenant, nil
}

func (r *Repo) CreateDraft(_ context.Context, slug string, cfg domain.ConfigDocument, note, by string) (domain.ConfigVersion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	v := domain.ConfigVersion{
		TenantSlug:    slug,
		VersionNumber: len(rec.versions) + 1,
		Status:        domain.VersionDraft,
		Note:          note,
		CreatedBy:     by,
		Config:        cfg,
		CreatedAt:     time.Now().UTC(),
	}
	rec.versions = append(rec.versions, v)
	return v, nil
}

func (r *Repo) Publish(_ context.Context, slug string, version int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return err
	}
	idx := version - 1
	if idx < 0 || idx >= len(rec.versions) {
		return domain.ErrVersionNotFound
	}
	rec.versions[idx].Status = domain.VersionPublished
	rec.tenant.ActiveVersionNumber = &version
	rec.tenant.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *Repo) Rollback(_ context.Context, slug string, targetVersion int, by string) (domain.ConfigVersion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	idx := targetVersion - 1
	if idx < 0 || idx >= len(rec.versions) {
		return domain.ConfigVersion{}, domain.ErrVersionNotFound
	}
	newNum := len(rec.versions) + 1
	v := domain.ConfigVersion{
		TenantSlug:    slug,
		VersionNumber: newNum,
		Status:        domain.VersionPublished,
		Note:          "rollback to v" + strconv.Itoa(targetVersion),
		CreatedBy:     by,
		Config:        rec.versions[idx].Config,
		CreatedAt:     time.Now().UTC(),
	}
	rec.versions = append(rec.versions, v)
	rec.tenant.ActiveVersionNumber = &newNum
	rec.tenant.UpdatedAt = time.Now().UTC()
	return v, nil
}

func (r *Repo) ListVersions(_ context.Context, slug string) ([]domain.ConfigVersion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return nil, err
	}
	out := make([]domain.ConfigVersion, len(rec.versions))
	for i, v := range rec.versions { // newest first
		out[len(rec.versions)-1-i] = v
	}
	return out, nil
}

func (r *Repo) GetVersion(_ context.Context, slug string, version int) (domain.ConfigVersion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	idx := version - 1
	if idx < 0 || idx >= len(rec.versions) {
		return domain.ConfigVersion{}, domain.ErrVersionNotFound
	}
	return rec.versions[idx], nil
}

func (r *Repo) GetActiveConfig(_ context.Context, slug string) (domain.ConfigDocument, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, err := r.get(slug)
	if err != nil {
		return domain.ConfigDocument{}, err
	}
	if rec.tenant.ActiveVersionNumber == nil {
		return domain.ConfigDocument{}, domain.ErrVersionNotFound
	}
	return rec.versions[*rec.tenant.ActiveVersionNumber-1].Config, nil
}
