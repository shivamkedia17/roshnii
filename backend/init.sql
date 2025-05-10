-- USE THIS SCRIPT ONLY WHEN MANUALLY CREATING THE DATABASE WITHOUT DOCKER
--
-- Connect to your PostgreSQL instance as a superuser (e.g., postgres) first
-- Create the database (if it doesn't exist)
CREATE DATABASE roshnii_db;

-- Create the user (if it doesn't exist)
-- Choose a strong password in practice!
CREATE USER roshnii
WITH
    PASSWORD 'abcd1234';

-- Grant privileges to the user on the database
GRANT ALL PRIVILEGES ON DATABASE roshnii_db TO roshnii;

-- \c roshnii_db roshnii  <-- Connect to the new database as the new user
