import type { ClaimDecision } from "../api/types";
import { Badge } from "./Badge";

export function DecisionView({ d }: { d: ClaimDecision }) {
  if (!d.accepted) {
    return (
      <div className="rounded border border-red-300 bg-red-50 p-3 text-sm">
        <b>Rejected:</b> {d.rejectionReasons.join("; ")}
      </div>
    );
  }
  const r = d.approval?.route;
  return (
    <div className="space-y-2 rounded border bg-white p-3 text-sm">
      <div>
        <b>Approval:</b>{" "}
        {d.approval?.outcome === "auto_approved" ? (
          <Badge tone="green">auto-approved</Badge>
        ) : (
          <Badge tone="blue">
            routed → {r?.committeeName ?? r?.tierLabel}
            {r?.approverRole ? ` (${r.approverRole})` : ""}
            {r?.requiredApprovals ? ` · ${r.requiredApprovals} approvals` : ""}
          </Badge>
        )}
      </div>
      <div>
        <b>Required documents:</b>{" "}
        {(d.requiredDocuments ?? []).join(", ") || "—"}
      </div>
      <div>
        <b>SLA:</b> {d.slaDays} days → {d.slaDeadline?.slice(0, 10)}
        {d.escalation
          ? ` · warn ${d.escalation.warnBeforeDays}d → ${d.escalation.notifyRole}`
          : ""}
      </div>
      <div>
        <b>Notifications:</b>{" "}
        {d.notifications
          .map((n) => `${n.event} [${n.channels.join(",")}]`)
          .join(" · ")}
      </div>
      {d.customFieldValidation && (
        <div>
          <b>Custom fields:</b>{" "}
          {d.customFieldValidation.valid ? (
            <Badge tone="green">valid</Badge>
          ) : (
            <Badge tone="red">
              {d.customFieldValidation.errors.map((e) => e.field).join(", ")}
            </Badge>
          )}
        </div>
      )}
      <details>
        <summary className="cursor-pointer font-medium">Trace</summary>
        <ul className="ml-4 list-disc">
          {d.trace.map((t, i) => (
            <li key={i}>
              <b>{t.dimension}:</b> {t.explanation}
            </li>
          ))}
        </ul>
      </details>
    </div>
  );
}
