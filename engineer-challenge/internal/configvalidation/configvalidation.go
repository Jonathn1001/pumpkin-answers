// Package configvalidation is the single source of truth for config-document validity:
// it composes every registered dimension's Validate into one field-error list.
package configvalidation

import (
	_ "claimsplatform/internal/dimensions" // blank import: register all six dimensions
	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

func Validate(cfg domain.ConfigDocument) []domain.FieldError {
	var errs []domain.FieldError
	for _, d := range registry.All() {
		errs = append(errs, d.Validate(cfg)...)
	}
	return errs
}
