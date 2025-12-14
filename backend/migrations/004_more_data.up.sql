-- Дополнительные проекты
INSERT INTO projects (id, name, stage, description, start_date, end_date, created_by)
VALUES
  ('cccccccc-cccc-cccc-cccc-ccccccccccaa', 'Бизнес-центр «Юг»', 'Монтаж фасада', 'Высотный бизнес-центр класса А', '2025-02-01', '2026-01-15', '11111111-1111-1111-1111-111111111111'),
  ('dddddddd-dddd-dddd-dddd-ddddddddddee', 'Производственный комплекс «Технопарк»', 'Пусконаладка', 'Индустриальный объект с логистическим центром', '2024-08-10', '2025-10-30', '11111111-1111-1111-1111-111111111111')
ON CONFLICT (id) DO NOTHING;

INSERT INTO project_members (project_id, user_id, role) VALUES
  ('cccccccc-cccc-cccc-cccc-ccccccccccaa', '11111111-1111-1111-1111-111111111111', 'manager'),
  ('cccccccc-cccc-cccc-cccc-ccccccccccaa', '22222222-2222-2222-2222-222222222222', 'engineer'),
  ('dddddddd-dddd-dddd-dddd-ddddddddddee', '11111111-1111-1111-1111-111111111111', 'manager'),
  ('dddddddd-dddd-dddd-dddd-ddddddddddee', '22222222-2222-2222-2222-222222222222', 'engineer')
ON CONFLICT DO NOTHING;

-- Дополнительные дефекты
INSERT INTO defects (id, project_id, title, description, priority, severity, status, assignee_id, due_date, created_by, updated_by)
VALUES
  (
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee05',
    'cccccccc-cccc-cccc-cccc-ccccccccccaa',
    'Коррозия металлоконструкций',
    'На кровле корпуса С обнаружены очаги коррозии несущих балок.',
    'HIGH',
    'MAJOR',
    'IN_PROGRESS',
    '22222222-2222-2222-2222-222222222222',
    '2025-04-10',
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee06',
    'dddddddd-dddd-dddd-dddd-ddddddddddee',
    'Нарушение герметичности трубопровода',
    'В блоке инженерных систем зафиксирована утечка теплоносителя.',
    'CRITICAL',
    'CRITICAL',
    'IN_REVIEW',
    '22222222-2222-2222-2222-222222222222',
    '2025-03-25',
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee07',
    'cccccccc-cccc-cccc-cccc-ccccccccccaa',
    'Отказ системы пожарной сигнализации',
    'Контроллер блока А выдаёт ложные тревоги при нагрузке > 60%.',
    'MEDIUM',
    'MAJOR',
    'NEW',
    '22222222-2222-2222-2222-222222222222',
    '2025-04-05',
    '11111111-1111-1111-1111-111111111111',
    NULL
  ),
  (
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee08',
    'dddddddd-dddd-dddd-dddd-ddddddddddee',
    'Деформация откатных ворот',
    'На логистической зоне выявлено залипание ворот при закрытии.',
    'LOW',
    'MINOR',
    'IN_PROGRESS',
    '22222222-2222-2222-2222-222222222222',
    '2025-05-15',
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222'
  )
ON CONFLICT (id) DO NOTHING;

INSERT INTO defect_comments (id, defect_id, author_id, body)
VALUES
  ('eeeeeeee-cccc-cccc-cccc-cccccccccc01', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee05', '22222222-2222-2222-2222-222222222222', 'Выполнена зачистка, требуется повторная проверка через 7 дней'),
  ('eeeeeeee-cccc-cccc-cccc-cccccccccc02', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee06', '11111111-1111-1111-1111-111111111111', 'Планируется поставка нового контроллера 25.03'),
  ('eeeeeeee-cccc-cccc-cccc-cccccccccc03', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee07', '22222222-2222-2222-2222-222222222222', 'Запрошена диагностика у подрядчика ЗАО «Борей»')
ON CONFLICT (id) DO NOTHING;
