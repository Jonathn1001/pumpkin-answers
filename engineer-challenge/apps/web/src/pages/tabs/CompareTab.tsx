import { useState } from 'react'
import { Button, Select, Space } from 'antd'
import { useDiff, useVersions } from '../../api/hooks'
import { DiffView } from '../../components/DiffView'
import type { Change } from '../../api/types'

export function CompareTab({ slug }: { slug: string }) {
  const versions = useVersions(slug)
  const diff = useDiff()
  const [left, setLeft] = useState<string>()
  const [right, setRight] = useState<string>()
  const [changes, setChanges] = useState<Change[] | null>(null)
  const opts = (versions.data ?? []).map((v) => ({
    value: `${slug}@${v.versionNumber}`,
    label: `${slug}@${v.versionNumber}`,
  }))
  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Space wrap>
        <Select placeholder="left…" style={{ width: 200 }} value={left} onChange={setLeft} options={opts} />
        <Select placeholder="right…" style={{ width: 200 }} value={right} onChange={setRight} options={opts} />
        <Button
          type="primary"
          disabled={!left || !right}
          loading={diff.isPending}
          onClick={() => left && right && diff.mutate({ left, right }, { onSuccess: (r) => setChanges(r.changes) })}
        >
          Diff
        </Button>
      </Space>
      {changes && <DiffView changes={changes} />}
    </Space>
  )
}
