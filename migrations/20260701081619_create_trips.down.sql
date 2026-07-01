-- 4.2.3 Откат походов
DROP INDEX IF EXISTS idx_trip_photos_trip;
DROP INDEX IF EXISTS idx_trips_route;
DROP INDEX IF EXISTS idx_trips_user;

DROP TABLE IF EXISTS trip_photos;
DROP TABLE IF EXISTS trips;
