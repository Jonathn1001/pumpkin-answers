ALTER TABLE tenants DROP CONSTRAINT IF EXISTS fk_active_version;
DROP TABLE IF EXISTS config_versions;
DROP TABLE IF EXISTS tenants;
