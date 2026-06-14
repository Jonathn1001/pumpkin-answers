import type { ComponentType } from 'react'
import type { WidgetProps } from './types'
import { TextInput, NumberInput, Toggle, Select, ColorInput, FallbackWidget } from './base'

export { FallbackWidget }

export const WIDGETS: Record<string, ComponentType<WidgetProps>> = {
  text: TextInput,
  number: NumberInput,
  toggle: Toggle,
  select: Select,
  color: ColorInput,
}
