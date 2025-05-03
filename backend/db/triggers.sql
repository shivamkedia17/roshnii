
-- Create a function to automatically update updated_at timestamp
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


-- Apply the timestamp trigger to albums table
CREATE TRIGGER set_albums_timestamp
BEFORE UPDATE ON albums
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
