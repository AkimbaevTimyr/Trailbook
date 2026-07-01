-- 4.2.2 Маршруты
CREATE TYPE route_type AS ENUM ('hike', 'alpine', 'ski', 'trail_run');
CREATE TYPE difficulty AS ENUM ('easy', 'moderate', 'hard', 'extreme');

CREATE TABLE IF NOT EXISTS routes (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    region_id       INTEGER REFERENCES regions(id) ON DELETE SET NULL,
    created_by      INTEGER REFERENCES users(id) ON DELETE SET NULL,

    route_type      route_type NOT NULL DEFAULT 'hike',
    difficulty      difficulty NOT NULL DEFAULT 'moderate',
    distance_km     DECIMAL(8,2),
    elevation_gain  INTEGER,                -- метры
    duration_hours  DECIMAL(4,1),           -- часы
    season          VARCHAR(50),            -- 'spring-summer', 'autumn', 'winter', 'year-round'

    -- Геоданные: GeoJSON LineString для универсальности (PostGIS опционально)
    track_geojson   JSONB,

    gpx_file_url    VARCHAR(500),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_routes_region ON routes(region_id);
CREATE INDEX idx_routes_type ON routes(route_type);
CREATE INDEX idx_routes_difficulty ON routes(difficulty);
CREATE INDEX idx_routes_distance ON routes(distance_km);
CREATE INDEX idx_routes_created_by ON routes(created_by);
