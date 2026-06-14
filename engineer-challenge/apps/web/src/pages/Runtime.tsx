import { useState } from 'react'
import { useTenants, useProcess } from '../api/hooks'
import { ClaimForm } from '../components/ClaimForm'
import { DecisionView } from '../components/DecisionView'
import type { Claim, ClaimDecision } from '../api/types'

type Run = { id: number; slug: string; claim: Claim; decision: ClaimDecision }

function outcomeLabel(d: ClaimDecision): string {
  if (!d.accepted) return 'rejected'
  if (d.approval?.outcome === 'auto_approved') return 'auto-approved'
  if (d.approval?.outcome === 'routed') {
    const r = d.approval.route
    return `routed → ${r?.committeeName ?? r?.tierLabel ?? r?.approverRole ?? ''}`
  }
  return 'accepted'
}

// Runtime runs sample claims against a tenant's active config. Live-only: results
// are kept in session state, not persisted.
export function Runtime() {
  const { data: tenants } = useTenants()
  const [slug, setSlug] = useState('')
  const process = useProcess(slug)
  const [runs, setRuns] = useState<Run[]>([])
  const [selected, setSelected] = useState<number | null>(null)
  const [nextId, setNextId] = useState(1)

  function run(claim: Claim) {
    process.mutate(claim, {
      onSuccess: (decision) => {
        const id = nextId
        setNextId(id + 1)
        setRuns((r) => [{ id, slug, claim, decision }, ...r])
        setSelected(id)
      },
    })
  }

  const detail = runs.find((r) => r.id === selected)

  return (
    <div className="space-y-6">
      <h2 className="text-lg font-semibold">Runtime — process a claim</h2>
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-3 rounded-lg border bg-white p-4">
          <label className="block text-sm">
            <span className="text-gray-600">Tenant</span>
            <select className="mt-1 w-full rounded border px-2 py-1" value={slug} onChange={(e) => setSlug(e.target.value)}>
              <option value="">Select a tenant…</option>
              {tenants?.map((t) => (
                <option key={t.slug} value={t.slug}>{t.name} ({t.slug})</option>
              ))}
            </select>
          </label>
          {slug ? (
            <ClaimForm onSubmit={run} />
          ) : (
            <p className="text-sm text-gray-500">Pick a tenant to run a claim against its active config.</p>
          )}
          {process.isError && <p className="text-sm text-red-600">Failed to process claim.</p>}
        </div>

        <div className="rounded-lg border bg-white p-4">
          {detail ? (
            <div className="space-y-2">
              <p className="text-sm text-gray-500">
                {detail.slug} · {detail.claim.type} · {detail.claim.amount}
              </p>
              <DecisionView d={detail.decision} />
            </div>
          ) : (
            <p className="text-sm text-gray-500">Run a claim to see the processing result.</p>
          )}
        </div>
      </div>

      {runs.length > 0 && (
        <div className="overflow-hidden rounded-lg border bg-white">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b text-left">
                <th className="p-2">#</th><th>Tenant</th><th>Type</th><th>Amount</th><th>Outcome</th>
              </tr>
            </thead>
            <tbody>
              {runs.map((r) => (
                <tr
                  key={r.id}
                  className={`cursor-pointer border-b hover:bg-gray-50 ${r.id === selected ? 'bg-blue-50' : ''}`}
                  onClick={() => setSelected(r.id)}
                >
                  <td className="p-2">{r.id}</td>
                  <td>{r.slug}</td>
                  <td>{r.claim.type}</td>
                  <td>{r.claim.amount}</td>
                  <td>{outcomeLabel(r.decision)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
