import { Alert, Collapse, Descriptions, Tag } from 'antd'
import type { ClaimDecision } from '../api/types'
import { dimensionLabel, eventLabel } from '../labels'

export function DecisionView({ d }: { d: ClaimDecision }) {
  if (!d.accepted) {
    return <Alert type="error" showIcon message="Rejected" description={d.rejectionReasons.join('; ')} />
  }
  const r = d.approval?.route
  return (
    <Descriptions column={1} bordered size="small">
      <Descriptions.Item label="Approval">
        {d.approval?.outcome === 'auto_approved' ? (
          <Tag color="green">auto-approved</Tag>
        ) : (
          <Tag color="blue">
            routed → {r?.committeeName ?? r?.tierLabel}
            {r?.approverRole ? ` (${r.approverRole})` : ''}
            {r?.requiredApprovals ? ` · ${r.requiredApprovals} approvals` : ''}
          </Tag>
        )}
      </Descriptions.Item>
      <Descriptions.Item label="Required documents">{(d.requiredDocuments ?? []).join(', ') || '—'}</Descriptions.Item>
      <Descriptions.Item label="SLA">
        {d.slaDays} days → {d.slaDeadline?.slice(0, 10)}
        {d.escalation ? ` · warn ${d.escalation.warnBeforeDays}d → ${d.escalation.notifyRole}` : ''}
      </Descriptions.Item>
      <Descriptions.Item label="Notifications">
        {d.notifications.map((n) => `${eventLabel(n.event)} [${n.channels.join(', ')}]`).join(' · ') || '—'}
      </Descriptions.Item>
      {d.customFieldValidation && (
        <Descriptions.Item label="Custom fields">
          {d.customFieldValidation.valid ? (
            <Tag color="green">valid</Tag>
          ) : (
            <Tag color="red">{d.customFieldValidation.errors.map((e) => e.field).join(', ')}</Tag>
          )}
        </Descriptions.Item>
      )}
      <Descriptions.Item label="Trace">
        <Collapse
          ghost
          size="small"
          items={[
            {
              key: 'trace',
              label: `${d.trace.length} steps`,
              children: (
                <ul style={{ margin: 0, paddingLeft: 18 }}>
                  {d.trace.map((t, i) => (
                    <li key={i}>
                      <b>{dimensionLabel(t.dimension)}:</b> {t.explanation}
                    </li>
                  ))}
                </ul>
              ),
            },
          ]}
        />
      </Descriptions.Item>
    </Descriptions>
  )
}
