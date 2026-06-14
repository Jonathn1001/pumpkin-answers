import { useState } from 'react'
import { Button, DatePicker, Form, Input, InputNumber, Select, Space } from 'antd'
import dayjs, { type Dayjs } from 'dayjs'
import type { Claim, ClaimType } from '../api/types'

const TYPES: ClaimType[] = ['OUTPATIENT', 'INPATIENT', 'DENTAL', 'MATERNITY', 'OPTICAL']
const DEFAULT_CUSTOM =
  '{"employeeId":"EMP1234","policyNumber":"HF-12345678","memberTier":"Gold","nationalId":"123456789012","citizenCategory":"General"}'

export function ClaimForm({
  onSubmit,
  submitText = 'Run preview',
}: {
  onSubmit: (c: Claim) => void
  submitText?: string
}) {
  const [type, setType] = useState<ClaimType>('OUTPATIENT')
  const [amount, setAmount] = useState<number>(10000)
  const [submittedAt, setSubmittedAt] = useState<Dayjs>(dayjs('2026-06-14'))
  const [custom, setCustom] = useState(DEFAULT_CUSTOM)

  function submit() {
    let cf: Record<string, unknown> = {}
    try {
      cf = JSON.parse(custom)
    } catch {
      /* ignore invalid json */
    }
    onSubmit({ type, amount, submittedAt: submittedAt.toISOString(), customFields: cf })
  }

  return (
    <Form layout="vertical" onFinish={submit}>
      <Space wrap align="end">
        <Form.Item label="Claim type" style={{ marginBottom: 12 }}>
          <Select value={type} onChange={setType} style={{ width: 160 }} options={TYPES.map((t) => ({ value: t, label: t }))} />
        </Form.Item>
        <Form.Item label="Amount" style={{ marginBottom: 12 }}>
          <InputNumber value={amount} onChange={(v) => setAmount(v ?? 0)} style={{ width: 140 }} min={0} />
        </Form.Item>
        <Form.Item label="Submitted at" style={{ marginBottom: 12 }}>
          <DatePicker value={submittedAt} onChange={(d) => d && setSubmittedAt(d)} allowClear={false} />
        </Form.Item>
      </Space>
      <Form.Item label="Custom fields (JSON)">
        <Input.TextArea value={custom} onChange={(e) => setCustom(e.target.value)} rows={3} style={{ fontFamily: 'monospace' }} />
      </Form.Item>
      <Button type="primary" htmlType="submit">
        {submitText}
      </Button>
    </Form>
  )
}
