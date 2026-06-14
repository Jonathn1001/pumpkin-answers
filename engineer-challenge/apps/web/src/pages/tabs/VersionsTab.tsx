import { useState } from "react";
import { useVersions, useRollback, useTenant, useDiff } from "../../api/hooks";
import { DiffView } from "../../components/DiffView";
import type { Change } from "../../api/types";

export function VersionsTab({ slug }: { slug: string }) {
  const versions = useVersions(slug);
  const tenant = useTenant(slug);
  const rollback = useRollback(slug);
  const diff = useDiff();
  const [selected, setSelected] = useState<number | null>(null);
  const [changes, setChanges] = useState<Change[] | null>(null);

  if (versions.isLoading) return <div>Loading…</div>;
  const active = tenant.data?.activeVersionNumber;

  function select(n: number) {
    setSelected(n);
    setChanges(null);
    if (n > 1) {
      diff.mutate(
        { left: `${slug}@${n - 1}`, right: `${slug}@${n}` },
        { onSuccess: (r) => setChanges(r.changes) },
      );
    }
  }

  return (
    <div className="grid gap-4 md:grid-cols-2">
      <table className="h-fit w-full text-sm">
        <thead>
          <tr className="border-b text-left">
            <th className="p-2">Version</th>
            <th>Status</th>
            <th>By</th>
            <th>When</th>
            <th>Note</th>
          </tr>
        </thead>
        <tbody>
          {versions.data?.map((v) => (
            <tr
              key={v.versionNumber}
              className={`cursor-pointer border-b hover:bg-gray-50 ${v.versionNumber === selected ? "bg-blue-50" : ""}`}
              onClick={() => select(v.versionNumber)}
            >
              <td className="p-2">
                v{v.versionNumber}
                {v.versionNumber === active ? " (active)" : ""}
              </td>
              <td>{v.status}</td>
              <td>{v.createdBy || "—"}</td>
              <td>{v.createdAt.slice(0, 10)}</td>
              <td>{v.note}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className="rounded-lg border bg-white p-4">
        {selected == null ? (
          <p className="text-sm text-gray-500">
            Select a version to see what changed.
          </p>
        ) : (
          <div className="space-y-3">
            <div className="flex items-center justify-between gap-2">
              <h3 className="font-semibold">
                Version details — v{selected}
                {selected === active ? " (active)" : ""}
              </h3>
              {selected !== active && (
                <button
                  className="rounded bg-blue-600 px-3 py-1 text-sm text-white disabled:opacity-50"
                  disabled={rollback.isPending}
                  onClick={() => rollback.mutate(selected)}
                >
                  Rollback to this version
                </button>
              )}
            </div>
            {selected === 1 ? (
              <p className="text-sm text-gray-500">
                Initial version — nothing to compare against.
              </p>
            ) : diff.isPending && !changes ? (
              <p className="text-sm text-gray-500">Loading diff…</p>
            ) : changes ? (
              <div className="space-y-1">
                <p className="text-xs text-gray-500">
                  Changes from v{selected - 1} → v{selected}:
                </p>
                <DiffView changes={changes} />
              </div>
            ) : null}
          </div>
        )}
      </div>
    </div>
  );
}
