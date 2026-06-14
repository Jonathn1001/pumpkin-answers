import { useState } from 'react'
import { Card, Empty, Select, Space, Table, Tag, Typography } from 'antd'
import { useTenants, useProcess } from '../api/hooks'
import { ClaimForm } from '../components/ClaimForm'
import { DecisionView } from '../components/DecisionView'
import type { Claim, ClaimDecision } from '../api/types'

type Run = { id: number; slug: string; claim: Claim; decision: ClaimDecision }

function outcomeTag(d: ClaimDecision) {
  if (!d.accepted) return <Tag color="red">rejected</Tag>
  if (d.approval?.outcome === 'auto_approved') return <Tag color="green">auto-approved</Tag>
  const r = d.approval?.route
  return <Tag color="blue">routed → {r?.committeeName ?? r?.tierLabel ?? r?.approverRole ?? ''}</Tag>
}

// Runtime runs sample claims against a tenant's active config. Live-only: results
// are kept in session state, not persisted.
export function Runtime() {
  const { data: tenants } = useTenants()
  const [slug, setSlug] = useState<string>()
  const process = useProcess(slug ?? '')
  const [runs, setRuns] = useState<Run[]>([])
  const [selected, setSelected] = useState<number | null>(null)
  const [nextId, setNextId] = useState(1)

  function run(claim: Claim) {
    if (!slug) return
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
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ margin: 0 }}>
        Runtime — process a claim
      </Typography.Title>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
        <Card size="small" title="Claim input">
          <Space direction="vertical" style={{ width: '100%' }}>
            <Select
              placeholder="Select a tenant…"
              style={{ width: '100%' }}
              value={slug}
              onChange={setSlug}
              options={(tenants ?? [])
                .filter((t) => t.status === 'active')
                .map((t) => ({ value: t.slug, label: `${t.name} (${t.slug})` }))}
            />
            {slug ? (
              <ClaimForm onSubmit={run} submitText="Process claim" />
            ) : (
              <Typography.Text type="secondary">Pick a tenant to run a claim against its active config.</Typography.Text>
            )}
            {process.isError && <Typography.Text type="danger">Failed to process claim.</Typography.Text>}
          </Space>
        </Card>
        <Card size="small" title="Processing result">
          {detail ? (
            <Space direction="vertical" style={{ width: '100%' }}>
              <Typography.Text type="secondary">
                {detail.slug} · {detail.claim.type} · {detail.claim.amount}
              </Typography.Text>
              <DecisionView d={detail.decision} />
            </Space>
          ) : (
            <Empty description="Run a claim to see the result." image={Empty.PRESENTED_IMAGE_SIMPLE} />
          )}
        </Card>
      </div>
      {runs.length > 0 && (
        <Card size="small" title="Processing log (this session)">
          <Table
            rowKey="id"
            size="small"
            pagination={false}
            dataSource={runs}
            onRow={(r) => ({ onClick: () => setSelected(r.id), style: { cursor: 'pointer' } })}
            rowClassName={(r) => (r.id === selected ? 'ant-table-row-selected' : '')}
            columns={[
              { title: '#', dataIndex: 'id', width: 60 },
              { title: 'Tenant', dataIndex: 'slug' },
              { title: 'Type', render: (_, r) => r.claim.type },
              { title: 'Amount', render: (_, r) => r.claim.amount },
              { title: 'Outcome', render: (_, r) => outcomeTag(r.decision) },
            ]}
          />
        </Card>
      )}
    </Space>
  )
}
