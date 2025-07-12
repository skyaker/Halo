INSERT INTO auth_credentials (user_id, login, password_hash, created_at) VALUES
(
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'alice',
  -- Пароль: "alice123" (bcrypt hash)
  '$2a$10$E4qLUYNHrIw23un1iuqoWeVFworsLVp/6jepd3PXwdujtWITgaK9e',
  '1752326323'
),
(
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'bob',
  -- Пароль: "bob123" 
  '$2a$10$SatWwxYMxKFc1wV5H3JO.ekb..ZBJo95rKhXNGkKVgHtu1mk6.aUS',
  '1752326323'
),
(
  'cccccccc-cccc-cccc-cccc-cccccccccccc',
  'charlie',
  -- Пароль: "charlie123"
  '$2a$10$DcUGsLX.uMnd3DqHyE4kQeUAZfg1YKo9Wnk23MOAZDbsR7JgkATF6',
  '1752326323'
);

SELECT * FROM auth_credentials;
