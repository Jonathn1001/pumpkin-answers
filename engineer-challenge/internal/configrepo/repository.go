package configrepo

import (
	"context"
	"errors"
	"strconv"
	"time"

	"claimsplatform/internal/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type repo struct{ db *gorm.DB }

// New returns a GORM-backed ConfigurationRepository. Tables must already exist (run Migrate first).
func New(db *gorm.DB) domain.ConfigurationRepository { return &repo{db: db} }

var _ domain.ConfigurationRepository = (*repo)(nil)

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func tenantIDBySlug(tx *gorm.DB, slug string) (int64, error) {
	var tm tenantModel
	res := tx.Select("id").Where("slug = ?", slug).Take(&tm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, domain.ErrTenantNotFound
		}
		return 0, res.Error
	}
	return tm.ID, nil
}

func nextVersionNumber(tx *gorm.DB, tenantID int64) (int, error) {
	var row struct{ MaxVersion *int }
	if err := tx.Model(&configVersionModel{}).Where("tenant_id = ?", tenantID).
		Select("MAX(version_number) AS max_version").Scan(&row).Error; err != nil {
		return 0, err
	}
	if row.MaxVersion == nil {
		return 1, nil
	}
	return *row.MaxVersion + 1, nil
}

func (r *repo) CreateTenant(ctx context.Context, t domain.Tenant, initial domain.ConfigDocument) (domain.Tenant, error) {
	cfgJSON, err := toJSON(initial)
	if err != nil {
		return domain.Tenant{}, err
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&tenantModel{}).Where("slug = ?", t.Slug).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return domain.ErrSlugTaken
		}
		tm := tenantModel{Slug: t.Slug, Name: t.Name, Status: orDefault(t.Status, domain.TenantActive)}
		if err := tx.Create(&tm).Error; err != nil {
			return err
		}
		vm := configVersionModel{TenantID: tm.ID, VersionNumber: 1, Status: domain.VersionPublished, Config: cfgJSON}
		if err := tx.Create(&vm).Error; err != nil {
			return err
		}
		return tx.Model(&tenantModel{}).Where("id = ?", tm.ID).
			Update("active_config_version_id", vm.ID).Error
	})
	if err != nil {
		return domain.Tenant{}, err
	}
	return r.GetTenantBySlug(ctx, t.Slug)
}

func (r *repo) GetTenantBySlug(ctx context.Context, slug string) (domain.Tenant, error) {
	var row tenantRow
	res := r.db.WithContext(ctx).
		Table("tenants").
		Select("tenants.*, av.version_number AS active_version_number").
		Joins("LEFT JOIN config_versions av ON av.id = tenants.active_config_version_id").
		Where("tenants.slug = ?", slug).
		Scan(&row)
	if res.Error != nil {
		return domain.Tenant{}, res.Error
	}
	if res.RowsAffected == 0 {
		return domain.Tenant{}, domain.ErrTenantNotFound
	}
	return toTenant(row), nil
}

func (r *repo) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	var rows []tenantRow
	if err := r.db.WithContext(ctx).
		Table("tenants").
		Select("tenants.*, av.version_number AS active_version_number").
		Joins("LEFT JOIN config_versions av ON av.id = tenants.active_config_version_id").
		Order("tenants.id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Tenant, 0, len(rows))
	for _, row := range rows {
		out = append(out, toTenant(row))
	}
	return out, nil
}

func (r *repo) UpdateTenantMeta(ctx context.Context, slug, name, status string) (domain.Tenant, error) {
	res := r.db.WithContext(ctx).
		Table("tenants").
		Where("slug = ?", slug).
		Updates(map[string]any{"name": name, "status": status, "updated_at": time.Now().UTC()})
	if res.Error != nil {
		return domain.Tenant{}, res.Error
	}
	if res.RowsAffected == 0 {
		return domain.Tenant{}, domain.ErrTenantNotFound
	}
	return r.GetTenantBySlug(ctx, slug)
}

func (r *repo) CreateDraft(ctx context.Context, slug string, cfg domain.ConfigDocument, note, by string) (domain.ConfigVersion, error) {
	cfgJSON, err := toJSON(cfg)
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	var vm configVersionModel
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tid, err := tenantIDBySlug(tx, slug)
		if err != nil {
			return err
		}
		next, err := nextVersionNumber(tx, tid)
		if err != nil {
			return err
		}
		vm = configVersionModel{TenantID: tid, VersionNumber: next, Status: domain.VersionDraft, Note: note, Config: cfgJSON, CreatedBy: by}
		return tx.Create(&vm).Error
	})
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	return toConfigVersion(slug, vm, cfg), nil
}

