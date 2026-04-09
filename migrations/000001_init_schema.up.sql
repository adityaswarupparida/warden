-- 000001_init_schema.up.sql

CREATE TYPE service_status AS ENUM ('healthy', 'down', 'degraded');                                                                                                          
CREATE TYPE incident_severity AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE incident_status AS ENUM ('identified', 'investigating', 'monitoring', 'resolved'); 

CREATE TABLE IF NOT EXISTS services (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    status      service_status NOT NULL DEFAULT 'healthy',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS incidents (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id  UUID NOT NULL REFERENCES services(id),
    title       TEXT NOT NULL,
    description TEXT NOT NULL,
    status      incident_status NOT NULL DEFAULT 'investigating',
    severity    incident_severity NOT NULL DEFAULT 'low',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS incident_updates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id),
    message     TEXT NOT NULL,
    status      incident_status NOT NULL DEFAULT 'investigating',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);