package domain

import "time"

type ApprovalOutcome string

const (
	AutoApproved ApprovalOutcome = "auto_approved"
	Routed       ApprovalOutcome = "routed"
)

type ApprovalRoute struct {
	Model             ApprovalModel `json:"model"` // tiered | committee
	TierLabel         string        `json:"tierLabel,omitempty"`
	ApproverRole      string        `json:"approverRole,omitempty"`
	CommitteeName     string        `json:"committeeName,omitempty"`
	RequiredApprovals int           `json:"requiredApprovals,omitempty"`
}

type ApprovalDecision struct {
	Outcome ApprovalOutcome `json:"outcome"`
	Route   *ApprovalRoute  `json:"route,omitempty"`
}

type NotificationFire struct {
	Event    string   `json:"event"`
	Channels []string `json:"channels"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CustomFieldValidation struct {
	Valid  bool         `json:"valid"`
	Errors []FieldError `json:"errors"`
}

type TraceEntry struct {
	Dimension   string `json:"dimension"`
	Explanation string `json:"explanation"`
}

type ClaimDecision struct {
	Accepted              bool                   `json:"accepted"`
	RejectionReasons      []string               `json:"rejectionReasons"`
	RequiredDocuments     []string               `json:"requiredDocuments"`
	Approval              *ApprovalDecision      `json:"approval,omitempty"`
	Notifications         []NotificationFire     `json:"notifications"`
	SLADeadline           *time.Time             `json:"slaDeadline,omitempty"`
	SLADays               int                    `json:"slaDays,omitempty"`
	Escalation            *Escalation            `json:"escalation,omitempty"`
	CustomFieldValidation *CustomFieldValidation `json:"customFieldValidation,omitempty"`
	Trace                 []TraceEntry           `json:"trace"`
	Extras                map[string]any         `json:"extras,omitempty"` // reserved for future dimensions
}
