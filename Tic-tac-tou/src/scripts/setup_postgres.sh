#!/usr/bin/env bash

set -euo pipefail

CONTAINER_NAME="${CONTAINER_NAME:-some-postgres}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-1}"
POSTGRES_USER="${POSTGRES_USER:-postgres}"
POSTGRES_DB="${POSTGRES_DB:-tic_tac_toe}"
POSTGRES_PORT="${POSTGRES_PORT:-5432}"
POSTGRES_IMAGE="${POSTGRES_IMAGE:-postgres}"
SCHEMA_PATH="${SCHEMA_PATH:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/schema.sql}"

if ! command -v docker >/dev/null 2>&1; then
	echo "Docker не найден. Установи Docker и попробуй снова."
	exit 1
fi

if [ ! -f "$SCHEMA_PATH" ]; then
	echo "Не найден schema.sql по пути: $SCHEMA_PATH"
	exit 1
fi

container_exists() {
	docker ps -a --format '{{.Names}}' | grep -Fxq "$CONTAINER_NAME"
}

container_running() {
	docker ps --format '{{.Names}}' | grep -Fxq "$CONTAINER_NAME"
}

if ! container_exists; then
	echo "Создаю контейнер $CONTAINER_NAME из образа $POSTGRES_IMAGE..."
	docker run -d \
		--name "$CONTAINER_NAME" \
		-e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
		-e POSTGRES_USER="$POSTGRES_USER" \
		-p "$POSTGRES_PORT":5432 \
		"$POSTGRES_IMAGE" >/dev/null
elif ! container_running; then
	echo "Запускаю существующий контейнер $CONTAINER_NAME..."
	docker start "$CONTAINER_NAME" >/dev/null
else
	echo "Контейнер $CONTAINER_NAME уже запущен."
fi

echo "Жду готовности PostgreSQL..."
until docker exec "$CONTAINER_NAME" pg_isready -U "$POSTGRES_USER" >/dev/null 2>&1; do
	sleep 1
done

DB_EXISTS="$(docker exec "$CONTAINER_NAME" psql -U "$POSTGRES_USER" -tAc "SELECT 1 FROM pg_database WHERE datname = '$POSTGRES_DB';" | tr -d '[:space:]')"

if [ "$DB_EXISTS" != "1" ]; then
	echo "Создаю базу данных $POSTGRES_DB..."
	docker exec "$CONTAINER_NAME" psql -U "$POSTGRES_USER" -c "CREATE DATABASE $POSTGRES_DB;" >/dev/null
else
	echo "База данных $POSTGRES_DB уже существует."
fi

echo "Применяю схему из $SCHEMA_PATH..."
docker exec -i "$CONTAINER_NAME" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" < "$SCHEMA_PATH" >/dev/null

cat <<MSG
Готово.

PostgreSQL container: $CONTAINER_NAME
Database: $POSTGRES_DB
User: $POSTGRES_USER
Port: $POSTGRES_PORT

DATABASE_URL=postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable
MSG