INSERT INTO notes (user_id, type_id, content, created_at, updated_at, ended_at, completed) VALUES
(
  'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 
  'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
  'Купить продукты: молоко, хлеб, яйца',
  '2024-01-15 09:00:00',
  '2024-01-15 10:30:00',
  '2024-01-15 10:30:00',
  TRUE
),
(
  'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
  'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a23',
  'Записаться на приём к врачу',
  '2024-01-16 14:00:00',
  '2024-01-17 18:00:00',
  '2024-01-17 18:00:00',
  TRUE
),
(
  'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a12',
  'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a24',
  'Написать отчёт по проекту до пятницы',
  '2024-02-01 10:00:00',
  '2024-02-03 16:45:00',
  NULL,
  FALSE
),
(
  'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a12',
  'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
  'Позвонить клиенту по поводу договора',
  '2024-02-05 11:20:00',
  '2024-02-05 11:20:00',
  NULL,
  FALSE
),
(
  'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a13',
  'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a23',
  'Оплатить счета за коммунальные услуги',
  '2024-01-25 08:00:00',
  '2024-01-30 12:00:00',
  NULL,
  FALSE
),
(
  'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
  'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a24',
  'Обновить резюме на hh.ru',
  '2024-01-10 15:30:00',
  '2024-01-20 09:15:00',
  NULL,
  FALSE
);

SELECT COUNT(*) AS total_notes FROM notes;
