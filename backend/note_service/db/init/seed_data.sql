INSERT INTO notes (id, user_id, category_id, content, created_at, updated_at, ended_at, completed) VALUES
(
  '11111111-1111-1111-1111-111111111111',
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'dddddddd-dddd-dddd-dddd-dddddddddddd',
  'Купить продукты: молоко, хлеб, яйца',
  '2024-01-15 09:00:00',
  '2024-01-15 10:30:00',
  '2024-01-15 10:30:00',
  TRUE
),
(
  '22222222-2222-2222-2222-222222222222',
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
  'Записаться на приём к врачу',
  '2024-01-16 14:00:00',
  '2024-01-17 18:00:00',
  '2024-01-17 18:00:00',
  TRUE
),
(
  '33333333-3333-3333-3333-333333333333',
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'ffffffff-ffff-ffff-ffff-ffffffffffff',
  'Написать отчёт по проекту до пятницы',
  '2024-02-01 10:00:00',
  '2024-02-03 16:45:00',
  NULL,
  FALSE
),
(
  '44444444-4444-4444-4444-444444444444',
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'dddddddd-dddd-dddd-dddd-dddddddddddd',
  'Позвонить клиенту по поводу договора',
  '2024-02-05 11:20:00',
  '2024-02-05 11:20:00',
  NULL,
  FALSE
),
(
  '55555555-5555-5555-5555-555555555555',
  'cccccccc-cccc-cccc-cccc-cccccccccccc',
  'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
  'Оплатить счета за коммунальные услуги',
  '2024-01-25 08:00:00',
  '2024-01-30 12:00:00',
  NULL,
  FALSE
),
(
  '66666666-6666-6666-6666-666666666666',
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'ffffffff-ffff-ffff-ffff-ffffffffffff',
  'Обновить резюме на hh.ru',
  '2024-01-10 15:30:00',
  '2024-01-20 09:15:00',
  NULL,
  FALSE
);

SELECT 
  id, 
  LEFT(content, 20) || '...' AS preview, 
  completed,
  CASE WHEN ended_at IS NULL THEN 'active' ELSE 'completed' END AS status
FROM notes;
