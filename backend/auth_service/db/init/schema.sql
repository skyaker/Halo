CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE auth_credentials (
  user_id UUID PRIMARY KEY,
  login VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at BIGINT NOT NULL
);
