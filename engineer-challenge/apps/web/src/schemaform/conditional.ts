import type { ConfigDocument, FieldDescriptor } from '../api/types'

// Hide descriptors that don't apply for the current config (tiered vs committee).
export function isVisible(d: FieldDescriptor, config: ConfigDocument): boolean {
  if (d.key === 'approval.tiers') return config.approval?.model === 'tiered'
  if (d.key === 'approval.committee') return config.approval?.model === 'committee'
  return true
}

// webhookUrl is required when the webhook channel is enabled.
export function webhookRequired(config: ConfigDocument): boolean {
  return (config.notifications?.channels ?? []).includes('webhook')
}

// Effective required-ness for a field given the current config (server stays source of truth).
export function isRequired(d: FieldDescriptor, config: ConfigDocument): boolean {
  if (d.key === 'notifications.webhookUrl') return webhookRequired(config)
  return d.required ?? false
}
