import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useCreateTenant, useActiveConfig } from '../api/hooks'
import { ApiError } from '../api/client'
import type { Tenant } from '../api/types'

// ClonePreview is a child so useActiveConfig only runs when a source is chosen
// (avoids a conditional hook / an empty-slug request).
function ClonePreview({ slug }: { slug: string }) {
  const { data, isLoading } = useActiveConfig(slug)
  if (isLoading) return <p className="text-sm text-gray-500">Loading source config…</p>
  if (!data) return null
  const enabled = Object.entries(data.claimTypes)
    .filter(([, c]) => c.enabled)
    .map(([k]) => k)
  const tiers = data.approval.model === 'tiered' ? (data.approval.tiers?.length ?? 0) : 0
  return (
    <ul className="space-y-0.5 text-sm text-gray-700">
      <li>Claim types: {enabled.join(', ') || '—'}</li>
      <li>
        Approval: {data.approval.model}
        {data.approval.model === 'tiered' ? ` (${tiers} tiers)` : ''}
      </li>
      <li>Auto-approve threshold: {data.approval.autoApproveThreshold}</li>
    </ul>
  )
}

export function CreateTenantWizard({
  tenants,
  initialClone,
  onClose,
}: {
  tenants: Tenant[]
  initialClone: string
  onClose: () => void
}) {
  const create = useCreateTenant()
  const [step, setStep] = useState<1 | 2>(1)
  const [name, setName] = useState('')
  const [cloneFrom, setCloneFrom] = useState(initialClone)
  const [err, setErr] = useState('')
  const [created, setCreated] = useState<Tenant | null>(null)

  function submit() {
    setErr('')
    create.mutate(
      { name, cloneFrom: cloneFrom || undefined },
      {
        onSuccess: (t) => setCreated(t),
        onError: (er) =>
          setErr(
            er instanceof ApiError
              ? er.fields?.map((f) => `${f.field}: ${f.message}`).join('; ') || er.message
              : 'error',
          ),
      },
    )
  }

  return (
    <div className="fixed inset-0 z-10 flex items-center justify-center bg-black/30 p-4" onClick={onClose}>
      <div className="w-full max-w-lg rounded-lg border bg-white p-6 shadow-lg" onClick={(e) => e.stopPropagation()}>
        {created ? (
          <div className="space-y-4 text-center">
            <div className="mx-auto flex h-10 w-10 items-center justify-center rounded-full bg-green-100 text-xl text-green-700">✓</div>
            <h3 className="text-lg font-semibold">Tenant Created Successfully!</h3>
            <p className="text-sm text-gray-600">
              {created.name} — <span className="font-mono">{created.slug}</span>
            </p>
            <div className="flex justify-center gap-2">
              <Link to={`/t/${created.slug}`} className="rounded bg-blue-600 px-4 py-1.5 text-white" onClick={onClose}>
                Go to Tenant
              </Link>
              <button className="rounded border px-4 py-1.5" onClick={onClose}>Close</button>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="font-semibold">Create tenant — step {step} of 2</h3>
              <button className="text-gray-400 hover:text-gray-600" onClick={onClose} aria-label="Close">✕</button>
            </div>

            {step === 1 ? (
              <div className="space-y-3">
                <label className="block text-sm">
                  <span className="text-gray-600">Name</span>
                  <input
                    className="mt-1 w-full rounded border px-2 py-1"
                    placeholder="e.g. SafeGuard Insurance"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    autoFocus
                  />
                  <span className="text-xs text-gray-500">The URL slug is generated from the name by the server.</span>
                </label>
                <label className="block text-sm">
                  <span className="text-gray-600">Start from</span>
                  <select className="mt-1 w-full rounded border px-2 py-1" value={cloneFrom} onChange={(e) => setCloneFrom(e.target.value)}>
                    <option value="">Default configuration</option>
                    {tenants.map((t) => (
                      <option key={t.slug} value={t.slug}>
                        Clone: {t.name} ({t.slug})
                      </option>
                    ))}
                  </select>
                </label>
                <div className="flex justify-end gap-2">
                  <button className="rounded border px-4 py-1.5" onClick={onClose}>Cancel</button>
                  <button
                    className="rounded bg-blue-600 px-4 py-1.5 text-white disabled:opacity-50"
                    disabled={!name.trim()}
                    onClick={() => setStep(2)}
                  >
                    Next
                  </button>
                </div>
              </div>
            ) : (
              <div className="space-y-3">
                <div className="rounded border bg-gray-50 p-3 text-sm">
                  <p className="mb-2 font-medium">Review</p>
                  <dl className="space-y-1">
                    <div className="flex justify-between gap-4">
                      <dt className="text-gray-500">Name</dt>
                      <dd className="text-right">{name}</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-gray-500">URL slug</dt>
                      <dd className="text-right text-gray-500">generated from name</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-gray-500">Source</dt>
                      <dd className="text-right">{cloneFrom ? `Cloned from ${cloneFrom}` : 'Default configuration'}</dd>
                    </div>
                  </dl>
                  {cloneFrom && (
                    <div className="mt-2 border-t pt-2">
                      <ClonePreview slug={cloneFrom} />
                    </div>
                  )}
                </div>
                {err && <div className="text-sm text-red-600">{err}</div>}
                <div className="flex justify-between">
                  <button className="rounded border px-4 py-1.5" onClick={() => setStep(1)}>Back</button>
                  <button
                    className="rounded bg-blue-600 px-4 py-1.5 text-white disabled:opacity-50"
                    disabled={create.isPending}
                    onClick={submit}
                  >
                    Create Tenant
                  </button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
