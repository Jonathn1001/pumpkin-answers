import { useState } from "react";
import type { Claim, ClaimType } from "../api/types";

const TYPES: ClaimType[] = [
  "OUTPATIENT",
  "INPATIENT",
  "DENTAL",
  "MATERNITY",
  "OPTICAL",
];
const DEFAULT_CUSTOM =
  '{"employeeId":"EMP1234","policyNumber":"HF-12345678","memberTier":"Gold","nationalId":"123456789012","citizenCategory":"General"}';

export function ClaimForm({ onSubmit }: { onSubmit: (c: Claim) => void }) {
  const [type, setType] = useState<ClaimType>("OUTPATIENT");
  const [amount, setAmount] = useState(10000);
  const [submittedAt, setSubmittedAt] = useState("2026-06-14");
  const [custom, setCustom] = useState(DEFAULT_CUSTOM);
  return (
    <form
      className="space-y-2"
      onSubmit={(e) => {
        e.preventDefault();
        let cf: Record<string, unknown> = {};
        try {
          cf = JSON.parse(custom);
        } catch {
          /* ignore */
        }
        onSubmit({
          type,
          amount,
          submittedAt: new Date(submittedAt).toISOString(),
          customFields: cf,
        });
      }}
    >
      <div className="flex flex-wrap gap-2">
        <select
          className="rounded border px-2 py-1"
          value={type}
          onChange={(e) => setType(e.target.value as ClaimType)}
        >
          {TYPES.map((t) => (
            <option key={t}>{t}</option>
          ))}
        </select>
        <input
          type="number"
          className="rounded border px-2 py-1"
          value={amount}
          onChange={(e) => setAmount(Number(e.target.value))}
        />
        <input
          type="date"
          className="rounded border px-2 py-1"
          value={submittedAt}
          onChange={(e) => setSubmittedAt(e.target.value)}
        />
      </div>
      <textarea
        className="w-full rounded border px-2 py-1 font-mono text-xs"
        rows={3}
        value={custom}
        onChange={(e) => setCustom(e.target.value)}
      />
      <button
        className="rounded bg-blue-600 px-3 py-1 text-white"
        type="submit"
      >
        Run preview
      </button>
    </form>
  );
}
