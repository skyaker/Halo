INSERT INTO auth_credentials (user_id, login, password_hash, created_at) VALUES
(
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'alice',
  -- Пароль: "alice123" (bcrypt hash)
  '$2a$10$N9qo8uLOickgx2ZMRZoMy.MH/rJkH3p1P3Y/7sYQQYHZ4Q8WZUQe',
  '2024-01-10 08:00:00'
),
(
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'bob',
  -- Пароль: "bob123" 
  '$2a$10$VE0tR5c7sYQQYHZ4Q8WZUO5v6z8i1J3p1P3Y/7sYQQYHZ4Q8WZUQe',
  '2024-01-11 09:15:00'
),
(
  'cccccccc-cccc-cccc-cccc-cccccccccccc',
  'charlie',
  -- Пароль: "charlie123"
  '$2a$10$HZ4Q8WZUO5v6z8i1J3p1P3Y/7sYQQYHZ4Q8WZUQeN9qo8uLOickgx',
  '2024-01-12 10:30:00'
);

SELECT * FROM auth_credentials;
