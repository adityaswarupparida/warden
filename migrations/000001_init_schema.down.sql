 -- 000001_init_schema.down.sql

DROP TABLE IF EXISTS incident_updates;
DROP TABLE IF EXISTS incidents;
DROP TABLE IF EXISTS services;

DROP TYPE IF EXISTS incident_severity;
DROP TYPE IF EXISTS incident_status;
DROP TYPE IF EXISTS service_status;