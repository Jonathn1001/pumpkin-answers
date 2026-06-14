import { useVersions, useRollback, useTenant } from "../../api/hooks";

export function VersionsTab({ slug }: { slug: string }) {
  const versions = useVersions(slug);
  const tenant = useTenant(slug);
  const rollback = useRollback(slug);
  if (versions.isLoading) return <div>Loading…</div>;
  const active = tenant.data?.activeVersionNumber;
  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b text-left">
          <th className="p-2">Version</th>
          <th>Status</th>
          <th>By</th>
          <th>When</th>
          <th>Note</th>
          <th />
        </tr>
      </thead>
      <tbody>
        {versions.data?.map((v) => (
          <tr key={v.versionNumber} className="border-b">
            <td className="p-2">
              v{v.versionNumber}
              {v.versionNumber === active ? " (active)" : ""}
            </td>
            <td>{v.status}</td>
            <td>{v.createdBy || "—"}</td>
            <td>{v.createdAt.slice(0, 10)}</td>
            <td>{v.note}</td>
            <td>
              {v.versionNumber !== active && (
                <button
                  className="text-blue-700"
                  onClick={() => rollback.mutate(v.versionNumber)}
                >
                  Rollback
                </button>
              )}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
