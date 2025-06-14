## LoadBalancer

## **Описание проекта**

Данный проект представляет собой балансировщик нагрузки (Load Balancer), который распределяет входящие HTTP-запросы между несколькими бэкенд-серверами. В проекте реализован функционал:

1. **Балансировка нагрузки**:
   - Реализован алгоритм Round-Robin для распределения запросов.

2. **Rate-Limiting**:
   - Ограничение частоты запросов с использованием алгоритма Token Bucket.
   - Настройка лимитов для каждого клиента через базу данных.

3. **Интеграция с PostgreSQL**:
   - Хранение информации о клиентах и их лимитах в базе данных.
   - CRUD API для управления клиентами(POST, DELETE).

4. **Контейнеризация**:
   - Dockerfile для сборки приложения.
   - Docker Compose для запуска балансировщика, БД и тестовых бэкендов.

---

## **Сборка и запуск**

### **1. Клонирование репозитория**
```bash
git clone https://github.com/v1adis1av28/LoadBalancer.git
cd LoadBalancer
```

### **2. Сборка и запуск проекта**
Соберите Docker-образы с помощью Docker Compose:
```bash
docker-compose up --build
```

После запуска:
- Балансировщик будет доступен на `localhost:8080`.
- Тестовые бэкенды будут доступны на портах `8081`, `8082` и `8083`.

### **4. Остановка проекта**
Чтобы остановить все сервисы:
```bash
CTRL + C
docker-compose down
```

---

## **Функционал**

### **1. Балансировка нагрузки**
- Входящие запросы распределяются между бэкендами с использованием алгоритма Round-Robin.

Пример запроса:
```bash
curl http://localhost:8080/
```

### **2. Rate-Limiting**
- Каждый клиент имеет ограничение на количество запросов в секунду.

Пример ошибки при превышении лимита:
```json
{
  "code": "429",
  "message": "Rate limit exceeded"
}
```

### **3. Управление клиентами**
API предоставляет возможность управлять клиентами через HTTP-запросы.

#### Добавление клиента
```bash
curl -X POST http://localhost:8080/clients \
-H "Content-Type: application/json" \
-d '{
  "client_id": "127.0.0.12",
  "capacity": 20,
  "rate_per_sec": 100
}'
```

#### Удаление клиента
```bash
curl -X DELETE http://localhost:8080/clients/127.0.0.12
```

---

## **Структура проекта**

Работа над структурой проекта велась основываясь на структуре стандартного проекта описанного -> https://github.com/golang-standards/project-layout

---

## **Конфигурация**

### **1. Файл `config.json`**
Файл содержит основные настройки балансировщика:
```json
{
  "port": "8080",
  "backends": [
    "http://backend1:8081",
    "http://backend2:8082",
    "http://backend3:8083"
  ]
}
```

### **2. Переменные окружения PostgreSQL**
Настройки подключения к базе данных задаются в `docker-compose.yml`:
```yaml
environment:
  POSTGRES_USER: postgres
  POSTGRES_PASSWORD: postgres
  POSTGRES_DB: loadbalancer
```

---

## **Логирование**

Логирование в приложении ведется с помощью slog, все логи записываются в файл app.log в виде JSON . Пример записи:
```json
{
  "time": "2025-04-29T14:55:01.383984156Z",
  "level": "INFO",
  "msg": "Server start on port %s",
  "!BADKEY": "8080"
}
```
