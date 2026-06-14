package domain

import (
	"strings"
	"time"
)

type ClaimType string

const (
	Outpatient ClaimType = "OUTPATIENT"
	Inpatient  ClaimType = "INPATIENT"
	Dental     ClaimType = "DENTAL"
	Maternity  ClaimType = "MATERNITY"
	Optical    ClaimType = "OPTICAL"
)

// AllClaimTypes returns the fixed claim-type order used for rendering/seeding (a fresh slice each call so callers cannot mutate it).
func AllClaimTypes() []ClaimType {
	return []ClaimType{Outpatient, Inpatient, Dental, Maternity, Optical}
}

type Claim struct {
	Type         ClaimType      `json:"type"`
	Amount       int64          `json:"amount"`
	SubmittedAt  time.Time      `json:"submittedAt"`
	CustomFields map[string]any `json:"customFields"`
}

// Valid reports whether t is one of the known claim types.
func (t ClaimType) Valid() bool {
	for _, v := range AllClaimTypes() {
		if t == v {
			return true
		}
	}
	return false
}

// Validate checks claim input shape (required type from the known set, non-negative
// amount, present submission time). Whether a valid claim is *accepted* is the engine's
// business decision; this only rejects structurally bad input (maps to HTTP 422).
func (c Claim) Validate() []FieldError {
	var errs []FieldError
	switch {
	case c.Type == "":
		errs = append(errs, FieldError{Field: "type", Message: "is required"})
	case !c.Type.Valid():
		errs = append(errs, FieldError{Field: "type", Message: "must be one of " + claimTypeList()})
	}
	if c.Amount < 0 {
		errs = append(errs, FieldError{Field: "amount", Message: "must not be negative"})
	}
	if c.SubmittedAt.IsZero() {
		errs = append(errs, FieldError{Field: "submittedAt", Message: "is required"})
	}
	return errs
}

func claimTypeList() string {
	ts := AllClaimTypes()
	ss := make([]string, len(ts))
	for i, t := range ts {
		ss[i] = string(t)
	}
	return strings.Join(ss, ", ")
}
