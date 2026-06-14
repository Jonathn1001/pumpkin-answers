import { useState } from 'react'
import { Button, Card, Empty, Table, Tag, Typography } from 'antd'
import type { TableColumnsType } from 'antd'
import { useVersions, useRollback, useTenant, useDiff } from '../../api/hooks'
import { DiffView } from '../../components/DiffView'
import type { Change, ConfigVersion } from '../../api/types'

export function VersionsTab({ slug }: { slug: string }) {
  const versions = useVersions(slug)
  const tenant = useTenant(slug)
  const rollback = useRollback(slug)
  const diff = useDiff()
  const [selected, setSelected] = useState<number | null>(null)
  const [changes, setChanges] = useState<Change[] | null>(null)
  const active = tenant.data?.activeVersionNumber

  function select(n: number) {
    setSelected(n)
    setChanges(null)
    if (n > 1) {
      diff.mutate({ left: `${slug}@${n - 1}`, right: `${slug}@${n}` }, { onSuccess: (r) => setChanges(r.changes) })
    }
  }

  const columns: TableColumnsType<ConfigVersion> = [
    {
      title: 'Version',
      dataIndex: 'versionNumber',
      render: (n: number) => (
        <span>
          v{n}
          {n === active ? (
            <Tag color="green" style={{ marginLeft: 6 }}>
              active
            </Tag>
          ) : null}
        </span>
      ),
    },
    { title: 'Status', dataIndex: 'status' },
    { title: 'By', dataIndex: 'createdBy', render: (b?: string) => b || '—' },
    { title: 'When', dataIndex: 'createdAt', render: (d: string) => d.slice(0, 10) },
    { title: 'Note', dataIndex: 'note' },
  ]

  return (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
      <Table
        rowKey="versionNumber"
        size="small"
        loading={versions.isLoading}
        columns={columns}
        dataSource={versions.data}
        pagination={false}
        onRow={(r) => ({ onClick: () => select(r.versionNumber), style: { cursor: 'pointer' } })}
        rowClassName={(r) => (r.versionNumber === selected ? 'ant-table-row-selected' : '')}
      />
      <Card
        size="small"
        title={selected ? `Version details — v${selected}${selected === active ? ' (active)' : ''}` : 'Version details'}
        extra={
          selected && selected !== active ? (
            <Button size="small" type="primary" loading={rollback.isPending} onClick={() => rollback.mutate(selected)}>
              Rollback to this version
            </Button>
          ) : null
        }
      >
        {selected == null ? (
          <Empty description="Select a version to see what changed." image={Empty.PRESENTED_IMAGE_SIMPLE} />
        ) : selected === 1 ? (
          <Typography.Text type="secondary">Initial version — nothing to compare against.</Typography.Text>
        ) : diff.isPending && !changes ? (
          <Typography.Text type="secondary">Loading diff…</Typography.Text>
        ) : changes ? (
          <>
            <Typography.Paragraph type="secondary" style={{ marginBottom: 8 }}>
              Changes from v{selected - 1} → v{selected}:
            </Typography.Paragraph>
            <DiffView changes={changes} />
          </>
        ) : null}
      </Card>
    </div>
  )
}
