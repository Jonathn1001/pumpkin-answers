package notifications

import (
	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

type dimension struct{}

func New() registry.Dimension { return dimension{} }

func init() { registry.Register(New()) }

func (dimension) Key() string { return "notifications" }

func (dimension) Evaluate(cfg domain.ConfigDocument, _ domain.Claim, dec *domain.ClaimDecision) {
	if !dec.Accepted {
		return
	}
	events := []string{"claim_submitted"}
	if dec.Approval != nil {
		switch dec.Approval.Outcome {
		case domain.AutoApproved:
			events = append(events, "claim_auto_approved")
		case domain.Routed:
			events = append(events, "claim_routed")
		}
	}
	for _, ev := range events {
		channels, ok := cfg.Notifications.Events[ev]
		if !ok {
			channels = cfg.Notifications.Channels
		}
		dec.Notifications = append(dec.Notifications, domain.NotificationFire{Event: ev, Channels: channels})
	}
	dec.Trace = append(dec.Trace, domain.TraceEntry{Dimension: "notifications",
		Explanation: "fired events: " + joinEvents(events)})
}

func joinEvents(ev []string) string {
	out := ""
	for i, e := range ev {
		if i > 0 {
			out += ", "
		}
		out += e
	}
	return out
}

func (dimension) Validate(cfg domain.ConfigDocument) []domain.FieldError {
	var errs []domain.FieldError
	hasWebhook := false
	for _, c := range cfg.Notifications.Channels {
		if c == "webhook" {
			hasWebhook = true
		}
	}
	if hasWebhook && cfg.Notifications.WebhookURL == "" {
		errs = append(errs, domain.FieldError{Field: "notifications.webhookUrl", Message: "webhook channel enabled but webhookUrl is empty"})
	}
	return errs
}

func (dimension) DefaultConfig() any {
	return domain.NotificationsConfig{Channels: []string{"email"},
		Events: map[string][]string{"claim_submitted": {"email"}}}
}

func (dimension) UISchema() []registry.FieldDescriptor {
	return []registry.FieldDescriptor{
		{Key: "notifications.channels", Label: "Channels", Type: "array", Widget: "channel-multiselect", Options: []string{"email", "sms", "webhook"}},
		{Key: "notifications.events", Label: "Event -> channel map", Type: "object", Widget: "events-grid"},
		{Key: "notifications.webhookUrl", Label: "Webhook URL", Type: "string", Widget: "text"},
	}
}
