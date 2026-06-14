package customfields

import (
	"fmt"
	"regexp"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

type dimension struct{}

func New() registry.Dimension { return dimension{} }

func init() { registry.Register(New()) }

func (dimension) Key() string { return "customFields" }

func (dimension) Evaluate(cfg domain.ConfigDocument, claim domain.Claim, dec *domain.ClaimDecision) {
	if !dec.Accepted {
		return
	}
	result := &domain.CustomFieldValidation{Valid: true}
	for _, f := range cfg.CustomFields {
		raw, present := claim.CustomFields[f.Key]
		if !present || raw == nil || raw == "" {
			if f.Required {
				result.Errors = append(result.Errors, domain.FieldError{Field: f.Key, Message: "required field is missing"})
			}
			continue
		}
		if f.Validation != nil && f.Validation.Pattern != nil {
			s := fmt.Sprintf("%v", raw)
			matched, err := regexp.MatchString(*f.Validation.Pattern, s)
			if err != nil || !matched {
				result.Errors = append(result.Errors, domain.FieldError{Field: f.Key, Message: "value does not match required pattern"})
			}
		}
	}
	result.Valid = len(result.Errors) == 0
	dec.CustomFieldValidation = result
	dec.Trace = append(dec.Trace, domain.TraceEntry{Dimension: "customFields",
		Explanation: fmt.Sprintf("validated %d custom field(s); valid=%v", len(cfg.CustomFields), result.Valid)})
}

func (dimension) Validate(cfg domain.ConfigDocument) []domain.FieldError {
	var errs []domain.FieldError
	seen := map[string]bool{}
	for i, f := range cfg.CustomFields {
		if f.Key == "" {
			errs = append(errs, domain.FieldError{Field: fmt.Sprintf("customFields[%d].key", i), Message: "key is required"})
		}
		if seen[f.Key] {
			errs = append(errs, domain.FieldError{Field: fmt.Sprintf("customFields[%d].key", i), Message: "duplicate key"})
		}
		seen[f.Key] = true
		if f.Type == "select" && len(f.Options) == 0 {
			errs = append(errs, domain.FieldError{Field: fmt.Sprintf("customFields[%d].options", i), Message: "select field needs options"})
		}
	}
	return errs
}

func (dimension) DefaultConfig() any { return []domain.CustomFieldConfig{} }

func (dimension) UISchema() []registry.FieldDescriptor {
	return []registry.FieldDescriptor{{Key: "customFields", Label: "Custom Fields", Type: "array", Widget: "customfields"}}
}
