import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { request } from "./client";
import type {
  Tenant,
  ConfigDocument,
  ConfigVersion,
  ClaimDecision,
  Claim,
  Change,
  ConfigSchemaResponse,
} from "./types";

export const keys = {
  tenants: ["tenants"] as const,
  tenant: (slug: string) => ["tenant", slug] as const,
  activeConfig: (slug: string) => ["config", slug] as const,
  versions: (slug: string) => ["versions", slug] as const,
  schema: ["config-schema"] as const,
};

export const useTenants = () =>
  useQuery({
    queryKey: keys.tenants,
    queryFn: () => request<Tenant[]>("GET", "/tenants"),
  });
export const useTenant = (slug: string) =>
  useQuery({
    queryKey: keys.tenant(slug),
    queryFn: () => request<Tenant>("GET", `/tenants/${slug}`),
  });
export const useActiveConfig = (slug: string) =>
  useQuery({
    queryKey: keys.activeConfig(slug),
    queryFn: () => request<ConfigDocument>("GET", `/tenants/${slug}/config`),
  });
export const useVersions = (slug: string) =>
  useQuery({
    queryKey: keys.versions(slug),
    queryFn: () => request<ConfigVersion[]>("GET", `/tenants/${slug}/versions`),
  });
export const useConfigSchema = () =>
  useQuery({
    queryKey: keys.schema,
    queryFn: () => request<ConfigSchemaResponse>("GET", "/config-schema"),
    staleTime: Infinity,
  });

export function useCreateTenant() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (b: { name: string; cloneFrom?: string }) =>
      request<Tenant>("POST", "/tenants", b),
    onSuccess: () => qc.invalidateQueries({ queryKey: keys.tenants }),
  });
}
export function useUpdateTenant(slug: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (b: { name: string; status?: string }) =>
      request<Tenant>("PATCH", `/tenants/${slug}`, b),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: keys.tenants });
      qc.invalidateQueries({ queryKey: keys.tenant(slug) });
    },
  });
}
export function useSaveDraft(slug: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (b: { config: ConfigDocument; note?: string }) =>
      request<ConfigVersion>("POST", `/tenants/${slug}/versions`, b),
    onSuccess: () => qc.invalidateQueries({ queryKey: keys.versions(slug) }),
  });
}
export function usePublish(slug: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (n: number) =>
      request<{ published: number }>(
        "POST",
        `/tenants/${slug}/versions/${n}/publish`,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: keys.activeConfig(slug) });
      qc.invalidateQueries({ queryKey: keys.versions(slug) });
      qc.invalidateQueries({ queryKey: keys.tenant(slug) });
    },
  });
}
export function useRollback(slug: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (targetVersion: number) =>
      request<ConfigVersion>("POST", `/tenants/${slug}/rollback`, {
        targetVersion,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: keys.activeConfig(slug) });
      qc.invalidateQueries({ queryKey: keys.versions(slug) });
      qc.invalidateQueries({ queryKey: keys.tenant(slug) });
    },
  });
}
export function usePreview(slug: string) {
  return useMutation({
    mutationFn: (b: {
      claim: Claim;
      versionNumber?: number;
      config?: ConfigDocument;
    }) => request<ClaimDecision>("POST", `/tenants/${slug}/preview`, b),
  });
}
export function useProcess(slug: string) {
  return useMutation({
    mutationFn: (claim: Claim) =>
      request<ClaimDecision>("POST", `/tenants/${slug}/process`, claim),
  });
}
export function useDiff() {
  return useMutation({
    mutationFn: ({ left, right }: { left: string; right: string }) =>
      request<{ changes: Change[] }>(
        "GET",
        `/diff?left=${encodeURIComponent(left)}&right=${encodeURIComponent(right)}`,
      ),
  });
}
