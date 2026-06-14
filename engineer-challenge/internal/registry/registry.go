package registry

import "claimsplatform/internal/domain"

// Dimension keys. The engine's ordering and lookups, and each dimension's Key(),
// reference these so a rename or typo is a compile error, not a silent break.
const (
	KeyBranding      = "branding"
	KeyClaimTypes    = "claimTypes"
	KeyApproval      = "approval"
	KeyNotifications = "notifications"
	KeySLA           = "sla"
	KeyCustomFields  = "customFields"
)

// FieldDescriptor is presentation metadata, kept separate from domain logic.
type FieldDescriptor struct {
	Key      string   `json:"key"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`   // string|number|boolean|select|array|object
	Widget   string   `json:"widget"` // text|number|toggle|select|tier-list|claimtype-grid|events-grid|customfields
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"`
}

// Dimension combines domain logic (Evaluate/Validate/DefaultConfig) with presentation metadata (UISchema).
type Dimension interface {
	Key() string
	Evaluate(cfg domain.ConfigDocument, claim domain.Claim, dec *domain.ClaimDecision)
	Validate(cfg domain.ConfigDocument) []domain.FieldError
	DefaultConfig() any
	UISchema() []FieldDescriptor
}

var registered []Dimension

func Register(d Dimension) { registered = append(registered, d) }

func All() []Dimension {
	out := make([]Dimension, len(registered))
	copy(out, registered)
	return out
}

// Get returns the dimension registered under key, or nil (the engine uses it to run claimTypes first).
func Get(key string) Dimension {
	for _, d := range registered {
		if d.Key() == key {
			return d
		}
	}
	return nil
}

// Reset clears the registry; for tests only.
func Reset() { registered = nil }
