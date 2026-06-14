import type { Change } from '../api/types'

const tone = { added: 'bg-green-50 text-green-800', removed: 'bg-red-50 text-red-800', changed: 'bg-amber-50 text-amber-800' }

export function DiffView({ changes }: { changes: Change[] }) {
  if (!changes.length) return <div className="text-sm text-gray-500">No differences.</div>
  return (
    <table className="w-full text-sm">
      <thead><tr className="text-left"><th>Path</th><th>Type</th><th>Left</th><th>Right</th></tr></thead>
      <tbody>{changes.map((c, i) => (
        <tr key={i} className={tone[c.type]}>
          <td className="font-mono">{c.path}</td><td>{c.type}</td>
          <td className="font-mono">{JSON.stringify(c.left)}</td><td className="font-mono">{JSON.stringify(c.right)}</td>
        </tr>
      ))}</tbody>
    </table>
  )
}
