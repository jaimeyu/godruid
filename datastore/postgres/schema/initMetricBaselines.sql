
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

COMMIT;