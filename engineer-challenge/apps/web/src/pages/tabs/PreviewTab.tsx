import { useState } from 'react'
import { usePreview } from '../../api/hooks'
import { ClaimForm } from '../../components/ClaimForm'
import { DecisionView } from '../../components/DecisionView'
import type { ClaimDecision } from '../../api/types'

export function PreviewTab({ slug }: { slug: string }) {
  const preview = usePreview(slug)
  const [decision, setDecision] = useState<ClaimDecision | null>(null)
  return (
    <div className="space-y-4">
      <ClaimForm onSubmit={claim => preview.mutate({ claim }, { onSuccess: setDecision })} />
      {decision && <DecisionView d={decision} />}
    </div>
  )
}
