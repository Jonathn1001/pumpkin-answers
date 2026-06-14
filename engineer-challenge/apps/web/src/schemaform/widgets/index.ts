import type { ComponentType } from 'react'
import type { WidgetProps } from './types'
import { TextInput, NumberInput, Toggle, Select, ColorInput, FallbackWidget } from './base'
import { ClaimTypeGrid, TierList, CommitteeForm, ChannelMultiSelect, EventsGrid, PerTypeNumberMap, EscalationForm, CustomFieldsEditor } from './complex'

export { FallbackWidget }

export const WIDGETS: Record<string, ComponentType<WidgetProps>> = {
  text: TextInput,
  number: NumberInput,
  toggle: Toggle,
  select: Select,
  color: ColorInput,
  'claimtype-grid': ClaimTypeGrid,
  'tier-list': TierList,
  committee: CommitteeForm,
  'channel-multiselect': ChannelMultiSelect,
  'events-grid': EventsGrid,
  'claimtype-number-map': PerTypeNumberMap,
  escalation: EscalationForm,
  customfields: CustomFieldsEditor,
}
