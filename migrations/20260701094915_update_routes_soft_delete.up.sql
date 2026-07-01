-- Drop route_points table - will be implemented later
DROP TABLE IF EXISTS route_points;

-- Add soft delete flag to routes
ALTER TABLE routes
    ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_routes_is_deleted ON routes(is_deleted);
