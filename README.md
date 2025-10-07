# Simple Comment Tree

Сервис для работы с древовидными комментариями, поддерживающий неограниченную вложенность, полнотекстовый поиск, сортировку и постраничный вывод. Реализован с использованием Clean Architecture и PostgreSQL с рекурсивными CTE.

## Функциональность

### HTTP API
- **POST /comments** — создание комментария (с указанием родительского)
- **GET /comments?parent={id}** — получение комментария и всех вложенных
- **DELETE /comments/{id}** — удаление комментария и всех вложенных под ним

### Дополнительные возможности
- Постраничная навигация и сортировка
- Полнотекстовый поиск по комментариям (PostgreSQL FTS)
- Мягкое удаление с каскадным обновлением
- Web-интерфейс для взаимодействия


## Технологии

- **Backend**: Go 1.24, Gin, PostgreSQL
- **Database**: PostgreSQL с рекурсивными CTE и полнотекстовым поиском
- **Frontend**: Vanilla HTML/CSS/JavaScript
- **Architecture**: Clean Architecture с DI
- **Testing**: Testify + Mockery

## Установка и запуск

### Docker

```bash
# Запуск всех сервисов
make up

# Просмотр логов
make logs

# Остановка
make down
```

## API Документация

### Создание комментария
```http
POST /comments
Content-Type: application/json

{
  "parent_id": 1,  // опционально
  "content": "Текст комментария",
  "author": "Имя автора"
}
```

### Получение комментариев
```http
GET /comments?parent=0&page=1&limit=20&sort=created_at_asc&search=текст
```

**Параметры:**
- `parent` - ID родительского комментария (0 для корневых)
- `page` - номер страницы
- `limit` - количество на странице
- `sort` - сортировка (`created_at_asc`, `created_at_desc`)
- `search` - поисковый запрос

### Удаление комментария
```http
DELETE /comments/{id}
```

## База данных

### Схема таблицы
```sql
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER REFERENCES comments(id),
    content TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
```

### Индексы
- `idx_comments_parent_id` - для рекурсивных запросов
- `idx_comments_created_at` - для сортировки
- `idx_comments_deleted_at` - для фильтрации
- `idx_comments_content_gin` - для полнотекстового поиска

## Web-интерфейс

Доступен по адресу `http://localhost:8080`

**Возможности:**
- Просмотр дерева комментариев с визуальной вложенностью
- Создание новых комментариев и ответов
- Поиск по содержимому комментариев
- Отступы по уровням вложенности

## Тестирование

```bash
# Запуск всех тестов
make test

# Генерация моков
go generate ./...

# Форматирование кода
make fmt
```

## Производительность

- **Рекурсивные CTE** для эффективного обхода дерева
- **GIN индексы** для быстрого полнотекстового поиска
- **Connection pooling** и retry стратегии
- **Мягкое удаление** без потери данных