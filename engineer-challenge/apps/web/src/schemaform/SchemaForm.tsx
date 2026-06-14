import { Card, Space } from 'antd'
import type { ConfigDocument, ConfigSchemaResponse, FieldError } from '../api/types'
import { getByPath, setByPath } from './path'
import { isVisible, isRequired } from './conditional'
import { WIDGETS, FallbackWidget } from './widgets'

interface Props {
  schema: ConfigSchemaResponse
  config: ConfigDocument
  onChange: (next: ConfigDocument) => void
  errors: FieldError[]
}

export function SchemaForm({ schema, config, onChange, errors }: Props) {
  const errorsFor = (key: string) =>
    errors.filter(
      (e) => e.field === key || e.field.startsWith(key + '.') || e.field.startsWith(key + '['),
    )
  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      {schema.dimensions.map((dim) => (
        <Card key={dim.key} size="small" title={<span style={{ textTransform: 'capitalize' }}>{dim.key}</span>}>
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            {dim.ui
              .filter((d) => isVisible(d, config))
              .map((d) => {
                const Widget = WIDGETS[d.widget] ?? FallbackWidget
                const descriptor = { ...d, required: isRequired(d, config) }
                return (
                  <Widget
                    key={d.key}
                    descriptor={descriptor}
                    value={getByPath(config, d.key)}
                    config={config}
                    errors={errorsFor(d.key)}
                    onChange={(v: unknown) => onChange(setByPath(config, d.key, v))}
                  />
                )
              })}
          </Space>
        </Card>
      ))}
    </Space>
  )
}
