import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { App, Button, Descriptions, Form, Input, Modal, Result, Steps } from 'antd'
import { useCreateTenant } from '../api/hooks'
import { ApiError } from '../api/client'
import type { Tenant } from '../api/types'

// Multi-step "Add New Tenant" wizard. Always starts from the default config —
// cloning an existing tenant is a separate flow (CloneTenantModal).
export function CreateTenantWizard({ open, onClose }: { open: boolean; onClose: () => void }) {
  const navigate = useNavigate()
  const { message } = App.useApp()
  const create = useCreateTenant()
  const [step, setStep] = useState(0)
  const [name, setName] = useState('')
  const [created, setCreated] = useState<Tenant | null>(null)

  function close() {
    setStep(0)
    setName('')
    setCreated(null)
    create.reset()
    onClose()
  }

  function submit() {
    create.mutate(
      { name },
      {
        onSuccess: setCreated,
        onError: (er) =>
          message.error(
            er instanceof ApiError
              ? er.fields?.map((f) => `${f.field}: ${f.message}`).join('; ') || er.message
              : 'Failed to create tenant',
          ),
      },
    )
  }

  return (
    <Modal open={open} onCancel={close} footer={null} title={created ? null : 'Add New Tenant'}>
      {created ? (
        <Result
          status="success"
          title="Tenant Created Successfully!"
          subTitle={`${created.name} · ${created.slug}`}
          extra={[
            <Button key="go" type="primary" onClick={() => { const s = created.slug; close(); navigate(`/t/${s}`) }}>
              Go to Tenant
            </Button>,
            <Button key="close" onClick={close}>Close</Button>,
          ]}
        />
      ) : (
        <>
          <Steps
            current={step}
            size="small"
            style={{ marginBottom: 24 }}
            items={[{ title: 'Basics' }, { title: 'Review & Create' }]}
          />
          {step === 0 ? (
            <>
              <Form layout="vertical">
                <Form.Item label="Tenant name" required help="The URL slug is generated from the name by the server.">
                  <Input
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    onPressEnter={() => name.trim() && setStep(1)}
                    placeholder="e.g. SafeGuard Insurance"
                    autoFocus
                  />
                </Form.Item>
              </Form>
              <div style={{ textAlign: 'right' }}>
                <Button onClick={close} style={{ marginRight: 8 }}>Cancel</Button>
                <Button type="primary" disabled={!name.trim()} onClick={() => setStep(1)}>Next</Button>
              </div>
            </>
          ) : (
            <>
              <Descriptions column={1} bordered size="small">
                <Descriptions.Item label="Name">{name}</Descriptions.Item>
                <Descriptions.Item label="URL slug">generated from name</Descriptions.Item>
                <Descriptions.Item label="Starts from">Default configuration</Descriptions.Item>
              </Descriptions>
              <div style={{ textAlign: 'right', marginTop: 24 }}>
                <Button onClick={() => setStep(0)} style={{ marginRight: 8 }}>Back</Button>
                <Button type="primary" loading={create.isPending} onClick={submit}>Create Tenant</Button>
              </div>
            </>
          )}
        </>
      )}
    </Modal>
  )
}
