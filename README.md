# TSV Processor

Сервис на Go для обработки файлов `.tsv`, загрузки данных в базу, генерации отчетов и предоставления API для доступа к данным.

## Описание

Сервис реализует следующие функции:

- Подключение к базе данных (PostgreSQL) и очереди сообщений (RabbitMQ) через конфиги из `.env`.
- Периодическая проверка указанной директории на новые `.tsv` файлы.
- Очередь на обработку файлов.
- Парсинг `.tsv` файлов и сохранение данных в базу.
- Генерация выходных файлов в формате PDF для каждого `unit_guid` из файла.
- Логирование ошибок парсинга в базу и формирование соответствующих файлов с ошибками.
- API для получения данных по `unit_guid` с поддержкой пагинации (`page` / `limit`).

## Требования

- Docker и Docker Compose
- Go 1.25
- PostgreSQL
- RabbitMQ

## Эндпоинты

 GET https://localhost:8000/api/v1/:guid?page=1&limit=20

Получение device_messages по unit_guid (с пагинцией)

## Конфигурация

Файл `.env`:

```env
APP_ENV=prod

POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=db

RABBIT_USERNAME=guest
RABBIT_PASSWORD=guest
RABBIT_VHOST=/
```
## Запуск
```bash
docker-compose -p tsv-processor up
```
