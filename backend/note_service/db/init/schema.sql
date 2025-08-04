CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notes (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  category_id UUID,
  content TEXT NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT,
  ended_at BIGINT,
  completed BOOLEAN DEFAULT FALSE
);