func (r *repo) Publish(ctx context.Context, slug string, version int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tid, err := tenantIDBySlug(tx, slug)
		if err != nil {
			return err
		}
		var vm configVersionModel
		res := tx.Where("tenant_id = ? AND version_number = ?", tid, version).Take(&vm)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return domain.ErrVersionNotFound
			}
			return res.Error
		}
		if err := tx.Model(&configVersionModel{}).Where("id = ?", vm.ID).
			Update("status", domain.VersionPublished).Error; err != nil {
			return err
		}
		return tx.Model(&tenantModel{}).Where("id = ?", tid).
			Updates(map[string]any{"active_config_version_id": vm.ID, "updated_at": time.Now().UTC()}).Error
	})
}

func (r *repo) Rollback(ctx context.Context, slug string, targetVersion int, by string) (domain.ConfigVersion, error) {
	var clone configVersionModel
	var doc domain.ConfigDocument
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tid, err := tenantIDBySlug(tx, slug)
		if err != nil {
			return err
		}
		var target configVersionModel
		res := tx.Where("tenant_id = ? AND version_number = ?", tid, targetVersion).Take(&target)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return domain.ErrVersionNotFound
			}
			return res.Error
		}
		next, err := nextVersionNumber(tx, tid)
		if err != nil {
			return err
		}
		clone = configVersionModel{TenantID: tid, VersionNumber: next, Status: domain.VersionPublished,
			Note: "rollback to v" + strconv.Itoa(targetVersion), CreatedBy: by, Config: target.Config}
		if err := tx.Create(&clone).Error; err != nil {
			return err
		}
		doc, err = toConfigDocument(clone.Config)
		if err != nil {
			return err
		}
		return tx.Model(&tenantModel{}).Where("id = ?", tid).
			Updates(map[string]any{"active_config_version_id": clone.ID, "updated_at": time.Now().UTC()}).Error
	})
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	return toConfigVersion(slug, clone, doc), nil
}

func (r *repo) ListVersions(ctx context.Context, slug string) ([]domain.ConfigVersion, error) {
	tid, err := tenantIDBySlug(r.db.WithContext(ctx), slug)
	if err != nil {
		return nil, err
	}
	var vms []configVersionModel
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tid).
		Order("version_number DESC").Find(&vms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ConfigVersion, 0, len(vms))
	for _, vm := range vms {
		d, err := toConfigDocument(vm.Config)
		if err != nil {
			return nil, err
		}
		out = append(out, toConfigVersion(slug, vm, d))
	}
	return out, nil
}

func (r *repo) GetVersion(ctx context.Context, slug string, version int) (domain.ConfigVersion, error) {
	tid, err := tenantIDBySlug(r.db.WithContext(ctx), slug)
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	var vm configVersionModel
	res := r.db.WithContext(ctx).Where("tenant_id = ? AND version_number = ?", tid, version).Take(&vm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return domain.ConfigVersion{}, domain.ErrVersionNotFound
		}
		return domain.ConfigVersion{}, res.Error
	}
	d, err := toConfigDocument(vm.Config)
	if err != nil {
		return domain.ConfigVersion{}, err
	}
	return toConfigVersion(slug, vm, d), nil
}

func (r *repo) GetActiveConfig(ctx context.Context, slug string) (domain.ConfigDocument, error) {
	var row struct{ Config datatypes.JSON }
	res := r.db.WithContext(ctx).
		Table("config_versions AS cv").
		Select("cv.config AS config").
		Joins("JOIN tenants t ON t.active_config_version_id = cv.id").
		Where("t.slug = ?", slug).
		Scan(&row)
	if res.Error != nil {
		return domain.ConfigDocument{}, res.Error
	}
	if res.RowsAffected == 0 {
		if _, err := r.GetTenantBySlug(ctx, slug); err != nil {
			return domain.ConfigDocument{}, err
		}
		return domain.ConfigDocument{}, domain.ErrVersionNotFound
	}
	return toConfigDocument(row.Config)
}
