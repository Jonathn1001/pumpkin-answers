package domain

type ConfigDocument struct {
	Branding      BrandingConfig                `json:"branding"`
	ClaimTypes    map[ClaimType]ClaimTypeConfig `json:"claimTypes"`
	Approval      ApprovalConfig                `json:"approval"`
	Notifications NotificationsConfig           `json:"notifications"`
	SLA           SLAConfig                     `json:"sla"`
	CustomFields  []CustomFieldConfig           `json:"customFields"`
}

type BrandingConfig struct {
	DisplayName    string `json:"displayName"`
	LogoURL        string `json:"logoUrl"`
	PrimaryColor   string `json:"primaryColor"`
	SecondaryColor string `json:"secondaryColor"`
	SupportEmail   string `json:"supportEmail"`
}

type ClaimTypeConfig struct {
	Enabled           bool     `json:"enabled"`
	RequiredDocuments []string `json:"requiredDocuments"`
}

type ApprovalTier struct {
	Label        string `json:"label"`
	MaxAmount    *int64 `json:"maxAmount"` // nil = open-ended
	ApproverRole string `json:"approverRole"`
}

type Committee struct {
	Name              string `json:"name"`
	RequiredApprovals int    `json:"requiredApprovals"`
}

type ApprovalModel string

const (
	ApprovalModelTiered    ApprovalModel = "tiered"
	ApprovalModelCommittee ApprovalModel = "committee"
)

type ApprovalConfig struct {
	AutoApproveThreshold int64          `json:"autoApproveThreshold"`
	Model                ApprovalModel  `json:"model"` // tiered | committee
	Tiers                []ApprovalTier `json:"tiers,omitempty"`
	Committee            *Committee     `json:"committee,omitempty"`
}

type NotificationsConfig struct {
	Channels   []string            `json:"channels"`
	Events     map[string][]string `json:"events"`
	WebhookURL string              `json:"webhookUrl,omitempty"`
}

type Escalation struct {
	WarnBeforeDays int    `json:"warnBeforeDays"`
	NotifyRole     string `json:"notifyRole"`
}

type SLAConfig struct {
	DefaultDays  int               `json:"defaultDays"`
	PerClaimType map[ClaimType]int `json:"perClaimType,omitempty"`
	Escalation   Escalation        `json:"escalation"`
}

type FieldValidation struct {
	Pattern *string  `json:"pattern,omitempty"`
	Min     *float64 `json:"min,omitempty"`
	Max     *float64 `json:"max,omitempty"`
}

type CustomFieldConfig struct {
	Key        string           `json:"key"`
	Label      string           `json:"label"`
	Type       string           `json:"type"` // string|number|date|select|boolean
	Required   bool             `json:"required"`
	Options    []string         `json:"options,omitempty"`
	Validation *FieldValidation `json:"validation,omitempty"`
}
