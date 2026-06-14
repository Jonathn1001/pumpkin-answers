import type { ComponentType } from 'react'
import type { WidgetProps } from './types'
import { TextInput, NumberInput, Toggle, SelectInput, ColorInput, LogoInput, FallbackWidget } from './base'
import { ClaimTypeGrid, TierList, CommitteeForm, ChannelMultiSelect, EventsGrid, PerTypeNumberMap, EscalationForm, CustomFieldsEditor } from './complex'

export { FallbackWidget }

export const WIDGETS: Record<string, ComponentType<WidgetProps>> = {
  text: TextInput,
  number: NumberInput,
  toggle: Toggle,
  select: SelectInput,
  color: ColorInput,
  logo: LogoInput,
  'claimtype-grid': ClaimTypeGrid,
  'tier-list': TierList,
  committee: CommitteeForm,
  'channel-multiselect': ChannelMultiSelect,
  'events-grid': EventsGrid,
  'claimtype-number-map': PerTypeNumberMap,
  escalation: EscalationForm,
  customfields: CustomFieldsEditor,
}
