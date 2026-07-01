-- Revert soft delete flag
DROP INDEX IF EXISTS idx_routes_is_deleted;
ALTER TABLE routes DROP COLUMN IF EXISTS is_deleted;

-- Recreate route_points table (simplified, full version will be added in a later migration)
CREATE TABLE IF NOT EXISTS route_points (
    id          SERIAL PRIMARY KEY,
    route_id    INTEGER NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    lat         DECIMAL(10, 8) NOT NULL,
    lon         DECIMAL(11, 8) NOT NULL,
    elevation   DECIMAL(8, 2),
    point_order INTEGER NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_route_points_route ON route_points(route_id);
CREATE INDEX IF NOT EXISTS idx_route_points_order ON route_points(route_id, point_order);
