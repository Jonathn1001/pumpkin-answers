import { useState } from 'react'
import { Card, Space } from 'antd'
import { usePreview } from '../../api/hooks'
import { ClaimForm } from '../../components/ClaimForm'
import { DecisionView } from '../../components/DecisionView'
import type { ClaimDecision } from '../../api/types'

export function PreviewTab({ slug }: { slug: string }) {
  const preview = usePreview(slug)
  const [decision, setDecision] = useState<ClaimDecision | null>(null)
  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Card size="small" title="Sample claim">
        <ClaimForm onSubmit={(claim) => preview.mutate({ claim }, { onSuccess: setDecision })} />
      </Card>
      {decision && (
        <Card size="small" title="Predicted outcome">
          <DecisionView d={decision} />
        </Card>
      )}
    </Space>
  )
}
