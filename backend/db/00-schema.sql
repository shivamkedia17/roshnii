-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    google_id VARCHAR(255) UNIQUE, -- Allow NULL if using dev login only initially
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    picture_url TEXT,
    auth_provider VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create the images table
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY, -- Changed to UUID type based on image_handler.go
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    storage_path TEXT, -- Path in blob storage
    content_type VARCHAR(100),
    size BIGINT,
    width INT,
    height INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    -- Add indexes later, e.g., ON user_id
);

-- Albums table
CREATE TABLE IF NOT EXISTS albums (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Album-Image relationship table (many-to-many)
CREATE TABLE IF NOT EXISTS album_images (
    album_id UUID NOT NULL REFERENCES albums (id) ON DELETE CASCADE,
    image_id UUID NOT NULL REFERENCES images (id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    PRIMARY KEY (album_id, image_id)
);
