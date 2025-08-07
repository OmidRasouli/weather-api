-- Weather API Database Initialization Script
-- This script sets up the initial database structure

-- Create the weather table
CREATE TABLE IF NOT EXISTS weather (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    temperature DECIMAL(5,2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    humidity INTEGER NOT NULL CHECK (humidity >= 0 AND humidity <= 100),
    wind_speed DECIMAL(5,2) NOT NULL CHECK (wind_speed >= 0),
    fetched_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_weather_city_country ON weather(city, country);
CREATE INDEX IF NOT EXISTS idx_weather_fetched_at ON weather(fetched_at);
CREATE INDEX IF NOT EXISTS idx_weather_created_at ON weather(created_at);

-- Create a function to automatically update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_weather_updated_at 
    BEFORE UPDATE ON weather 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
