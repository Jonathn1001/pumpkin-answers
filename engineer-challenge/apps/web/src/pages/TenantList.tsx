import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Avatar, Button, Input, Select, Space, Table, Tag, Typography } from 'antd'
import type { TableColumnsType } from 'antd'
import { CopyOutlined, EditOutlined, InboxOutlined, PlusOutlined, UndoOutlined } from '@ant-design/icons'
import { useTenants, useUpdateTenantMeta } from '../api/hooks'
import { CreateTenantWizard } from '../components/CreateTenantWizard'
import { CloneTenantModal } from '../components/CloneTenantModal'
import type { Tenant } from '../api/types'

export function TenantList() {
  const { data: tenants, isLoading } = useTenants()
  const update = useUpdateTenantMeta()
  const [createOpen, setCreateOpen] = useState(false)
  const [cloneSource, setCloneSource] = useState<{ slug: string; name: string } | null>(null)
  const [q, setQ] = useState('')
  const [status, setStatus] = useState('all')
  const [editingSlug, setEditingSlug] = useState<string | null>(null)
  const [editName, setEditName] = useState('')

  const data = useMemo(
    () =>
      (tenants ?? []).filter(
        (t) =>
          (status === 'all' || t.status === status) &&
          (q === '' || t.name.toLowerCase().includes(q.toLowerCase()) || t.slug.includes(q.toLowerCase())),
      ),
    [tenants, q, status],
  )

  function saveEdit(t: Tenant) {
    if (!editName.trim()) return
    update.mutate({ slug: t.slug, name: editName, status: t.status }, { onSuccess: () => setEditingSlug(null) })
  }
  function toggleArchive(t: Tenant) {
    update.mutate({ slug: t.slug, name: t.name, status: t.status === 'archived' ? 'active' : 'archived' })
  }

  const columns: TableColumnsType<Tenant> = [
    {
      title: 'Name',
      dataIndex: 'name',
      render: (_, t) =>
        editingSlug === t.slug ? (
          <Space>
            <Input
              size="small"
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
              onPressEnter={() => saveEdit(t)}
              style={{ width: 180 }}
              autoFocus
            />
            <Button size="small" type="link" loading={update.isPending} disabled={!editName.trim()} onClick={() => saveEdit(t)}>
              Save
            </Button>
            <Button size="small" type="link" onClick={() => setEditingSlug(null)}>
              Cancel
            </Button>
          </Space>
        ) : (
          <Space>
            <Avatar style={{ backgroundColor: '#2f54eb', flexShrink: 0 }}>{t.name[0]?.toUpperCase()}</Avatar>
            <Link to={`/t/${t.slug}`}>{t.name}</Link>
          </Space>
        ),
    },
    { title: 'Slug', dataIndex: 'slug', render: (s: string) => <Typography.Text code>{s}</Typography.Text> },
    {
      title: 'Status',
      dataIndex: 'status',
      render: (s: string) => <Tag color={s === 'active' ? 'green' : 'default'}>{s}</Tag>,
    },
    { title: 'Active v', dataIndex: 'activeVersionNumber', render: (v?: number) => v ?? '—' },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, t) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => { setEditingSlug(t.slug); setEditName(t.name) }}>
            Edit
          </Button>
          <Button
            size="small"
            icon={t.status === 'archived' ? <UndoOutlined /> : <InboxOutlined />}
            onClick={() => toggleArchive(t)}
          >
            {t.status === 'archived' ? 'Activate' : 'Archive'}
          </Button>
          <Button size="small" icon={<CopyOutlined />} onClick={() => setCloneSource({ slug: t.slug, name: t.name })}>
            Clone
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography.Title level={3} style={{ margin: 0 }}>
          Tenants
        </Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>
          Create Tenant
        </Button>
      </div>
      <Space wrap>
        <Input.Search placeholder="Search name or slug" allowClear value={q} onChange={(e) => setQ(e.target.value)} style={{ width: 240 }} />
        <Select
          value={status}
          onChange={setStatus}
          style={{ width: 150 }}
          options={[
            { value: 'all', label: 'All status' },
            { value: 'active', label: 'Active' },
            { value: 'archived', label: 'Archived' },
          ]}
        />
      </Space>
      <Table rowKey="slug" loading={isLoading} columns={columns} dataSource={data} pagination={false} />

      <CreateTenantWizard open={createOpen} onClose={() => setCreateOpen(false)} />
      <CloneTenantModal source={cloneSource} onClose={() => setCloneSource(null)} />
    </Space>
  )
}
