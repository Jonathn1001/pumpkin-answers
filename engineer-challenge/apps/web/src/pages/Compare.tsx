import { useState } from 'react'
import { Button, Select, Space, Typography } from 'antd'
import { useTenants, useDiff } from '../api/hooks'
import { DiffView } from '../components/DiffView'
import type { Change } from '../api/types'

export function Compare() {
  const tenants = useTenants()
  const diff = useDiff()
  const [left, setLeft] = useState<string>()
  const [right, setRight] = useState<string>()
  const [changes, setChanges] = useState<Change[] | null>(null)
  const opts = (tenants.data ?? []).map((t) => ({ value: t.slug, label: `${t.name} (${t.slug})` }))
  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ margin: 0 }}>
        Compare tenants
      </Typography.Title>
      <Space wrap>
        <Select placeholder="left tenant…" style={{ width: 240 }} value={left} onChange={setLeft} options={opts} />
        <Select placeholder="right tenant…" style={{ width: 240 }} value={right} onChange={setRight} options={opts} />
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
