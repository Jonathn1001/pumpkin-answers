import type { BrandingConfig } from "../api/types";

export function applyBranding(b?: BrandingConfig) {
  const root = document.documentElement;
  root.style.setProperty("--brand-primary", b?.primaryColor || "#1f2937");
  root.style.setProperty("--brand-secondary", b?.secondaryColor || "#374151");
}
