import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useActiveConfig, useTenant } from "../api/hooks";
import { applyBranding } from "../theme/applyBranding";
import { ConfigTab } from "./tabs/ConfigTab";
import { PreviewTab } from "./tabs/PreviewTab";
import { VersionsTab } from "./tabs/VersionsTab";
import { CompareTab } from "./tabs/CompareTab";

const TABS = ["Config", "Preview", "Versions", "Compare"] as const;

export function TenantDetail() {
  const { slug = "" } = useParams();
  const tenant = useTenant(slug);
  const config = useActiveConfig(slug);
  const [tab, setTab] = useState<(typeof TABS)[number]>("Config");
  useEffect(() => {
    applyBranding(config.data?.branding);
    return () => applyBranding(undefined);
  }, [config.data]);
  if (tenant.isError) return <div>Tenant not found.</div>;
  return (
    <div className="space-y-4">
      <h2
        className="text-xl font-semibold"
        style={{ color: "var(--brand-primary)" }}
      >
        {tenant.data?.name ?? slug}
      </h2>
      <nav className="flex gap-2 border-b">
        {TABS.map((t) => (
          <button
            key={t}
            className={`px-3 py-1 ${tab === t ? "border-b-2 border-blue-600 font-medium" : ""}`}
            onClick={() => setTab(t)}
          >
            {t}
          </button>
        ))}
      </nav>
      {tab === "Config" && <ConfigTab key={slug} slug={slug} />}
      {tab === "Preview" && <PreviewTab key={slug} slug={slug} />}
      {tab === "Versions" && <VersionsTab key={slug} slug={slug} />}
      {tab === "Compare" && <CompareTab key={slug} slug={slug} />}
    </div>
  );
}
