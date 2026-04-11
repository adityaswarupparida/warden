-- 000002_add_occurred_at_to_incidents.up.sql

ALTER TABLE incidents ADD COLUMN occurred_at TIMESTAMPTZ;