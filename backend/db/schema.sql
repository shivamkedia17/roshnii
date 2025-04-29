-- Connect to your PostgreSQL instance as a superuser (e.g., postgres) first
-- then run these commands using psql or a GUI tool:

-- Create the database (if it doesn't exist)
CREATE DATABASE roshnii_db;

-- Create the user (if it doesn't exist)
-- Choose a strong password in practice!
CREATE USER roshnii WITH PASSWORD 'abcd1234';

-- Grant privileges to the user on the database
GRANT ALL PRIVILEGES ON DATABASE roshnii_db TO roshnii;

-- \c roshnii_db roshnii  <-- Connect to the new database as the new user

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    google_id VARCHAR(255) UNIQUE, -- Allow NULL if using dev login only initially
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    picture_url TEXT,
    auth_provider VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create the images table
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY, -- Changed to UUID type based on image_handler.go
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    storage_path TEXT, -- Path in blob storage
    content_type VARCHAR(100),
    size BIGINT,
    width INT,
    height INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    -- Add indexes later, e.g., ON user_id
);

-- Optional: Create a function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply the trigger to users table
CREATE TRIGGER set_users_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Apply the trigger to images table
CREATE TRIGGER set_images_timestamp
BEFORE UPDATE ON images
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Grant necessary permissions on tables to the user (if not covered by GRANT ALL on DB)
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE users TO roshnii;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE images TO roshnii;
GRANT USAGE, SELECT ON SEQUENCE users_id_seq TO roshnii; -- Grant usage on the sequence for id generation
