-- CREATE DATABASE "TeamUp";

-- \c "TeamUp"

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE TABLE IF NOT EXISTS users (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     username VARCHAR(50) UNIQUE NOT NULL,
--     name VARCHAR(100) NOT NULL,
--     lastname VARCHAR(100) NOT NULL,
--     phone_number VARCHAR(20),
--     email VARCHAR(255) UNIQUE NOT NULL,
--     password TEXT,
--     oauth_provider VARCHAR(50),
--     oauth_provider_id VARCHAR(255)
-- );