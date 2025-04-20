# Event Flow Platform

Event Flow - это современная платформа для управления мероприятиями, предоставляющая комплексное решение для организации, планирования и проведения различных типов событий.

## 🚀 Технологический стек

- **Backend**: Go (Golang)
- **База данных**: PostgreSQL
- **Контейнеризация**: Docker & Docker Compose
- **API**: RESTful API

## 📁 Структура проекта

```
backend/
├── cmd/                    # Точки входа приложения
├── internal/              # Внутренний код приложения
│   ├── core/             # Основная бизнес-логика
│   │   ├── models/      # Модели данных
│   │   ├── ports/       # Интерфейсы (ports & adapters)
│   │   └── services/    # Бизнес-сервисы
│   ├── infrastructure/   # Инфраструктурный код
│   └── config/          # Конфигурация приложения
├── migrations/           # Миграции базы данных
└── docker/              # Docker конфигурации
```

## 🔑 Основные функции

- Управление мероприятиями
- Регистрация участников
- Управление билетами
- Аналитика и отчеты
- Интеграция с платежными системами

## 🛠 Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/your-username/event-flow.git
```

2. Создайте файл .env на основе .env.example:
```bash
cp .env.example .env
```

3. Запустите приложение с помощью Docker Compose:
```bash
docker-compose up -d
```

## 📚 API Endpoints

### Проверка работоспособности
- `GET /health` - Проверка работоспособности сервера
  - Response: 200 OK
    ```json
    {
      "status": "ok"
    }
    ```

### Аутентификация
- `POST /auth/register` - Регистрация нового пользователя
  - Request Body:
    ```json
    {
      "email": "string",
      "password": "string",
      "name": "string",
      "role": "string",
      "description": "string",
      "activity_area": "string"
    }
    ```
  - Response: 200 OK
    ```json
    {
      "access_token": "string"
    }
    ```

- `POST /auth/login` - Вход в систему
  - Request Body:
    ```json
    {
      "email": "string",
      "password": "string"
    }
    ```
  - Response: 200 OK
    ```json
    {
      "access_token": "string"
    }
    ```

### Пользователи
- `GET /users/getInfo` - Получение информации о текущем пользователе
  - Headers: `Authorization: Bearer {token}`
  - Response: 200 OK
    ```json
    {
      "id": "string",
      "email": "string",
      "name": "string",
      "avatar": "string",
      "role": "string",
      "description": "string",
      "activity_area": "string",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
    ```

- `PUT /users/editInfo` - Редактирование информации пользователя
  - Headers: `Authorization: Bearer {token}`
  - Request Body:
    ```json
    {
      "email": "string",
      "name": "string",
      "avatar": "string"
    }
    ```
  - Response: 200 OK
    ```json
    {
      "id": "string",
      "email": "string",
      "name": "string",
      "avatar": "string",
      "role": "string",
      "description": "string",
      "activity_area": "string",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
    ```

- `POST /users/uploadAvatar` - Загрузка аватара пользователя
  - Headers: `Authorization: Bearer {token}`
  - Request Body:
    ```json
    {
      "base64_data": "string"
    }
    ```
  - Response: 200 OK
    ```json
    {
      "url": "string"
    }
    ```

- `GET /users/search` - Поиск пользователей по имени
  - Headers: `Authorization: Bearer {token}`
  - Query Parameters:
    - `name`: строка поиска (обязательный параметр)
  - Response: 200 OK
    ```json
    [
      {
        "id": "string",
        "name": "string",
        "avatar": "string"
      }
    ]
    ```

### Друзья
- `GET /users/friends` - Получение списка друзей
  - Headers: `Authorization: Bearer {token}`
  - Response: 200 OK
    ```json
    {
      "friends": [
        {
          "id": "string",
          "email": "string",
          "name": "string",
          "avatar": "string",
          "role": "string",
          "description": "string",
          "activity_area": "string",
          "created_at": "datetime",
          "updated_at": "datetime"
        }
      ]
    }
    ```

- `GET /users/friends/incoming` - Получение входящих запросов в друзья
  - Headers: `Authorization: Bearer {token}`
  - Response: 200 OK
    ```json
    [
      {
        "id": "string",
        "name": "string",
        "avatar": "string"
      }
    ]
    ```

- `POST /users/friends/request` - Отправка запроса в друзья
  - Headers: `Authorization: Bearer {token}`
  - Request Body:
    ```json
    {
      "to_id": "string"
    }
    ```
  - Response: 200 OK
    ```json
    {
      "id": "string",
      "from_id": "string",
      "to_id": "string",
      "status": "string",
      "created_at": "datetime",
      "from_name": "string",
      "from_avatar": "string"
    }
    ```

- `PUT /users/friends/respond` - Ответ на запрос в друзья
  - Headers: `Authorization: Bearer {token}`
  - Request Body:
    ```json
    {
      "friend_id": "string",
      "accept": "boolean"
    }
    ```
  - Response: 200 OK

- `DELETE /users/friends/:friendId` - Удаление друга
  - Headers: `Authorization: Bearer {token}`
  - Response: 200 OK

## 📊 База данных

Проект использует PostgreSQL в качестве основной базы данных. Миграции находятся в директории `migrations/`.

## 🔒 Безопасность

- JWT аутентификация
- HTTPS шифрование
- Защита от SQL-инъекций
- Rate limiting
- CORS политики

## 🤝 Вклад в проект

Мы приветствуем вклад в развитие проекта! Пожалуйста, следуйте этим шагам:

1. Форкните репозиторий
2. Создайте ветку для ваших изменений
3. Внесите изменения
4. Создайте Pull Request

## 📝 Лицензия

MIT License