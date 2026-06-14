package domain

import "errors"

var (
	ErrTenantNotFound  = errors.New("tenant not found")
	ErrVersionNotFound = errors.New("config version not found")
	ErrSlugTaken       = errors.New("tenant slug already exists")
)

// ValidationError carries field-level config validation failures (maps to HTTP 422).
type ValidationError struct {
	Fields []FieldError
}

func (e ValidationError) Error() string { return "config validation failed" }
