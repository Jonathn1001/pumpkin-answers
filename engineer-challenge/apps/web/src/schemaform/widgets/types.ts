import type { ConfigDocument, FieldDescriptor, FieldError } from '../../api/types'

export interface WidgetProps {
  descriptor: FieldDescriptor
  value: unknown
  onChange: (v: unknown) => void
  config: ConfigDocument
  errors: FieldError[]
}
