package configrepo

import (
	"encoding/json"
	"time"

	"claimsplatform/internal/domain"

	"gorm.io/datatypes"
)

type tenantModel struct {
	ID                    int64 `gorm:"primaryKey"`
	Slug                  string
	Name                  string
	Status                string
	ActiveConfigVersionID *int64
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (tenantModel) TableName() string { return "tenants" }

type configVersionModel struct {
	ID            int64 `gorm:"primaryKey"`
	TenantID      int64
	VersionNumber int
	Status        string
	Note          string
	Config        datatypes.JSON
	CreatedBy     string
	CreatedAt     time.Time
}

func (configVersionModel) TableName() string { return "config_versions" }

// tenantRow is a flat struct for scanning tenant + active version_number from a JOIN query.
type tenantRow struct {
	ID                    int64     `gorm:"column:id"`
	Slug                  string    `gorm:"column:slug"`
	Name                  string    `gorm:"column:name"`
	Status                string    `gorm:"column:status"`
	ActiveConfigVersionID *int64    `gorm:"column:active_config_version_id"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
	ActiveVersionNumber   *int      `gorm:"column:active_version_number"`
}

func toTenant(row tenantRow) domain.Tenant {
	return domain.Tenant{
		Slug:                row.Slug,
		Name:                row.Name,
		Status:              row.Status,
		ActiveVersionNumber: row.ActiveVersionNumber,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
}

func toConfigVersion(slug string, vm configVersionModel, doc domain.ConfigDocument) domain.ConfigVersion {
	return domain.ConfigVersion{
		TenantSlug:    slug,
		VersionNumber: vm.VersionNumber,
		Status:        vm.Status,
		Note:          vm.Note,
		CreatedBy:     vm.CreatedBy,
		Config:        doc,
		CreatedAt:     vm.CreatedAt,
	}
}

// JSONB <-> domain.ConfigDocument boundary (forward-compatible: empty => zero value).
func toJSON(doc domain.ConfigDocument) (datatypes.JSON, error) {
	b, err := json.Marshal(doc)
	return datatypes.JSON(b), err
}

func toConfigDocument(j datatypes.JSON) (domain.ConfigDocument, error) {
	var doc domain.ConfigDocument
	if len(j) == 0 {
		return doc, nil
	}
	err := json.Unmarshal(j, &doc)
	return doc, err
}
