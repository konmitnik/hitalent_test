# hitalent_test — API организационной структуры

REST API приложения для управления деревом подразделений и сотрудниками. Написано на Go с использованием net/http, GORM, PostgreSQL, goose (миграции), Docker.

---

## Быстрый старт

```bash
docker-compose up --build
```

API приложение будет доступно по ссылке **http://localhost:8080** как только все 3 службы буду в порядке:

| Служба    | Роль                                      |
|-----------|-------------------------------------------|
| `db`      | PostgreSQL 17                             |
| `migrate` | Выполняет goose миграции при запуске      |
| `app`     | Go HTTP сервер на порте 8080              |

---

## Руководство по API

### Подразделения (Departments)

#### Создание подразделения
```
POST /departments/
Content-Type: application/json

{ "name": "Engineering", "parent_id": null }
```

#### Получение подразделения (с поддеревом + сотрудниками)
```
GET /departments/{id}?depth=2&include_employees=true
```
| Параметр запроса   | По умолчанию | Заметка                               |
|--------------------|--------------|---------------------------------------|
| `depth`            | `1`          | Уровень вложенности детей, максимум 5 |
| `include_employees`| `true`       | Добавляет список сотрудников в ответ  |

#### Обновление подразделения (Название / родитель)
```
PATCH /departments/{id}
Content-Type: application/json

{ "name": "Backend", "parent_id": 3 }
```
Чтобы переместить подразделение в корень передайте `"parent_id": null`.

#### Удаление подразделения
```
DELETE /departments/{id}?mode=cascade
DELETE /departments/{id}?mode=reassign&reassign_to_department_id=5
```
| `Режим`    | Поведение                                                              |
|------------|------------------------------------------------------------------------|
| `cascade`  | Удаляет департамент и все его поддеревья и сотрудников                 |
| `reassign` | Переназначает сотрудников в указанный департамент, а потом удаляет     |

---

### Сотрудники

#### Создание сотрудника
```
POST /departments/{id}/employees/
Content-Type: application/json

{ "full_name": "Alice Smith", "position": "Engineer", "hired_at": "2023-06-01" }
```
`hired_at` опциональный параметр (дата в формате ISO 8601).

---

## Валидация

- `name` — обязательный параметр, 1–200 символов, пробелы по краям убраны; должно быть уникальное имя в рамках 1 родителя
- `full_name` / `position` — обязательный параметр, 1–200 символов
- Перемещение подразделения внутрь своего поддерева вернет  **409 Conflict**
- Создание подразделения с несуществующим `parent_id` вернет **404**
- Создание сотрудника в несуществующем подразделении вернет **404**

---

## Структура проекта

```
cmd/api/          точка входа
internal/
  config/         соединения с БД
  handlers/       обработка HTTP запросов
  helpers/        JSON response helper
  models/         GORM модели (Department, Employee)
  repository/     запросы к БД
migrations/       goose SQL миграции
Dockerfile
docker-compose.yaml
```

---

## Локальная разработка (без Docker)

Требования: Go 1.21+, PostgreSQL запускается локально.

```bash
# выполнение миграций
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=test_db sslmode=disable" up

# запуск сервера
go run ./cmd/api/
```

Переменные окружения (все опциональны, показаны дефолтные значения):

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=test_db
DB_SSLMODE=disable
```
