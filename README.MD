# Organization Structure API

REST API сервис для управления организационной структурой предприятия (департаменты и сотрудники).
Реализовано в рамках тестового задания.

## 🛠 Технологический стек

*   **Язык:** Go 1.25
*   **Архитектура:** Clean Architecture (Domain -> Service -> Repository -> Handler)
*   **База данных:** PostgreSQL
*   **ORM:** GORM
*   **Миграции:** Goose
*   **Конфигурация:** Cleanenv (YAML + Environment variables)
*   **Тесты:** Testify
*   **Деплой:** Docker & Docker Compose

## 🚀 Быстрый запуск

Для запуска требуется установленный **Docker** и **Docker Compose**.

1.  Клонируйте репозиторий:
    ```bash
    git clone https://github.com/Valery223/org-structure-api.git
    cd org-structure-api
    ```

2.  Создайте файл с переменными окружения (или используйте дефолтный из docker-compose):
    *Проект уже настроен на работу "из коробки" с настройками в `docker-compose.yml` и `.env`.*

3.  Запустите приложение:
    ```bash
    docker-compose up --build
    ```

API будет доступно по адресу: `http://localhost:8080`.
Миграции применятся автоматически при старте.

## 🧪 Тестирование

Запуск юнит-тестов (бизнес-логика):
```bash
go test ./internal/service/... -v
```

## 📚 API Документация

### Департаменты

*   **POST** `/departments/` - Создать департамент
    ```json
    { "name": "Backend", "parent_id": null }
    ```
*   **GET** `/departments/{id}?depth=1&include_employees=true` - Получить дерево
*   **PATCH** `/departments/{id}` - Переместить или переименовать
    ```json
    { "parent_id": 5 }
    ```
*   **DELETE** `/departments/{id}?mode=cascade` - Каскадное удаление
*   **DELETE** `/departments/{id}?mode=reassign&reassign_to_department_id=10` - Удаление с переносом сотрудников

### Сотрудники

*   **POST** `/departments/{id}/employees/` - Добавить сотрудника
    ```json
    {
      "full_name": "Ivan Ivanov",
      "position": "Developer",
      "hired_at": "2025-10-01"
    }
    ```

## 🏗 Архитектура

Проект следует принципам чистой архитектуры:
*   `cmd/app` - Точка входа (Main).
*   `internal/domain` - Бизнес-сущности и интерфейсы.
*   `internal/service` - Бизнес-логика (проверка циклов, валидация).
*   `internal/repository` - Работа с БД (PostgreSQL implementation).
*   `internal/handler` - HTTP обработчики (REST).
*   `internal/config` - Загрузка конфигурации.
*   `migrations` - SQL файлы миграций.