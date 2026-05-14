-- +goose Up
-- Seed data to exercise all API features:
--
-- Дерево подразделений (4 уровня):
--
--   Компания (1) [root]
--   ├── Разработка (2)
--   │   ├── Backend (4)
--   │   │   └── Go-команда (7)
--   │   └── Frontend (5)
--   ├── HR (3)
--   │   └── Рекрутинг (6)
--   └── Финансы (8)
--
-- Покрывает следующие сценарии:
--   GET  depth=1 → только непосредственные подразделения (Разработка, HR, Финансы)
--   GET  depth=2 → 2 уровня
--   GET  depth=3 → дерево целиком, включая Go-команду
--   PATCH parent_id → переместить Go-команду в Frontend (корректное обновление иерархии)
--   PATCH parent_id → переместить Разработка в Backend → 409 определение циклической зависимости
--   DELETE cascade → удалание поддеревьев
--   DELETE reassign → сотрудники из удаляемого подразделения переходят в другое подразделение
--   include_employees=false → без сотрудников в ответе
--   Employees with and without hired_at → проверка обработки nullable полей

INSERT INTO departments (id, name, parent_id) VALUES
  (1, 'Компания',    NULL),
  (2, 'Разработка',  1),
  (3, 'HR',          1),
  (4, 'Backend',     2),
  (5, 'Frontend',    2),
  (6, 'Рекрутинг',   3),
  (7, 'Go-команда',  4),
  (8, 'Финансы',     1);

SELECT setval('departments_id_seq', (SELECT MAX(id) FROM departments));

INSERT INTO employees (department_id, full_name, position, hired_at) VALUES
  -- Разработка
  (2, 'Иванов Иван Иванович',        'Руководитель разработки',  '2018-03-01'),
  -- Backend
  (4, 'Петров Пётр Петрович',        'Senior Backend Developer',  '2019-07-15'),
  (4, 'Сидоров Алексей Николаевич',  'Backend Developer',         '2021-04-20'),
  -- Frontend
  (5, 'Козлова Мария Сергеевна',     'Senior Frontend Developer', '2020-09-01'),
  (5, 'Лебедев Роман Андреевич',     'Frontend Developer',        '2022-11-14'),
  -- Go-команда
  (7, 'Новиков Дмитрий Олегович',    'Go Developer',              '2023-01-09'),
  -- HR (без hired_at для проверки обработки nullable полей)
  (3, 'Смирнова Елена Викторовна',   'HR-менеджер',               NULL),
  -- Рекрутинг
  (6, 'Волков Андрей Михайлович',    'Рекрутер',                  '2021-11-30'),
  (6, 'Орлова Наталья Игоревна',     'Lead Recruiter',            '2020-02-10'),
  -- Финансы
  (8, 'Фёдоров Виктор Семёнович',    'Финансовый директор',       '2017-06-01'),
  (8, 'Белова Ксения Дмитриевна',    'Бухгалтер',                 NULL);

-- +goose Down
-- сотрудники удаляются автоматически при удалении подразделений благодаря ON DELETE CASCADE
DELETE FROM departments WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8);
SELECT setval('departments_id_seq', COALESCE((SELECT MAX(id) FROM departments), 0));
