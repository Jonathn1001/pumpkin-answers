import type {
  ConfigDocument,
  ConfigSchemaResponse,
  FieldError,
} from "../api/types";
import { getByPath, setByPath } from "./path";
import { isVisible } from "./conditional";
import { WIDGETS, FallbackWidget } from "./widgets";

interface Props {
  schema: ConfigSchemaResponse;
  config: ConfigDocument;
  onChange: (next: ConfigDocument) => void;
  errors: FieldError[];
}

export function SchemaForm({ schema, config, onChange, errors }: Props) {
  const errorsFor = (key: string) =>
    errors.filter(
      (e) =>
        e.field === key ||
        e.field.startsWith(key + ".") ||
        e.field.startsWith(key + "["),
    );
  return (
    <div className="space-y-8">
      {schema.dimensions.map((dim) => (
        <section key={dim.key} className="rounded-lg border bg-white p-4">
          <h3 className="mb-3 font-semibold">{dim.key}</h3>
          <div className="space-y-4">
            {dim.ui
              .filter((d) => isVisible(d, config))
              .map((d) => {
                const Widget = WIDGETS[d.widget] ?? FallbackWidget;
                return (
                  <Widget
                    key={d.key}
                    descriptor={d}
                    value={getByPath(config, d.key)}
                    config={config}
                    errors={errorsFor(d.key)}
                    onChange={(v: unknown) =>
                      onChange(setByPath(config, d.key, v))
                    }
                  />
                );
              })}
          </div>
        </section>
      ))}
    </div>
  );
}
