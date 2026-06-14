import { useState } from 'react'
import { Button, ColorPicker, Input, InputNumber, Select, Space, Switch, Typography } from 'antd'
import type { WidgetProps } from './types'

export function FieldLabel({ p }: { p: WidgetProps }) {
  return (
    <div style={{ marginBottom: 4 }}>
      <Typography.Text strong>{p.descriptor.label}</Typography.Text>
      {p.descriptor.required && <Typography.Text type="danger"> *</Typography.Text>}
      {p.errors.map((e, i) => (
        <Typography.Text key={i} type="danger" style={{ marginLeft: 8, fontSize: 12 }}>
          {e.message}
        </Typography.Text>
      ))}
    </div>
  )
}

export function TextInput(p: WidgetProps) {
  return (
    <div>
      <FieldLabel p={p} />
      <Input value={(p.value as string) ?? ''} onChange={(e) => p.onChange(e.target.value)} />
    </div>
  )
}

export function NumberInput(p: WidgetProps) {
  return (
    <div>
      <FieldLabel p={p} />
      <InputNumber style={{ width: '100%' }} value={(p.value as number) ?? 0} onChange={(v) => p.onChange(v ?? 0)} />
    </div>
  )
}

export function Toggle(p: WidgetProps) {
  return (
    <Space>
      <Switch checked={Boolean(p.value)} onChange={(v) => p.onChange(v)} />
      <Typography.Text>{p.descriptor.label}</Typography.Text>
    </Space>
  )
}

export function SelectInput(p: WidgetProps) {
  return (
    <div>
      <FieldLabel p={p} />
      <Select
        style={{ width: '100%' }}
        value={(p.value as string) || undefined}
        onChange={(v) => p.onChange(v)}
        options={(p.descriptor.options ?? []).map((o) => ({ value: o, label: o }))}
      />
    </div>
  )
}

export function ColorInput(p: WidgetProps) {
  return (
    <div>
      <FieldLabel p={p} />
      <ColorPicker value={(p.value as string) || '#000000'} onChange={(_, hex) => p.onChange(hex)} showText />
    </div>
  )
}

export function LogoInput(p: WidgetProps) {
  const url = (p.value as string) ?? ''
  const [brokenUrl, setBrokenUrl] = useState('')
  const showImg = url !== '' && url !== brokenUrl
  return (
    <div>
      <FieldLabel p={p} />
      <Space>
        {showImg ? (
          <img
            src={url}
            alt="logo preview"
            style={{ height: 40, width: 40, objectFit: 'contain', border: '1px solid #eee', borderRadius: 4 }}
            onError={() => setBrokenUrl(url)}
          />
        ) : (
          <div
            style={{
              height: 40,
              width: 40,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              border: '1px solid #eee',
              borderRadius: 4,
              fontSize: 10,
              color: '#aaa',
            }}
          >
            {url ? 'invalid' : 'no logo'}
          </div>
        )}
        <Input style={{ width: 280 }} placeholder="https://…" value={url} onChange={(e) => p.onChange(e.target.value)} />
        {url && (
          <Button type="link" onClick={() => p.onChange('')}>
            Remove
          </Button>
        )}
      </Space>
    </div>
  )
}

export function FallbackWidget(p: WidgetProps) {
  return (
    <div>
      <FieldLabel p={p} />
      <Input.TextArea
        rows={4}
        defaultValue={JSON.stringify(p.value, null, 2)}
        onBlur={(e) => {
          try {
            p.onChange(JSON.parse(e.target.value))
          } catch {
            /* ignore invalid json */
          }
        }}
        style={{ fontFamily: 'monospace', fontSize: 12 }}
      />
    </div>
  )
}
