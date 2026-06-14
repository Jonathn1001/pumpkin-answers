import { useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { App, Avatar, Button, Form, Input, Modal, Segmented, Space, Table, Tag, Tooltip, Typography } from 'antd'
import type { TableColumnsType } from 'antd'
import {
  ClockCircleOutlined,
  DeleteOutlined,
  EditOutlined,
  EyeOutlined,
  PlusOutlined,
  UndoOutlined,
} from '@ant-design/icons'
import { useTenants, useUpdateTenantMeta } from '../api/hooks'
import type { Tenant } from '../api/types'

export function TenantList() {
  const { data: tenants, isLoading } = useTenants()
  const update = useUpdateTenantMeta()
  const { modal, message } = App.useApp()
  const navigate = useNavigate()
  const [q, setQ] = useState('')
  // Soft-deleted (archived) tenants are hidden from the default view; the
  // Archived view is where they can be restored.
  const [view, setView] = useState<'active' | 'archived'>('active')
  const [renameTarget, setRenameTarget] = useState<Tenant | null>(null)
  const [renameName, setRenameName] = useState('')

  const data = useMemo(
    () =>
      (tenants ?? []).filter(
        (t) => t.status === view && (q === '' || t.name.toLowerCase().includes(q.toLowerCase())),
      ),
    [tenants, q, view],
  )

  function openRename(t: Tenant) {
    setRenameTarget(t)
    setRenameName(t.name)
  }
  function saveRename() {
    if (!renameTarget || !renameName.trim()) return
    update.mutate(
      { slug: renameTarget.slug, name: renameName, status: renameTarget.status },
      {
        onSuccess: () => {
          setRenameTarget(null)
          message.success('Tenant renamed')
        },
      },
    )
  }
  // Delete = soft-delete (archive). Only once the server confirms do we toast and
  // let the table reload — the mutation invalidates the tenants query on success.
  function confirmDelete(t: Tenant) {
    modal.confirm({
      title: `Delete ${t.name}?`,
      content: 'The tenant is archived and hidden from the list. You can restore it from the Archived view.',
      okText: 'Delete',
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          await update.mutateAsync({ slug: t.slug, name: t.name, status: 'archived' })
          message.success(`${t.name} deleted`)
        } catch {
          message.error(`Failed to delete ${t.name}`)
        }
      },
    })
  }
  function restore(t: Tenant) {
    update.mutate({ slug: t.slug, name: t.name, status: 'active' }, { onSuccess: () => message.success('Tenant restored') })
  }

  const columns: TableColumnsType<Tenant> = [
    {
      title: 'Name',
      dataIndex: 'name',
      render: (_, t) => (
        <Space>
          <Avatar style={{ backgroundColor: '#2f54eb', flexShrink: 0 }}>{t.name[0]?.toUpperCase()}</Avatar>
          <Link to={`/t/${t.slug}`}>{t.name}</Link>
        </Space>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: (s: string) => <Tag color={s === 'active' ? 'green' : 'default'}>{s}</Tag>,
    },
    { title: 'Active v', dataIndex: 'activeVersionNumber', render: (v?: number) => v ?? '—' },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, t) =>
        t.status === 'archived' ? (
          <Button size="small" icon={<UndoOutlined />} onClick={() => restore(t)}>
            Restore
          </Button>
        ) : (
          <Space>
            <Tooltip title="Edit name">
              <Button size="small" aria-label="Edit name" icon={<EditOutlined />} onClick={() => openRename(t)} />
            </Tooltip>
            <Tooltip title="Preview">
              <Button size="small" aria-label="Preview" icon={<EyeOutlined />} onClick={() => navigate(`/t/${t.slug}?tab=preview`)} />
            </Tooltip>
            <Tooltip title="Version history">
              <Button size="small" aria-label="Version history" icon={<ClockCircleOutlined />} onClick={() => navigate(`/t/${t.slug}?tab=versions`)} />
            </Tooltip>
            <Tooltip title="Delete">
              <Button size="small" danger aria-label="Delete" icon={<DeleteOutlined />} onClick={() => confirmDelete(t)} />
            </Tooltip>
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
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/tenants/new')}>
          Create Tenant
        </Button>
      </div>
      <Space wrap>
        <Input.Search
          placeholder="Search by name"
          allowClear
          value={q}
          onChange={(e) => setQ(e.target.value)}
          style={{ width: 240 }}
        />
        <Segmented
          value={view}
          onChange={(v) => setView(v as 'active' | 'archived')}
          options={[
            { label: 'Active', value: 'active' },
            { label: 'Archived', value: 'archived' },
          ]}
        />
      </Space>
      <Table rowKey="slug" loading={isLoading} columns={columns} dataSource={data} pagination={false} />

      <Modal
        open={!!renameTarget}
        title={renameTarget ? `Rename "${renameTarget.name}"` : ''}
        onCancel={() => setRenameTarget(null)}
        onOk={saveRename}
        okText="Save"
        confirmLoading={update.isPending}
        okButtonProps={{ disabled: !renameName.trim() }}
      >
        <Form layout="vertical">
          <Form.Item label="Tenant name" required>
            <Input value={renameName} onChange={(e) => setRenameName(e.target.value)} onPressEnter={saveRename} autoFocus />
          </Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}
