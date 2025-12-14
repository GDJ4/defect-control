INSERT INTO users (id, email, password_hash, full_name, role)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'manager@systemacontrola.ru', '$2a$12$JH6d6qxGn8pYtBXexdFvcuoS4x2ma/pk3jfbQ3VIAtF5NCz7L3A2K', 'Лазарев Михаил', 'manager'),
  ('22222222-2222-2222-2222-222222222222', 'engineer@systemacontrola.ru', '$2a$12$JH6d6qxGn8pYtBXexdFvcuoS4x2ma/pk3jfbQ3VIAtF5NCz7L3A2K', 'Кузнецов Андрей', 'engineer'),
  ('33333333-3333-3333-3333-333333333333', 'observer@systemacontrola.ru', '$2a$12$JH6d6qxGn8pYtBXexdFvcuoS4x2ma/pk3jfbQ3VIAtF5NCz7L3A2K', 'Петрова Мария', 'observer')
ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, name, stage, description, start_date, end_date, created_by)
VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ЖК «Север»', 'Монтаж инженерных сетей', 'Жилой комплекс из 4 корпусов рядом с МКАД', '2025-01-15', '2025-12-30', '11111111-1111-1111-1111-111111111111'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'ЖК «Восток»', 'Отделочные работы', 'Многофункциональный комплекс с ТЦ', '2024-11-01', '2025-09-15', '11111111-1111-1111-1111-111111111111')
ON CONFLICT (id) DO NOTHING;

INSERT INTO project_members (project_id, user_id, role) VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '22222222-2222-2222-2222-222222222222', 'engineer'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'manager'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 'engineer'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'manager')
ON CONFLICT DO NOTHING;

INSERT INTO defects (id, project_id, title, description, priority, severity, status, assignee_id, due_date, created_by, updated_by)
VALUES
  (
    'dddddddd-dddd-dddd-dddd-dddddddddd01',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'Трещина в несущей стене 3 этажа',
    'Обнаружено растрескивание панели в секции Б. Требуется обследование и усиление.',
    'HIGH',
    'MAJOR',
    'IN_PROGRESS',
    '22222222-2222-2222-2222-222222222222',
    '2025-02-15',
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    'dddddddd-dddd-dddd-dddd-dddddddddd02',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'Отслоение плитки в шахте лифта',
    'В шахте лифта секции Г наблюдается отслоение плитки на площади ~3м².',
    'MEDIUM',
    'MINOR',
    'NEW',
    '22222222-2222-2222-2222-222222222222',
    '2025-03-01',
    '11111111-1111-1111-1111-111111111111',
    NULL
  ),
  (
    'dddddddd-dddd-dddd-dddd-dddddddddd03',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'Отсутствует герметизация витражей',
    'На фасаде корпуса А отсутствует часть герметика, возможна утечка влаги.',
    'CRITICAL',
    'CRITICAL',
    'IN_REVIEW',
    '22222222-2222-2222-2222-222222222222',
    '2025-01-30',
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    'dddddddd-dddd-dddd-dddd-dddddddddd04',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'Некачественная шпаклевка',
    'В помещениях 5-6 этажей появляются наплывы и трещины на шпаклевке.',
    'LOW',
    'MINOR',
    'CLOSED',
    '22222222-2222-2222-2222-222222222222',
    '2024-12-20',
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222'
  )
ON CONFLICT (id) DO NOTHING;

INSERT INTO defect_comments (id, defect_id, author_id, body)
VALUES
  ('cccccccc-cccc-cccc-cccc-cccccccccc01', 'dddddddd-dddd-dddd-dddd-dddddddddd01', '22222222-2222-2222-2222-222222222222', 'Проведено обследование, требуется усиление каркаса.'),
  ('cccccccc-cccc-cccc-cccc-cccccccccc02', 'dddddddd-dddd-dddd-dddd-dddddddddd02', '11111111-1111-1111-1111-111111111111', 'Контрольный осмотр запланирован на 15.02.')
ON CONFLICT (id) DO NOTHING;
