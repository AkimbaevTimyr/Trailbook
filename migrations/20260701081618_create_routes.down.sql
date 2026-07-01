-- 4.2.2 Откат маршрутов
DROP INDEX IF EXISTS idx_route_points_order;
DROP INDEX IF EXISTS idx_route_points_route;
DROP INDEX IF EXISTS idx_routes_created_by;
DROP INDEX IF EXISTS idx_routes_distance;
DROP INDEX IF EXISTS idx_routes_difficulty;
DROP INDEX IF EXISTS idx_routes_type;
DROP INDEX IF EXISTS idx_routes_region;

DROP TABLE IF EXISTS route_points;
DROP TABLE IF EXISTS routes;

DROP TYPE IF EXISTS difficulty;
DROP TYPE IF EXISTS route_type;
