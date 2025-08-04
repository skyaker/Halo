CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE categories (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT
);

CREATE UNIQUE INDEX unique_user_category_name
  ON categories(user_id, name);
