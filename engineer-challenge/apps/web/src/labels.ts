// Human-readable display labels for the machine identifiers that the API uses
// (notification event keys, config dimension keys, claim types).
//
// Display-only: these keys are part of the stored config + engine contract and
// are NEVER renamed in the data — we only prettify them for the UI. Unknown keys
// fall back to a generic humanizer so a new backend key still renders sensibly.

// "claim_submitted" -> "Claim submitted", "claimTypes" -> "Claim types"
export function humanize(key: string): string {
  const spaced = key
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2') // camelCase -> camel Case
    .replace(/[_-]+/g, ' ') // snake_case / kebab-case -> spaces
    .trim()
  return spaced.charAt(0).toUpperCase() + spaced.slice(1)
}

const EVENT_LABELS: Record<string, string> = {
  claim_submitted: 'Claim submitted',
  claim_auto_approved: 'Claim auto-approved',
  claim_routed: 'Claim routed',
  sla_breach_warning: 'SLA breach warning',
}

const DIMENSION_LABELS: Record<string, string> = {
  branding: 'Branding',
  claimTypes: 'Claim Types',
  approval: 'Approval',
  notifications: 'Notifications',
  sla: 'SLA',
  customFields: 'Custom Fields',
}

const CLAIM_TYPE_LABELS: Record<string, string> = {
  OUTPATIENT: 'Outpatient',
  INPATIENT: 'Inpatient',
  DENTAL: 'Dental',
  MATERNITY: 'Maternity',
  OPTICAL: 'Optical',
}

/** Friendly label for a notification event key (e.g. "claim_submitted"). */
export const eventLabel = (key: string): string => EVENT_LABELS[key] ?? humanize(key)

/** Friendly label for a config/trace dimension key (e.g. "claimTypes", "sla"). */
export const dimensionLabel = (key: string): string => DIMENSION_LABELS[key] ?? humanize(key)

/** Friendly label for a claim-type key (e.g. "OUTPATIENT" -> "Outpatient"). */
export const claimTypeLabel = (key: string): string => CLAIM_TYPE_LABELS[key] ?? humanize(key.toLowerCase())

/**
 * Friendly breadcrumb for a dotted diff path like "branding.displayName" or
 * "notifications.events.claim_submitted.0". The leading segment is a config
 * dimension; numeric segments are 1-based list positions; segments nested under
 * "events" are notification event keys; segments keyed by claim type (under
 * "claimTypes"/"perClaimType") are claim-type keys; everything else is humanized.
 */
export function pathLabel(path: string): string {
  const segments = path.split('.')
  return segments
    .map((seg, i) => {
      if (i === 0) return dimensionLabel(seg)
      if (/^\d+$/.test(seg)) return `#${Number(seg) + 1}`
      if (segments[i - 1] === 'events') return eventLabel(seg)
      if (segments[i - 1] === 'claimTypes' || segments[i - 1] === 'perClaimType') return claimTypeLabel(seg)
      return humanize(seg)
    })
    .join(' › ')
}
