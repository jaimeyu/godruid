-- DROP INDEX IF EXISTS idx_monitored_object_id;
-- DROP INDEX IF EXISTS idx_hour_of_week;
-- DROP INDEX IF EXISTS idx_tenant_id;
-- DROP INDEX IF EXISTS idx_monitored_object_id_hour_of_week;
create TABLE IF NOT EXISTS metric_baselines (
    tenant_id varchar(256) NOT NULL,
    monitored_object_id varchar(256) NOT NULL,
    hour_of_week int NOT NULL,
    baselines jsonb,
    created_timestamp bigint NOT NULL default 0,
    last_modified_timestamp bigint NOT NULL,
    last_reset_timestamp bigint NOT NULL default 0,
    PRIMARY KEY (tenant_id, monitored_object_id, hour_of_week)
);
CREATE INDEX IF NOT EXISTS idx_monitored_object_id ON metric_baselines(monitored_object_id);
CREATE INDEX IF NOT EXISTS dx_hour_of_week ON metric_baselines(hour_of_week);
CREATE INDEX IF NOT EXISTS idx_tenant_id ON metric_baselines(tenant_id);
CREATE INDEX IF NOT EXISTS idx_monitored_object_id_hour_of_week ON metric_baselines(monitored_object_id, hour_of_week);

COMMIT;