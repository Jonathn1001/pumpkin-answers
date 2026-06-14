import Ajv, { type ValidateFunction } from "ajv";
import { useMemo } from "react";
import type { DimensionSchema, FieldError } from "../api/types";

// Returns a validator: config -> structural field errors (server is the real source of truth).
export function useAjv(dimensions: DimensionSchema[]) {
  return useMemo(() => {
    const ajv = new Ajv({ allErrors: true, strict: false });
    const validators: Record<string, ValidateFunction> = {};
    for (const d of dimensions) {
      try {
        validators[d.key] = ajv.compile(d.jsonSchema);
      } catch {
        /* ignore unreflectable schema */
      }
    }
    return (config: Record<string, unknown>): FieldError[] => {
      const errs: FieldError[] = [];
      for (const [key, validate] of Object.entries(validators)) {
        if (!validate(config[key])) {
          for (const e of validate.errors ?? []) {
            const sub = (e.instancePath || "")
              .replace(/^\//, "")
              .replace(/\//g, ".");
            errs.push({
              field: sub ? `${key}.${sub}` : key,
              message: e.message ?? "invalid",
            });
          }
        }
      }
      return errs;
    };
  }, [dimensions]);
}
