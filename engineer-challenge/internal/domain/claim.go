package domain

import "time"

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
