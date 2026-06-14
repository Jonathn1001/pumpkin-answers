CREATE TABLE tenants (
  id                       BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  slug                     TEXT NOT NULL UNIQUE,
  name                     TEXT NOT NULL,
  status                   TEXT NOT NULL DEFAULT 'active',
  active_config_version_id BIGINT NULL,
  created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE config_versions (
  id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  tenant_id      BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  version_number INT  NOT NULL,
  status         TEXT NOT NULL,
  note           TEXT,
  config         JSONB NOT NULL,
  created_by     TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, version_number)
);

ALTER TABLE tenants ADD CONSTRAINT fk_active_version
  FOREIGN KEY (active_config_version_id) REFERENCES config_versions(id) ON DELETE SET NULL;
