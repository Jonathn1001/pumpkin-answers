import { useState } from 'react'
import { App, Form, Input, Modal, Typography } from 'antd'
import { useCreateTenant } from '../api/hooks'
import { ApiError } from '../api/client'

type Source = { slug: string; name: string }

// Quick "Duplicate" of an existing tenant: one field (new name), copies the
// source's active config. Distinct from the multi-step create wizard.
export function CloneTenantModal({ source, onClose }: { source: Source | null; onClose: () => void }) {
  const { message } = App.useApp()
  const create = useCreateTenant()
  const [name, setName] = useState('')
  const [seenSlug, setSeenSlug] = useState<string | null>(null)

  // Prefill the name when a new source opens — adjust state during render rather
  // than in an effect (see react.dev "you might not need an effect").
  if (source && source.slug !== seenSlug) {
    setSeenSlug(source.slug)
    setName(`${source.name} copy`)
  }

  function submit() {
    if (!source || !name.trim()) return
    create.mutate(
      { name, cloneFrom: source.slug },
      {
        onSuccess: (t) => {
          message.success(`Created ${t.name} (${t.slug}) from ${source.slug}`)
          onClose()
        },
        onError: (er) =>
          message.error(
            er instanceof ApiError
              ? er.fields?.map((f) => `${f.field}: ${f.message}`).join('; ') || er.message
              : 'Failed to duplicate tenant',
          ),
      },
    )
  }

  return (
    <Modal
      open={!!source}
      title={source ? `Duplicate "${source.name}"` : ''}
      onCancel={onClose}
      onOk={submit}
      okText="Duplicate"
      confirmLoading={create.isPending}
      okButtonProps={{ disabled: !name.trim() }}
    >
      <Typography.Paragraph type="secondary">
        Copies the active configuration of <Typography.Text code>{source?.slug}</Typography.Text> into a new tenant.
      </Typography.Paragraph>
      <Form layout="vertical">
        <Form.Item label="New tenant name" required help="The slug is generated from the name by the server.">
          <Input value={name} onChange={(e) => setName(e.target.value)} onPressEnter={submit} autoFocus />
        </Form.Item>
      </Form>
    </Modal>
  )
}
