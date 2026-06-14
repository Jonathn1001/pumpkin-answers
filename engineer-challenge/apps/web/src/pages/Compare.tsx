import { useState } from 'react'
import { useTenants, useDiff } from '../api/hooks'
import { DiffView } from '../components/DiffView'
import type { Change } from '../api/types'

export function Compare() {
  const tenants = useTenants()
  const diff = useDiff()
  const [left, setLeft] = useState('')
  const [right, setRight] = useState('')
  const [changes, setChanges] = useState<Change[] | null>(null)
  return (
    <div className="space-y-3">
      <h2 className="text-xl font-semibold">Compare tenants</h2>
      <div className="flex gap-2">
        <select className="rounded border px-2 py-1" value={left} onChange={e => setLeft(e.target.value)}><option value="">left…</option>{tenants.data?.map(t => <option key={t.slug} value={t.slug}>{t.slug}</option>)}</select>
        <select className="rounded border px-2 py-1" value={right} onChange={e => setRight(e.target.value)}><option value="">right…</option>{tenants.data?.map(t => <option key={t.slug} value={t.slug}>{t.slug}</option>)}</select>
        <button className="rounded bg-blue-600 px-3 py-1 text-white" disabled={!left || !right} onClick={() => diff.mutate({ left, right }, { onSuccess: r => setChanges(r.changes) })}>Diff</button>
      </div>
      {changes && <DiffView changes={changes} />}
    </div>
  )
}
