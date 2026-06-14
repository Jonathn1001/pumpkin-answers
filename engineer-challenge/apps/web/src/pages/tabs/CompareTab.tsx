import { useState } from 'react'
import { useDiff, useVersions } from '../../api/hooks'
import { DiffView } from '../../components/DiffView'
import type { Change } from '../../api/types'

export function CompareTab({ slug }: { slug: string }) {
  const versions = useVersions(slug)
  const diff = useDiff()
  const [left, setLeft] = useState('')
  const [right, setRight] = useState('')
  const [changes, setChanges] = useState<Change[] | null>(null)
  const opts = versions.data?.map(v => `${slug}@${v.versionNumber}`) ?? []
  return (
    <div className="space-y-3">
      <div className="flex gap-2">
        <select className="rounded border px-2 py-1" value={left} onChange={e => setLeft(e.target.value)}><option value="">left…</option>{opts.map(o => <option key={o}>{o}</option>)}</select>
        <select className="rounded border px-2 py-1" value={right} onChange={e => setRight(e.target.value)}><option value="">right…</option>{opts.map(o => <option key={o}>{o}</option>)}</select>
        <button className="rounded bg-blue-600 px-3 py-1 text-white" disabled={!left || !right} onClick={() => diff.mutate({ left, right }, { onSuccess: r => setChanges(r.changes) })}>Diff</button>
      </div>
      {changes && <DiffView changes={changes} />}
    </div>
  )
}
