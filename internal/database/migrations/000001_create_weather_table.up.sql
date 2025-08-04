CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS weather (
    id UUID PRIMARY KEY,
    city_name TEXT NOT NULL,
    country TEXT NOT NULL,
    temperature DOUBLE PRECISION,
    description TEXT,
    humidity INTEGER,
    wind_speed DOUBLE PRECISION,
    fetched_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
