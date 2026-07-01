-- 4.2.3 Походы
CREATE TABLE IF NOT EXISTS trips (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    route_id        INTEGER REFERENCES routes(id) ON DELETE SET NULL,

    started_at      TIMESTAMP WITH TIME ZONE,
    finished_at     TIMESTAMP WITH TIME ZONE,

    distance_km     DECIMAL(8,2),
    elevation_gain  INTEGER,
    duration_hours  DECIMAL(4,1),

    gpx_file_url    VARCHAR(500),
    notes           TEXT,

    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_trips_user ON trips(user_id);
CREATE INDEX idx_trips_route ON trips(route_id);

CREATE TABLE IF NOT EXISTS trip_photos (
    id          SERIAL PRIMARY KEY,
    trip_id     INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    photo_url   VARCHAR(500) NOT NULL,
    lat         DECIMAL(10, 8),
    lon         DECIMAL(11, 8),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_trip_photos_trip ON trip_photos(trip_id);
