INSERT INTO users (id, username, email, password_hash, created_at) VALUES
(
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'alice',
  'alice@example.com',
),
(
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'bob',
  'bob@example.com',
),
(
  'cccccccc-cccc-cccc-cccc-cccccccccccc',
  'charlie',
  'charlie@example.com',
);

SELECT * FROM users;
