# Финальный проект 1 семестра (простой уровень сложности)

REST API сервис для загрузки и выгрузки данных о ценах.

## Требования к системе

### Предпочтительные OS
- Linux/macOS/Windows

### Минимальное аппратаное обеспечение
- Процессор: 2 core
- RAM: 2GB
- HDD: 4Gb

## Установка и запуск

1. Установите Go
2. Установите PostgreSQL, настройка базы данных
```bash
sudo apt update
sudo apt install postgresql
sudo su - postgres
psql
CREATE USER validator WITH PASSWORD 'val1dat0r';
CREATE DATABASE "project-sem-1" OWNER validator;
\q
```
3. Клонируйте репозиторий и запустите скрипт подготовки
```bash
git clone git@github.com:panov-a-st/itmo-devops-sem1-project-template.git
cd itmo-devops-sem1-project-template
bash ./scripts/prepare.sh
```
4. Запустите локальный сервер
```bash
bash ./scripts/run.sh
```

## Тестирование

### Тестирование API-запросов:

```bash
bash ./scripts/tests.sh 1
```

### Скрипт `tests.sh` проверяет:

- Корректность работы эндпоинта загрузки `POST /api/v0/prices`
- Корректность работы эндпоинта выгрузки `GET /api/v0/prices`

### Примеры запросов
- Получение zip архива с помощью GET запроса:
```bash
curl -X GET -o output_test.zip http://localhost:8080/api/v0/prices
```
- Загрузка данных из архива с помощью POST
```bash
curl -X POST -F "file=@sample_data.zip" http://localhost:8080/api/v0/prices
```

## Контакт

В случае вопросов можно обращаться:
- Telegram @orionisman
