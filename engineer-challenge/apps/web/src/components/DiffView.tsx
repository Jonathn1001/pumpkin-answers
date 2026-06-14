import { Table, Tag, Typography } from 'antd'
import type { Change } from '../api/types'
import { pathLabel } from '../labels'

const tone: Record<Change['type'], string> = { added: 'green', removed: 'red', changed: 'gold' }

export function DiffView({ changes }: { changes: Change[] }) {
  if (!changes.length) return <Typography.Text type="secondary">No differences.</Typography.Text>
  const rows = changes.map((c, i) => ({ ...c, key: i }))
  return (
    <Table
      size="small"
      pagination={false}
      dataSource={rows}
      columns={[
        { title: 'Path', dataIndex: 'path', render: (p: string) => <Typography.Text>{pathLabel(p)}</Typography.Text> },
        { title: 'Type', dataIndex: 'type', render: (t: Change['type']) => <Tag color={tone[t]}>{t}</Tag> },
        { title: 'Left', dataIndex: 'left', render: (v: unknown) => <Typography.Text code>{JSON.stringify(v) ?? '—'}</Typography.Text> },
        { title: 'Right', dataIndex: 'right', render: (v: unknown) => <Typography.Text code>{JSON.stringify(v) ?? '—'}</Typography.Text> },
      ]}
    />
  )
}
