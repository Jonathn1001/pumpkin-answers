package registry

import "claimsplatform/internal/domain"

// FieldDescriptor là metadata trình bày (presentation) — tách khỏi domain logic.
type FieldDescriptor struct {
	Key      string   `json:"key"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`   // string|number|boolean|select|array|object
	Widget   string   `json:"widget"` // text|number|toggle|select|tier-list|claimtype-grid|events-grid|customfields
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"`
}

// Dimension: phần domain (Evaluate/Validate/DefaultConfig) + presentation (UISchema).
type Dimension interface {
	Key() string
	Evaluate(cfg domain.ConfigDocument, claim domain.Claim, dec *domain.ClaimDecision)
	Validate(cfg domain.ConfigDocument) []domain.FieldError
	DefaultConfig() any
	UISchema() []FieldDescriptor
}

var registered []Dimension

func Register(d Dimension) { registered = append(registered, d) }

func All() []Dimension { return registered }

// Get trả dimension theo key (engine dùng để chạy claimTypes trước).
func Get(key string) Dimension {
	for _, d := range registered {
		if d.Key() == key {
			return d
		}
	}
	return nil
}

// Reset chỉ dùng trong test.
func Reset() { registered = nil }
